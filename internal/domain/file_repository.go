package domain

type FileRepositoryCriteria struct {
	Ids       []string
	Paths     []string
	PathsLike []string
}

type FileRepositoryResult struct {
	Files []*File
}

type FileRepository interface {
	Save(file *File) error
	FindOne(id string) (*File, error)
	Find(criteria *FileRepositoryCriteria) (*FileRepositoryResult, error)
	Delete(id string) error
}
