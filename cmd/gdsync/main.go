package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
	// path := "/home/joao/Projects/teste"

	fs := local.NewFileSystem()

	fileRepository, err := local.NewFileRepo()

	if err != nil {
		log.Fatal("Could not initiate the database")
	}

	path, err := fileRepository.HasConfigurations()

	if err != nil {
		fmt.Println(err)
	}

	if path == "" {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Choose the folder to sync: ")

		path, err = reader.ReadString('\n')

		if err != nil {
			log.Fatal("Could not read the path name")
		}

		err = fileRepository.Init(path)

		if err != nil {
			log.Fatal(err)
		}
	}

	syncUsecase := usecases.NewSyncFile(path, fs, fileRepository)

	syncUsecase.Execute()
}
