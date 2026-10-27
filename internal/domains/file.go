// Package domains - Files
package domains

import (
	"fmt"
	"os"
)

type File struct {
	ID       string
	DriveID  string
	Name     string
	Path     string
	Size     int64
	ModTime  int64
	Checksum string
}

func NewFile() {
}

func (f *File) generateChecksum(fileInfo os.FileInfo) {
	f.Checksum = fmt.Sprintf("%d:%d", fileInfo.Size(), fileInfo.ModTime().UnixNano())
}
