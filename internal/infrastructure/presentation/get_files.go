package presentation

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gihyeonsung/file/internal/domain"
)

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
