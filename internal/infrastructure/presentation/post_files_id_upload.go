package presentation

import (
	"encoding/json"
	"io"
	"net/http"
)

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
