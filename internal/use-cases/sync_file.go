package usecases

import (
	"fmt"

	"github.com/jota-oliveira/gdsync/internal/interfaces"
)

type SyncFile struct {
	localFileSystem interfaces.FileSystem
	path            string
}

func NewSyncFile(path string, fs interfaces.FileSystem) *SyncFile {
	return &SyncFile{
		path:            path,
		localFileSystem: fs,
	}
}

func (s *SyncFile) Execute() error {
	filesInFolder, err := s.localFileSystem.ListFiles(s.path)

	if err != nil {
		return err
	}

	for _, file := range filesInFolder {
		fmt.Println(file.Name)
	}

	return nil
}
