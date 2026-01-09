package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	// marshal data to json
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// write status code
	w.WriteHeader(status)

	// send response
	if _, err := w.Write(response); err != nil {
		return err
	}
	return nil
}

func RespondWithError(w http.ResponseWriter, status int, message string) error {
	payload := map[string]string{"error": message}
	return RespondWithJSON(w, status, payload)
}
