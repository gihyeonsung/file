package application

import "io"

type FileService interface {
	Write(path string, r io.Reader) (int64, error)
	Read(path string) (io.ReadCloser, error)
	Delete(path string) error
}
