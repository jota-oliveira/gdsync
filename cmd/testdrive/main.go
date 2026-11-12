package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jota-oliveira/gdsync/internal/domains"
	driveinfra "github.com/jota-oliveira/gdsync/internal/infra/drive"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	client, err := driveinfra.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	service, err := drive.NewService(
		ctx,
		option.WithHTTPClient(client),
	)
	if err != nil {
		log.Fatal(err)
	}

	result, err := service.Files.List().
		PageSize(10).
		Fields("nextPageToken, files(id, name)").
		Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Google Drive conectado!")

	for _, file := range result.Files {
		fmt.Printf("%s - %s\n", file.Id, file.Name)
	}

	driveRepository := driveinfra.NewRepository(service)

	//	file := &domains.File{
	//	Path: "test-upload.txt",
	//}

	/*
		driveID, err := repo.Upload(file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Upload realizado!")
		fmt.Println("Drive ID:", driveID)
	*/
	folder := &domains.Folder{
		Name: "gdsync-test",
		Path: "/tmp/gdsync-test",
	}

	driveID, err := driveRepository.CreateFolder(folder, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Folder created:", driveID)
}
