package local

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/jota-oliveira/gdsync/internal/domains"
)

type FileRepo struct {
	db *sql.DB
}

func NewFileRepo() (*FileRepo, error) {
	db, err := sql.Open("sqlite", "gdsync.db")

	if err != nil {
		return nil, err
	}

	return &FileRepo{
		db: db,
	}, nil
}

func (fr *FileRepo) Delete(id string) error {
	return nil
}

func (fr *FileRepo) HasConfigurations() (string, error) {
	var path string

	err := fr.db.QueryRow(`
		SELECT sync_path
		FROM config
		WHERE id = 1
	`).Scan(&path)

	return path, err
}

func (fr *FileRepo) Init(path string) error {
	_, err := fr.db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			id INTEGER PRIMARY KEY,
			sync_path TEXT NOT NULL
		);
	`)

	if err != nil {
		return err
	}

	_, err = fr.db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			local_id TEXT PRIMARY KEY,
			path TEXT NOT NULL,
			size INTEGER NOT NULL,
			mod_time INTEGER NOT NULL,
			drive_id TEXT
		);
	`)

	if err != nil {
		return err
	}

	_, err = fr.db.Exec(`
		INSERT OR REPLACE INTO config (
			id,
			sync_path
		)
		VALUES (1, ?)
	`, path)

	return err
}

func (fr *FileRepo) NewFile(file domains.File) error {
	return nil
}

func (fr *FileRepo) UpdateFile(file domains.File) error {
	return nil
}

func (fr *FileRepo) ListFiles(path string) ([]*domains.File, error) {
	rows, err := fr.db.Query(`
		SELECT
			local_id as ID,
			path as Path,
			size as Size,
			mod_time as ModTime,
			drive_id as DriveId
		FROM files
		WHERE path = ?
	`, strings.Trim(path, "\n"))

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var files []*domains.File

	for rows.Next() {
		var file *domains.File

		err := rows.Scan(
			&file.ID,
			&file.Path,
			&file.UpdatedAt,
			&file.Checksum,
		)

		if err != nil {
			return nil, err
		}

		files = append(files, file)

		fmt.Println("from database: ", file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
