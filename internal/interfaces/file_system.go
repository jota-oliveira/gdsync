// internal/interfaces/filesystem.go

package interfaces

import "github.com/jota-oliveira/gdsync/internal/domains"

type FileSystem interface {
	ListFiles(path string) ([]*domains.File, error)
}
