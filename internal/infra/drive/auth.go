package drive

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	driveapi "google.golang.org/api/drive/v3"
)

func NewClient(ctx context.Context) (*http.Client, error) {
	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("read credentials: %w", err)
	}

	config, err := google.ConfigFromJSON(
		credentials,
		driveapi.DriveFileScope,
	)
	if err != nil {
		return nil, fmt.Errorf("parse credentials: %w", err)
	}

	token, err := loadToken("token.json")
	if errors.Is(err, os.ErrNotExist) {
		token, err = authorize(ctx, config)
		if err != nil {
			return nil, err
		}

		if err := saveToken("token.json", token); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return config.Client(ctx, token), nil
}

func authorize(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	config.RedirectURL = fmt.Sprintf(
		"http://localhost:%d/callback",
		listener.Addr().(*net.TCPAddr).Port,
	)

	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		return nil, err
	}
	state := hex.EncodeToString(stateBytes)
	verifier := oauth2.GenerateVerifier()

	result := make(chan string, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		if authErr := r.URL.Query().Get("error"); authErr != "" {
			fmt.Fprintln(w, "Authorization denied:", authErr)
			result <- ""
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "Authorized! You can close this tab.")
		result <- code
	})

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Close()

	url := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)

	fmt.Println("Open this URL in your browser:")
	fmt.Println(url)

	select {
	case code := <-result:
		if code == "" {
			return nil, errors.New("authorization denied")
		}

		return config.Exchange(
			ctx, code, oauth2.VerifierOption(verifier),
		)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func loadToken(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func saveToken(path string, token *oauth2.Token) error {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}
