package interfaces

import (
	"github.com/jota-oliveira/gdsync/internal/domains"
)

type CloudStorage interface {
	ListFiles() ([]domains.File, error)

	UploadFile(file domains.File) error
	RemoveFile(file domains.File) error
	ManageConflit(file domains.File) error

	UploadFiles(files []domains.File) []error
	RemoveFiles(files []domains.File) []error
}
