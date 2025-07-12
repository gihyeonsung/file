package presentation

import (
	"encoding/json"
	"net/http"
)

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
