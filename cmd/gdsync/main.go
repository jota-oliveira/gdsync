package main

import (
	"github.com/jota-oliveira/gdsync/internal/infra/local"
	usecases "github.com/jota-oliveira/gdsync/internal/use-cases"
)

func main() {
	// reader := bufio.NewReader(os.Stdin)

	// fmt.Println("Choose the folder to sync: ")

	// path, err := reader.ReadString('\n')

	// if err != nil {
	// 	fmt.Errorf("Sorry! We could'nt find you file", err)
	// }

	// path = strings.TrimSpace(path)

	// For testing purposes
	path := "/home/joao/Projects/teste"

	fs := local.NewFileSystem()

	syncUsecase := usecases.NewSyncFile(path, fs)

	syncUsecase.Execute()
}
