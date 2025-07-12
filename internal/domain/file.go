package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Path       string    `json:"path"`
	PathRemote *string   `json:"pathRemote"` // path at remote storage
	Size       *int      `json:"size"`       // in bytes
	MimeType   *string   `json:"mimeType"`
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
		MimeType:   nil,
	}

	return file, nil
}

func (f *File) Upload(pathRemote string, size int, mimeType string) {
	f.PathRemote = &pathRemote
	f.Size = &size
	f.MimeType = &mimeType
}
