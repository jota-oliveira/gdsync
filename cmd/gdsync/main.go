package main

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/jota-oliveira/gdsync/internal/infra/local"
	usecases "github.com/jota-oliveira/gdsync/internal/use-cases"

	driveinfra "github.com/jota-oliveira/gdsync/internal/infra/drive"
	driveapi "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {
	// reader := bufio.NewReader(os.Stdin)

	// fmt.Println("Choose the folder to sync: ")

	// path, err := reader.ReadString('\n')

	// if err != nil {
	// 	fmt.Errorf("Sorry! We could'nt find you file", err)
	// }

	// path = strings.TrimSpace(path)

	fs := local.NewFileSystem()

	fileRepository, err := local.NewFileRepo()
	if err != nil {
		log.Fatal("Could not initiate the database")
	}

	path, err := fileRepository.HasConfigurations()

	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	path = strings.TrimSpace(path)

	/*
		if path == "" {
			reader := bufio.NewReader(os.Stdin)

			fmt.Println("Choose the folder to sync: ")

			path, err = reader.ReadString('\n')
			if err != nil {
				log.Fatal("Could not read the path name")
			}

			path = strings.TrimSpace(path)

			path, err = filepath.Abs(path)
			if err != nil {
				log.Fatal(err)
			}
	*/
	err = fileRepository.Init(path)
	if err != nil {
		log.Fatal(err)
	}
	//}

	ctx := context.Background()

	client, err := driveinfra.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	service, err := driveapi.NewService(
		ctx,
		option.WithHTTPClient(client),
	)
	if err != nil {
		log.Fatal(err)
	}

	driveRepository := driveinfra.NewRepository(service)

	/*
		syncUsecase := usecases.NewSyncFile(
			path,
			fs,
			fileRepository,
			driveRepository,
		)*/

	/*
		if err := syncUsecase.Execute(); err != nil {
			log.Fatal(err)
		}*/

	folderSync := usecases.NewSyncFolder(
		path,
		fs,
		fileRepository,
		driveRepository,
	)

	if err := folderSync.Execute(); err != nil {
		log.Fatal(err)
	}
}
