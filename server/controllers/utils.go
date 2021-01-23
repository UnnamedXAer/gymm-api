package controllers

import (
	"encoding/json"
	"net/http"
)

func responseWithError(w http.ResponseWriter, code int, err error) {
	responseWithJSON(w, code, map[string]string{"error": err.Error()})
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	output, err := json.Marshal(payload)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
	}
	w.Write(output)
}
