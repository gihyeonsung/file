package persistence

import (
	"database/sql"
	"strings"

	"github.com/gihyeonsung/file/internal/domain"
)

type SqliteFileRepository struct {
	db *sql.DB
}

func NewSqliteFileRepository(db *sql.DB) *SqliteFileRepository {
	return &SqliteFileRepository{db: db}
}

func (r *SqliteFileRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS files (
			id TEXT PRIMARY KEY NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			path TEXT NOT NULL,
			path_remote TEXT,
			size INTEGER,
			mime_type TEXT
		)
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *SqliteFileRepository) Save(file *domain.File) error {
	query := `
		INSERT OR REPLACE INTO files (id, created_at, updated_at, path, path_remote, size, mime_type)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var pathRemote sql.NullString
	var size sql.NullInt64
	var mimeType sql.NullString

	if file.PathRemote != nil {
		pathRemote.String = *file.PathRemote
		pathRemote.Valid = true
	}

	if file.Size != nil {
		size.Int64 = int64(*file.Size)
		size.Valid = true
	}

	if file.MimeType != nil {
		mimeType.String = *file.MimeType
		mimeType.Valid = true
	}

	_, err := r.db.Exec(query, file.Id, file.CreatedAt, file.UpdatedAt, file.Path, pathRemote, size, mimeType)
	return err
}

func (r *SqliteFileRepository) FindOne(id string) (*domain.File, error) {
	query := `
		SELECT id, created_at, updated_at, path, path_remote, size, mime_type
		FROM files
		WHERE id = ?
	`

	var file domain.File
	var pathRemote sql.NullString
	var size sql.NullInt64
	var mimeType sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&file.Id,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.Path,
		&pathRemote,
		&size,
		&mimeType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if pathRemote.Valid {
		file.PathRemote = &pathRemote.String
	}

	if size.Valid {
		file.Size = &size.Int64
	}

	if mimeType.Valid {
		file.MimeType = &mimeType.String
	}

	return &file, nil
}

func (r *SqliteFileRepository) Find(criteria *domain.FileRepositoryCriteria) (*domain.FileRepositoryResult, error) {
	query := `
		SELECT id, created_at, updated_at, path, path_remote, size, mime_type
		FROM files
		WHERE true
	`
	var args []interface{}

	if criteria != nil && len(criteria.Ids) > 0 {
		placeholders := make([]string, len(criteria.Ids))
		for i := range criteria.Ids {
			placeholders[i] = "?"
		}

		query += " AND id IN (" + strings.Join(placeholders, ",") + ")"
		for _, id := range criteria.Ids {
			args = append(args, id)
		}
	}

	if criteria != nil && len(criteria.Paths) > 0 {
		placeholders := make([]string, len(criteria.Paths))
		for i := range criteria.Paths {
			placeholders[i] = "?"
		}
		query += " AND path IN (" + strings.Join(placeholders, ",") + ")"
		for _, path := range criteria.Paths {
			args = append(args, path)
		}
	}

	if criteria != nil && len(criteria.PathsLike) > 0 {
		placeholders := make([]string, len(criteria.PathsLike))
		for i := range criteria.PathsLike {
			placeholders[i] = "?"
		}
		query += " AND path LIKE (" + strings.Join(placeholders, ",") + ")"
		for _, path := range criteria.PathsLike {
			args = append(args, "%"+path+"%")
		}
	}

	query += " ORDER BY path ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]*domain.File, 0)

	for rows.Next() {
		var file domain.File
		var pathRemote sql.NullString
		var size sql.NullInt64
		var mimeType sql.NullString

		err := rows.Scan(
			&file.Id,
			&file.CreatedAt,
			&file.UpdatedAt,
			&file.Path,
			&pathRemote,
			&size,
			&mimeType,
		)

		if err != nil {
			return nil, err
		}

		if pathRemote.Valid {
			file.PathRemote = &pathRemote.String
		}

		if size.Valid {
			file.Size = &size.Int64
		}

		if mimeType.Valid {
			file.MimeType = &mimeType.String
		}

		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &domain.FileRepositoryResult{Files: files}, nil
}

func (r *SqliteFileRepository) Delete(id string) error {
	query := `DELETE FROM files WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
