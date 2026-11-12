package local

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/jota-oliveira/gdsync/internal/domains"
)

func (fr *FileRepo) NewFolder(folder domains.Folder) error {
	_, err := fr.db.Exec(
		`
        INSERT INTO folders (
            local_id,
            drive_id,
            name,
            path,
            parent_id
        )
        VALUES (?, ?, ?, ?, ?)
    `,
		folder.ID,
		folder.DriveID,
		folder.Name,
		folder.Path,
		folder.ParentID,
	)
	if err != nil {
		return fmt.Errorf("insert folder %q: %w", folder.Path, err)
	}

	return nil
}

func (fr *FileRepo) ListFolders(root string) ([]*domains.Folder, error) {
	root = filepath.Clean(root)

	rows, err := fr.db.Query(`
    SELECT
        local_id,
        drive_id,
        name,
        path,
        parent_id
    FROM folders
    WHERE path = ?
       OR substr(path, 1, length(?) + 1) = ? || '/'
`, root, root, root)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*domains.Folder

	for rows.Next() {
		folder := &domains.Folder{}
		var driveID sql.NullString
		var parentID sql.NullString

		if err := rows.Scan(
			&folder.ID,
			&driveID,
			&folder.Name,
			&folder.Path,
			&parentID,
		); err != nil {
			return nil, err
		}

		if driveID.Valid {
			folder.DriveID = driveID.String
		}

		if parentID.Valid {
			folder.ParentID = parentID.String
		}

		folders = append(folders, folder)
	}

	return folders, rows.Err()
}
