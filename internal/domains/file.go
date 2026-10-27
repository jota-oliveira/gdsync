// Package domains - Files
package domains

import (
	"fmt"
	"os"
)

type File struct {
	ID        string
	Name      string
	Path      string
	Checksum  string
	UpdatedAt int64
}

func NewFile() {

}

func (f *File) generateChecksum(fileInfo os.FileInfo) {
	f.Checksum = fmt.Sprintf("%d:%d", fileInfo.Size(), fileInfo.ModTime().UnixNano())
}
