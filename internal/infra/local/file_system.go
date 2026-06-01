package local

import (
	"os"

	"github.com/jota-oliveira/gdsync/internal/domains"
)

type LocalFileSystem struct{}

func NewFileSystem() *LocalFileSystem {
	return &LocalFileSystem{}
}

func (fs *LocalFileSystem) ListFiles(path string) ([]*domains.File, error) {
	entries, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var files []*domains.File

	for _, entry := range entries {
		info, err := entry.Info()

		if err != nil {
			continue
		}

		files = append(files, &domains.File{
			Name:      info.Name(),
			Path:      path + "/" + info.Name(),
			Checksum:  "123",
			ID:        info.Name(),
			UpdatedAt: info.ModTime().Unix(),
		})
	}

	return files, nil
}
