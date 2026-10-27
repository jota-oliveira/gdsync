package repositories

import "github.com/jota-oliveira/gdsync/internal/domains"

type DriveRepository interface {
	Upload(file *domains.File) (string, error)
}
