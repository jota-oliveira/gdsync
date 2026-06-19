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
}

func NewSyncFile(path string, fs interfaces.FileSystem, fileRepository repositories.FileRepository) *SyncFile {
	return &SyncFile{
		path:            path,
		localFileSystem: fs,
		repository:      fileRepository,
	}
}

func (s *SyncFile) syncronize(filesInFolder []*domains.File, filesInDB []*domains.File) {
	for i, file := range filesInFolder {
		fmt.Println(file, i)
	}
}

func (s *SyncFile) Execute() error {

	// if s.repository.HasConfigurations() {

	// } else {
	// 	if err := s.repository.Init(s.path); err != nil {
	// 		log.Fatal("Could not init the database config")
	// 	}
	// }

	filesInFolder, err := s.localFileSystem.ListFiles(s.path)

	if err != nil {
		return err
	}

	filesInDB, err := s.repository.ListFiles(s.path)

	s.syncronize(filesInFolder, filesInDB)

	if err != nil {
		return err
	}

	for _, file := range filesInFolder {
		fmt.Println(file.Checksum)
	}

	return nil
}
