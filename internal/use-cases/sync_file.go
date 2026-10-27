package usecases

import (
	"fmt"

	"github.com/jota-oliveira/gdsync/internal/domains"
	"github.com/jota-oliveira/gdsync/internal/interfaces"
	"github.com/jota-oliveira/gdsync/internal/repositories"
)

type SyncFile struct {
	localFileSystem interfaces.FileSystem
	path            string
	repository      repositories.FileRepository
	driveRepository repositories.DriveRepository
}

func NewSyncFile(path string, fs interfaces.FileSystem, fileRepository repositories.FileRepository) *SyncFile {
	return &SyncFile{
		path:            path,
		localFileSystem: fs,
		repository:      fileRepository,
	}
}

func (s *SyncFile) generateSyncPlan(filesInFolder []*domains.File, filesInDB []*domains.File) domains.SyncPlan {
	filesDatabaseMap := make(map[string]*domains.File)
	filesFolderMap := make(map[string]*domains.File)

	for _, file := range filesInDB {
		filesDatabaseMap[file.ID] = file
	}

	for _, file := range filesInFolder {
		filesFolderMap[file.ID] = file
	}

	plan := domains.SyncPlan{}

	for fileID, file := range filesFolderMap {
		fileInDB := filesDatabaseMap[fileID]

		if fileInDB != nil {
			if file.Checksum != fileInDB.Checksum {
				plan.Modified = append(plan.Modified, file)
			}
		} else {
			plan.New = append(plan.New, file)
		}
	}

	for fileID, file := range filesDatabaseMap {
		fileInFolder := filesFolderMap[fileID]

		if fileInFolder == nil {
			plan.Deleted = append(plan.Deleted, file)
		}
	}

	return plan
}

func (s *SyncFile) Execute() error {
	filesInFolder, err := s.localFileSystem.ListFiles(s.path)
	if err != nil {
		return err
	}

	filesInDB, err := s.repository.ListFiles(s.path)
	if err != nil {
		return err
	}

	plan := s.generateSyncPlan(filesInFolder, filesInDB)

	// Apenas para visualizar o resultado em tempo de desenvolvimento:
	fmt.Printf("Novos: %d | Modificados: %d | Deletados: %d\n", len(plan.New), len(plan.Modified), len(plan.Deleted))

	if len(plan.Deleted) > 0 {
		for _, file := range plan.Deleted {
			err := s.repository.Delete(file.ID)
			if err != nil {
				return err
			}
		}
	}

	if len(plan.New) > 0 {
		for _, file := range plan.New {
			err := s.repository.NewFile(*file)
			if err != nil {
				return err
			}
		}
	}

	if len(plan.Modified) > 0 {
		for _, file := range plan.Modified {
			err := s.repository.UpdateFile(*file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
