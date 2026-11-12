package repositories

import "github.com/jota-oliveira/gdsync/internal/domains"

type FolderRepository interface {
	NewFolder(folder domains.Folder) error
	ListFolders(root string) ([]*domains.Folder, error)
	// UpdateFolder(folder domains.Folder) error
}
