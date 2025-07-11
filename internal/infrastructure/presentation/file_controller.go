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

func (c *FileController) handleFilesGet(w http.ResponseWriter, r *http.Request) {
	ids := []string{}
	paths := []string{}
	pathsLike := []string{}

	if r.URL.Query().Get("ids") != "" {
		ids = strings.Split(r.URL.Query().Get("ids"), ",")
	}
	if r.URL.Query().Get("paths") != "" {
		paths = strings.Split(r.URL.Query().Get("paths"), ",")
	}
	if r.URL.Query().Get("paths-like") != "" {
		pathsLike = strings.Split(r.URL.Query().Get("paths-like"), ",")
	}

	criteria := &domain.FileRepositoryCriteria{
		Ids:       ids,
		Paths:     paths,
		PathsLike: pathsLike,
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

func (c *FileController) handleFilesIdDownload(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	bytes, err := c.fileDownload.Execute(id)
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
	w.Write(bytes)
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
