package response

import (
	"encoding/json"
	"net/http"
)

func StatusOK(w http.ResponseWriter, response interface{}) error {
	return ResponseJson(w, http.StatusOK, response)
}

func BadRequest(w http.ResponseWriter, response interface{}) error {
	return ResponseJson(w, http.StatusBadRequest, response)
}

func Unauthorized(w http.ResponseWriter, response interface{}) error {
	return ResponseJson(w, http.StatusUnauthorized, response)
}

func NotFound(w http.ResponseWriter, response interface{}) error {
	return ResponseJson(w, http.StatusNotFound, response)
}

func InternalServerError(w http.ResponseWriter, response interface{}) error {
	return ResponseJson(w, http.StatusInternalServerError, response)
}

func ResponseJson(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}
