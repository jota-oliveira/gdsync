// Package google - Client for Google Drive

package google

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveClient struct {
	service *drive.Service
}

func NewDriveClient(credentialsPath *google.Credentials) (*DriveClient, error) {
	ctx := context.Background()

	srv, err := drive.NewService(
		ctx,
		option.WithCredentials(credentialsPath),
		option.WithScopes(drive.DriveScope),
	)

	if err != nil {
		return nil, err
	}

	return &DriveClient{
		service: srv,
	}, nil
}
