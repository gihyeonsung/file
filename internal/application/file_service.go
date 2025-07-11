package application

type FileService interface {
	Write(path string, data []byte) error
	Read(path string) ([]byte, error)
	Delete(path string) error
}
