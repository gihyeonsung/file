package presentation

import (
	"net/http"

	"github.com/gihyeonsung/file/internal/application"
)

type FileController struct {
	fileCreate   *application.FileCreate
	fileDelete   *application.FileDelete
	fileFind     *application.FileFind
	fileDownload *application.FileDownload
	fileUpload   *application.FileUpload
}

func NewFileController(
	mux *http.ServeMux,
	fileCreate *application.FileCreate,
	fileDelete *application.FileDelete,
	fileFind *application.FileFind,
	fileDownload *application.FileDownload,
	fileUpload *application.FileUpload,
) *FileController {
	controller := &FileController{fileCreate: fileCreate, fileDelete: fileDelete, fileFind: fileFind, fileDownload: fileDownload, fileUpload: fileUpload}

	mux.HandleFunc("/api/v1/files", controller.handleFiles)
	mux.HandleFunc("/api/v1/files/{id}", controller.handleFilesId)

	return controller
}

func (c *FileController) handleFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.handleFilesGet(w, r)
	case http.MethodPost:
		c.handleFilesPost(w, r)
	}
}

func (c *FileController) handleFilesId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.handleFilesIdDownload(w, r)
	case http.MethodPost:
		c.handleFilesPostIdUpload(w, r)
	case http.MethodDelete:
		c.handleFilesDeleteId(w, r)
	}
}
