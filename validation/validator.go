package validation

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// StructValidError represents an error for the struct validation
// use `ValidationErrors` to get the map of fields with error text
type StructValidError struct {
	s    string
	errs map[string]string
}

func (err StructValidError) Error() string {
	return err.s
}

// ValidationErrors returns map of fields errors
func (err *StructValidError) ValidationErrors() map[string]string {
	return err.errs
}

// NewStructValidError creates new structValidError with given map of errors and optional error text
func NewStructValidError(errsMap map[string]string, errTxt ...string) *StructValidError {
	err := StructValidError{
		errs: errsMap,
	}
	if errTxt != nil {
		err.s = errTxt[0]
	} else {
		err.s = "Validation failed"
	}
	return &err
}

// New returns new *validator.Validate
func New() *validator.Validate {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("pwdStrength", pwdStrengthValidateFunc)
	validate.RegisterAlias("pwd", "pwdStrength")

	return validate
}

func pwdStrengthValidateFunc(fdl validator.FieldLevel) bool {
	fldValue := fdl.Field().String() // debugging
	fldValueRune := []rune(fldValue)
	pwdCount := len(fldValueRune)
	charCounts := make(map[rune]int, pwdCount)
	var hasLetter, hasDigit bool

	for _, v := range fldValueRune {
		if hasLetter == false && unicode.IsLetter(v) {
			hasLetter = true
		}
		if hasDigit == false && unicode.IsDigit(v) {
			hasDigit = true
		}

		charCounts[v]++
	}

	if hasLetter == false && hasDigit == false {
		return false
	}

	countF := float32(pwdCount)
	ratio := float32(len(charCounts)) / countF
	fmt.Println(fmt.Sprintf("pwd ratio: %f", ratio))
	if ratio > 0.8 {
		return false
	}
	for _, v := range charCounts {
		if float32(v)/countF > 0.35 {
			return false
		}
	}

	return true
}

// PrintValidationErrorInfo prints `err` info to the console
func PrintValidationErrorInfo(err validator.FieldError) {
	fmt.Println("Error: ", err.Error())
	fmt.Println("Namespace: ", err.Namespace()) // can differ when a custom TagNameFunc is registered or
	fmt.Println("Field: ", err.Field())         // by passing alt name to ReportError like below
	fmt.Println("StructNamespace: ", err.StructNamespace())
	fmt.Println("StructField: ", err.StructField())
	fmt.Println("Tag: ", err.Tag())
	fmt.Println("ActualTag: ", err.ActualTag())
	fmt.Println("Kind: ", err.Kind())
	fmt.Println("Type: ", err.Type())
	fmt.Println("Value: ", err.Value())
	fmt.Println("Param: ", err.Param())
	fmt.Println("")
}

// GetFieldJSONTag returns field's `json` tag name
func GetFieldJSONTag(u interface{}, fldName string) (string, bool) {
	v := reflect.ValueOf(u)
	i := reflect.Indirect(v)
	s := i.Type()
	field, found := s.FieldByName(fldName)
	if found == false {
		return "", false
	}
	return field.Tag.Get("json"), true
}
