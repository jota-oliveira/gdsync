package repositories

import "github.com/jota-oliveira/gdsync/internal/domains"

type FileRepository interface {
	HasConfigurations() (string, error)
	Init(path string) error
	NewFile(file domains.File) error
	Delete(id string) error
	UpdateFile(file domains.File) error
	ListFiles() ([]domains.File, error)
}
