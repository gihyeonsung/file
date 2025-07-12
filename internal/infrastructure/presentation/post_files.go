package presentation

import (
	"encoding/json"
	"net/http"
)

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
