package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func responseWithError(w http.ResponseWriter, code int, err error) {
	log.Println("[responseWithError] code: " + strconv.Itoa(code) + ", err: " + err.Error())
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
