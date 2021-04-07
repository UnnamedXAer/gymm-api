package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/unnamedxaer/gymm-api/helpers"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	output, err := json.Marshal(payload)
	if err != nil {
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
	}
	w.Write(output)
}

func getErrOfMalformedInput(u interface{}, excluded []string) string {
	errInfo := make(map[string]string)

	val := reflect.ValueOf(u).Elem()
	for i := 0; i < val.NumField(); i++ {
		fld := val.Type().Field(i)

		if helpers.StrSliceIndexOf(excluded, fld.Name) == -1 {
			jsonFldName := fld.Tag.Get("json")
			if jsonFldName == "" {
				jsonFldName = fld.Name
			}

			errInfo[jsonFldName] = fld.Type.Name()
		}
	}
	/**/
	resErrText := "Malformed payload."
	errInfoTxt, err := json.MarshalIndent(errInfo, "", "  ")
	if err == nil {
		return resErrText + " The payload should look like: \n" + string(errInfoTxt)
	}

	return resErrText
}
