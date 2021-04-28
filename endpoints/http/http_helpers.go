package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/helpers"
)

// sends error response with given code and error's text as a response
func responseWithError(w http.ResponseWriter, code int, err error) {
	responseWithJSON(w, code, map[string]string{"error": err.Error()})
}

// sends error response with code 500 - Internal Server Error
func responseWithInternalError(w http.ResponseWriter) {
	responseWithJSON(w, http.StatusInternalServerError,
		map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
}

// sends error response with code 401 - Unauthorized, and optional v as a error message otherwise status text will be used
func responseWithUnauthorized(w http.ResponseWriter, v ...interface{}) {
	var errTxt string
	if len(v) > 0 {
		errTxt = fmt.Sprintf("%v", v[0])
	} else {
		errTxt = http.StatusText(http.StatusUnauthorized)
	}
	responseWithJSON(w, http.StatusUnauthorized,
		map[string]string{"error": errTxt})
}

// sends error response with given code and message as a response
func responseWithErrorTxt(w http.ResponseWriter, code int, errTxt string) {
	responseWithJSON(w, code, map[string]string{"error": errTxt})
}

// sends http response with given payload
func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	output, err := json.Marshal(payload)
	if err != nil {
		responseWithErrorTxt(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}
	w.Write(output)
}

// returns error message about malformed payload, message includes info about
// correct payload structure
//
// Parameter excludedFields can contain struct fields that should be not included
// in massage.
func getErrOfMalformedInput(v interface{}, excludedFields []string) string {
	errInfo := make(map[string]string)

	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {

		fld := val.Type().Field(i)

		if helpers.StrSliceIndexOf(excludedFields, fld.Name) == -1 {
			jsonFldName := fld.Tag.Get("json")
			if jsonFldName == "" {
				jsonFldName = fld.Name
			}

			// errInfo[jsonFldName] = fld.Type.Kind().String()

			errInfo[jsonFldName] = getSimpleType(fld.Type)
		}
	}
	/**/
	resErrText := "Malformed payload."
	errInfoTxt, err := json.Marshal(errInfo)
	if err == nil {
		return resErrText + " The payload should look like: \n" + string(errInfoTxt)
	}

	return resErrText
}

func getSimpleType(fld reflect.Type) string {
	switch fld.Kind() {
	case reflect.Invalid:
		return "-invalid-"
	case reflect.Ptr:
		return getSimpleType(fld.Elem())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return " positive int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Complex64, reflect.Complex128:
		return "complex"
	case reflect.String:
		return "string"
	case reflect.Map, reflect.Struct:
		return "object"
	case reflect.Array, reflect.Slice:
		return "array"
	case reflect.Bool:
		return "bool"
	}
	return "???"
}

// unmarTypeErrorFormat creates new  error from json.UnmarshalTypeError
// for use as http response
func unmarTypeErrorFormat(unmarshalTypeErr *json.UnmarshalTypeError) error {
	return fmt.Errorf("%q has to be of type %s, got %s",
		unmarshalTypeErr.Field, getSimpleType(unmarshalTypeErr.Type), unmarshalTypeErr.Value)
}

// check if error is one of json errors and create new one
// with more usefull message for client, otherwise return false and original error
func formatParseErrors(err error) (bool, error) {
	syntaxErr, ok := err.(*json.SyntaxError)
	if ok {
		return true, syntaxErr
	}

	invalidUnmarshalErr, ok := err.(*json.InvalidUnmarshalError)
	if ok {
		return true, invalidUnmarshalErr
	}

	unmarshalTypeErr, ok := err.(*json.UnmarshalTypeError)
	if ok {
		return true, unmarTypeErrorFormat(unmarshalTypeErr)
	}
	return false, err
}

// return new error saying user has not permissions for requested data
//
// dataName param can be used to pass name of requested data
// eg.: formatUnauthorizedError("exercise") == "unauthorized: ... requested exercise"
func formatUnauthorizedError(dataName ...string) error {
	name := "data"
	if len(dataName) > 0 {
		name = dataName[0]
	}
	return fmt.Errorf("unauthorized: you do not have permissons to requested %s", name)
}

func logDebugError(l *zerolog.Logger, req *http.Request, err error) {
	l.Debug().Msgf("[%s %s]: error: %v", req.Method, req.RequestURI, err)
}
func logDebug(l *zerolog.Logger, req *http.Request, v interface{}) {
	l.Debug().Msgf("[%s %s]: %v", req.Method, req.RequestURI, v)
}
