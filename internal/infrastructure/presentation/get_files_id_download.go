package presentation

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (c *FileController) handleFilesIdDownload(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	f, mimeType, fileSize, err := c.fileDownload.Execute(id)
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

	w.Header().Set("Content-Type", *mimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(*fileSize, 10))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}
