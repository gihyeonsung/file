package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Path       string
	PathRemote *string // path at remote storage
	Size       *int    // in bytes
}

func NewFile(path string) (*File, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	file := &File{
		Id:         id.String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Path:       path,
		PathRemote: nil,
		Size:       nil,
	}

	return file, nil
}

func (f *File) Upload(pathRemote string, size int) {
	f.PathRemote = &pathRemote
	f.Size = &size
}
