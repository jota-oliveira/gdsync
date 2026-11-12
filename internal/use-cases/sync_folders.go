package usecases

import (
	"fmt"

	"github.com/jota-oliveira/gdsync/internal/interfaces"
	"github.com/jota-oliveira/gdsync/internal/repositories"
)

type SyncFolder struct {
	localFileSystem  interfaces.FileSystem
	folderRepository repositories.FolderRepository
	driveRepository  repositories.DriveRepository
	path             string
}

func NewSyncFolder(
	path string,
	fs interfaces.FileSystem,
	folderRepository repositories.FolderRepository,
	driveRepository repositories.DriveRepository,
) *SyncFolder {
	return &SyncFolder{
		path:             path,
		localFileSystem:  fs,
		folderRepository: folderRepository,
		driveRepository:  driveRepository,
	}
}

func (s *SyncFolder) Execute() error {
	localFolders, err := s.localFileSystem.ListFolders(s.path)
	if err != nil {
		return fmt.Errorf("list local folders: %w", err)
	}

	savedFolders, err := s.folderRepository.ListFolders(s.path)
	if err != nil {
		return fmt.Errorf("list saved folders: %w", err)
	}

	// Relaciona o ID local ao ID remoto.
	driveIDs := make(map[string]string)

	for _, folder := range savedFolders {
		if folder.DriveID != "" {
			driveIDs[folder.ID] = folder.DriveID
		}
	}

	for _, folder := range localFolders {
		// Já existe no Drive?
		if driveIDs[folder.ID] != "" {
			continue
		}

		var parentDriveID string

		if folder.ParentID != "" {
			var found bool

			parentDriveID, found = driveIDs[folder.ParentID]
			if !found {
				return fmt.Errorf(
					"remote parent not found for %s",
					folder.Path,
				)
			}
		}

		// Cria no Drive.
		driveID, err := s.driveRepository.CreateFolder(
			folder,
			parentDriveID,
		)
		if err != nil {
			return fmt.Errorf(
				"create remote folder %s: %w",
				folder.Path,
				err,
			)
		}

		folder.DriveID = driveID

		// Persiste no SQLite.
		if err := s.folderRepository.NewFolder(*folder); err != nil {
			return fmt.Errorf(
				"save folder %s: %w",
				folder.Path,
				err,
			)
		}

		driveIDs[folder.ID] = driveID
	}

	return nil
}
