package drive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jota-oliveira/gdsync/internal/domains"
	driveapi "google.golang.org/api/drive/v3"
)

type Repository struct {
	service *driveapi.Service
}

func NewRepository(service *driveapi.Service) *Repository {
	return &Repository{
		service: service,
	}
}

func (r *Repository) Upload(file *domains.File) (string, error) {
	localFile, err := os.Open(file.Path)
	if err != nil {
		return "", fmt.Errorf("open file %q: %w", file.Path, err)
	}
	defer localFile.Close()

	metadata := &driveapi.File{
		Name: filepath.Base(file.Path),
	}

	uploaded, err := r.service.Files.Create(metadata).
		Media(localFile).
		Fields("id").
		Do()
	if err != nil {
		return "", fmt.Errorf("upload %q: %w", file.Path, err)
	}

	if uploaded.Id == "" {
		return "", fmt.Errorf("drive returned empty ID for %q", file.Path)
	}

	return uploaded.Id, nil
}

func (r *Repository) CreateFolder(
	folder *domains.Folder,
	parentDriveID string,
) (string, error) {
	metadata := &driveapi.File{
		Name:     folder.Name,
		MimeType: "application/vnd.google-apps.folder",
	}

	if parentDriveID != "" {
		metadata.Parents = []string{parentDriveID}
	}

	created, err := r.service.Files.Create(metadata).
		Fields("id").
		Do()
	if err != nil {
		return "", fmt.Errorf(
			"could not create folder %q: %w",
			folder.Path,
			err,
		)
	}

	if created.Id == "" {
		return "", fmt.Errorf(
			"Drive returned empty ID for %q",
			folder.Path,
		)
	}

	return created.Id, nil
}
