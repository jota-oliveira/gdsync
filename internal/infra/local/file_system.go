package local

import (
	"fmt"
	"os"
	"syscall"

	"github.com/jota-oliveira/gdsync/internal/domains"
)

type LocalFileSystem struct{}

func NewFileSystem() *LocalFileSystem {
	return &LocalFileSystem{}
}

func generateChecksum(fileInfo os.FileInfo) string {
	// hash := sha256.New()
	return fmt.Sprintf("%d:%d", fileInfo.Size(), fileInfo.ModTime().UnixNano())

	// _, err := hash.Write([]byte(fileCondensedData))
	// if err != nil {
	// 	return "", err
	// }

	// return hex.EncodeToString(hash.Sum(nil)), nil
}

func generateID(fileInfo os.FileInfo) string {
	stat := fileInfo.Sys().(*syscall.Stat_t)

	return fmt.Sprintf("%d:%d", stat.Dev, stat.Ino)
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

		checksum := generateChecksum(info)
		id := generateID(info)

		files = append(files, &domains.File{
			ID:        id,
			Name:      info.Name(),
			Path:      path + "/" + info.Name(),
			Checksum:  checksum,
			UpdatedAt: info.ModTime().Unix(),
		})
	}

	return files, nil
}
