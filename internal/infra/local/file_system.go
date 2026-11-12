package local

import (
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
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

func (fs *LocalFileSystem) ListFiles(
	root string,
) ([]*domains.File, error) {
	var files []*domains.File

	err := filepath.WalkDir(
		root,
		func(path string, entry iofs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if entry.IsDir() {
				return nil
			}

			info, err := entry.Info()
			if err != nil {
				return err
			}

			id := generateID(info)

			files = append(files, &domains.File{
				ID:       id,
				Name:     info.Name(),
				Path:     path,
				Size:     info.Size(),
				ModTime:  info.ModTime().UnixNano(),
				Checksum: generateChecksum(info),
			})

			return nil
		},
	)

	return files, err
}

func (lfs *LocalFileSystem) ListFolders(
	root string,
) ([]*domains.Folder, error) {
	var folders []*domains.Folder

	folderIDs := make(map[string]string)

	root = filepath.Clean(root)

	err := filepath.WalkDir(
		root,
		func(path string, entry iofs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !entry.IsDir() {
				return nil
			}

			info, err := entry.Info()
			if err != nil {
				return err
			}

			id := generateID(info)

			var parentID string

			if path != root {
				parentPath := filepath.Dir(path)
				parentID = folderIDs[parentPath]
			}

			folderIDs[path] = id

			folders = append(folders, &domains.Folder{
				ID:       id,
				Name:     entry.Name(),
				Path:     path,
				ParentID: parentID,
			})

			return nil
		},
	)

	return folders, err
}
