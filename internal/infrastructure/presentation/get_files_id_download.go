package presentation

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *FileController) handleFilesIdDownload(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	f, mimeType, err := c.fileDownload.Execute(id)
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
	defer f.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", *mimeType)
	io.Copy(w, f)
}
