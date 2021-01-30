package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func responseWithErrorMsg(w http.ResponseWriter, code int, err error) {
	log.Println(fmt.Sprintf("[responseWithErrorMsg] code: %d, err: %#v", code, err))
	responseWithJSON(w, code, map[string]string{"error": err.Error()})
}

func responseWithErrorJSON(w http.ResponseWriter, code int, errObj interface{}) {
	log.Println(fmt.Sprintf("[responseWithErrorJSON] code: %d, err: %#v", code, errObj))
	responseWithJSON(w, code, errObj)
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	output, err := json.Marshal(payload)
	if err != nil {
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
	}
	w.Write(output)
}
