package http

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func validateExerciseInput(validate *validator.Validate, exercise *usecases.ExerciseInput) error {
	errs := validate.Struct(exercise)
	if errs == nil {
		return nil
	}

	validateErrs, ok := errs.(validator.ValidationErrors)
	if !ok {
		return errs
	}

	formattedErrors := make(map[string]string, len(validateErrs))
	var errText, txt string
	for _, err := range validateErrs {
		fieldName := err.Field()
		fn, ok := validation.GetFieldJSONTag(exercise, fieldName)
		if ok {
			fieldName = fn
		}
		switch err.Tag() {
		case "set_unit":
			txt = fmt.Sprintf("The '%s' is incorrect, allowed values: 1 - 'weight', 2 - 'time'. ", fieldName)
		case "required":
			txt = fmt.Sprintf("The '%s' field value is required and cannot be empty. ", fieldName)
		case "min":
			var objLengthUnit string
			if err.Kind() == reflect.String {
				objLengthUnit = "characters"
			} else {
				objLengthUnit = "elements"
			}

			txt = fmt.Sprintf("The '%s' has to be at least %s %s long. ", fieldName, err.Param(), objLengthUnit)
		case "max":
			var objLengthUnit string
			if err.Kind() == reflect.String {
				objLengthUnit = "characters"
			} else {
				objLengthUnit = "elements"
			}

			txt = fmt.Sprintf("The '%s' has to be at max %s %s long. ", fieldName, err.Param(), objLengthUnit)
		case "ex_name_chars":
			txt = fmt.Sprintf("The '%s' is incorrect, allowed are: letters and numbers. ", fieldName)
		case "printascii":
			txt = fmt.Sprintf("The '%s' is incorrect, allowed are: printable characters. ", fieldName)
		default:
			txt = fmt.Sprintf("The '%s' field failed on the '%s' tag validation. ", fieldName, err.Tag())
		}
		errText += txt
		formattedErrors[fieldName] = txt
	}
	return validation.NewStructValidError(formattedErrors, errText)
}
