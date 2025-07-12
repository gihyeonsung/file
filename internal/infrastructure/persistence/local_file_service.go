package persistence

import (
	"io"
	"os"
	"path/filepath"
)

type LocalFileService struct {
	pathBase string
}

func NewLocalFileService(pathBase string) *LocalFileService {
	return &LocalFileService{pathBase: pathBase}
}

func (s *LocalFileService) Write(path string, r io.Reader) (int, error) {
	path = filepath.Join(s.pathBase, path)

	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	size, err := io.Copy(file, r)
	if err != nil {
		return 0, err
	}

	return int(size), nil
}

func (s *LocalFileService) Read(path string) (io.ReadCloser, error) {
	path = filepath.Join(s.pathBase, path)

	return os.Open(path)
}

func (s *LocalFileService) Delete(path string) error {
	path = filepath.Join(s.pathBase, path)

	return os.Remove(path)
}
