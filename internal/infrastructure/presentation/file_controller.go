package presentation

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gihyeonsung/file/internal/application"
	"github.com/gihyeonsung/file/internal/domain"
)

type FileController struct {
	fileCreate *application.FileCreate
	fileDelete *application.FileDelete
	fileFind   *application.FileFind
	fileUpload *application.FileUpload
}

func NewFileController(
	mux *http.ServeMux,
	fileCreate *application.FileCreate,
	fileDelete *application.FileDelete,
	fileFind *application.FileFind,
	fileUpload *application.FileUpload,
) *FileController {
	controller := &FileController{fileCreate: fileCreate, fileDelete: fileDelete, fileFind: fileFind, fileUpload: fileUpload}

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
	case http.MethodPost:
		c.handleFilesPostIdUpload(w, r)
	case http.MethodDelete:
		c.handleFilesDeleteId(w, r)
	}
}

func (c *FileController) handleFilesGet(w http.ResponseWriter, r *http.Request) {
	criteria := domain.FileRepositoryCriteria{
		Ids:       strings.Split(r.URL.Query().Get("ids"), ","),
		Paths:     strings.Split(r.URL.Query().Get("paths"), ","),
		PathsLike: strings.Split(r.URL.Query().Get("paths-like"), ","),
	}

	result, err := c.fileFind.Execute(criteria)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Status int    `json:"status"`
			Error  string `json:"error"`
		}{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status int            `json:"status"`
		Data   []*domain.File `json:"data"`
	}{
		Status: http.StatusOK,
		Data:   result.Files,
	})
}

func (c *FileController) handleFilesPost(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Path string `json:"path"`
	}{}

	json.NewDecoder(r.Body).Decode(&body)

	err := c.fileCreate.Execute(body.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Status int    `json:"status"`
			Error  string `json:"error"`
		}{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Status int `json:"status"`
		Data   any `json:"data"`
	}{
		Status: http.StatusCreated,
		Data:   nil,
	})
}

func (c *FileController) handleFilesPostIdUpload(w http.ResponseWriter, r *http.Request) {
	form, err := r.MultipartReader()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Status int    `json:"status"`
			Error  string `json:"error"`
		}{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	id := r.PathValue("id")

	data := []byte{}
	for {
		part, err := form.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(struct {
				Status int    `json:"status"`
				Error  string `json:"error"`
			}{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			})
			return
		}

		partData, err := io.ReadAll(part)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(struct {
				Status int    `json:"status"`
				Error  string `json:"error"`
			}{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			})
			return
		}

		data = append(data, partData...)
	}

	err = c.fileUpload.Execute(id, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Status int    `json:"status"`
			Error  string `json:"error"`
		}{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status int `json:"status"`
		Data   any `json:"data"`
	}{
		Status: http.StatusOK,
		Data:   nil,
	})
}

func (c *FileController) handleFilesDeleteId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := c.fileDelete.Execute(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Status int    `json:"status"`
			Error  string `json:"error"`
		}{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status int `json:"status"`
		Data   any `json:"data"`
	}{
		Status: http.StatusOK,
		Data:   nil,
	})
}
