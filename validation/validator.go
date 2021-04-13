package validation

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/unnamedxaer/gymm-api/entities"
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

	// @todo: use snake case
	validate.RegisterValidation("pwdStrength", pwdStrengthValidateFunc)
	validate.RegisterAlias("pwd", "pwdStrength")

	validate.RegisterValidation("set_unit", setUnitValidateFunc)

	return validate
}

func setUnitValidateFunc(fldLev validator.FieldLevel) bool {
	fld := fldLev.Field()
	return validateSetUnit(fld)
}

func validateSetUnit(fld reflect.Value) bool {
	switch fld.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fldValue := fld.Int()
		if fldValue == int64(entities.Weight) || fldValue == int64(entities.Time) {
			return true
		}
	}

	return false
}

func pwdStrengthValidateFunc(fdl validator.FieldLevel) bool {
	fldValue := fdl.Field().String()
	return validatePassword(fldValue)
}

func validatePassword(pwd string) bool {
	pwdCount := len(pwd)
	if pwdCount < 6 {
		return false
	}
	pwdRunes := []rune(pwd)
	charCounts := make(map[rune]int, pwdCount)
	var diffCharTypesCnt int
	var hasLetter, hasDigit, hasSpecial, hasSpecialExtended bool

	for _, v := range pwdRunes {
		if v <= 0x2f || (v >= 0x3a && v <= 0x40) || (v >= 0x5b && v <= 0x60) {
			if hasSpecial == false {
				hasSpecial = true
				diffCharTypesCnt++
			}
		} else if (v >= 0x61 && v <= 0x7a) || (v >= 0x41 && v <= 0x5a) {
			if hasLetter == false {
				hasLetter = true
				diffCharTypesCnt++
			}
		} else if v >= 0x30 && v <= 0x39 {
			if hasDigit == false {
				hasDigit = true
				diffCharTypesCnt++
			}
		} else if v >= 0x7b {
			if hasSpecialExtended == false {
				hasSpecialExtended = true
				diffCharTypesCnt++
			}
		}

		charCounts[v]++
		if diffCharTypesCnt > 1 && len(charCounts) >= 5 {
			return true
		}
	}
	return false
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
