package presentation

import (
	"encoding/json"
	"net/http"
)

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
