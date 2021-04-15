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
	for _, err := range validateErrs {
		fieldName := err.Field()
		fn, ok := validation.GetFieldJSONTag(exercise, fieldName)
		if ok {
			fieldName = fn
		}
		formattedErrors[fieldName] += getErrorTranslation(&err, fieldName)
	}
	if len(formattedErrors) > 0 {
		errText := concatErrors(formattedErrors)
		return validation.NewStructValidError(formattedErrors, errText)
	}
	return nil
}

func validateExerciseInput4Update(validate *validator.Validate, exercise *usecases.ExerciseInput) error {
	formattedErrors := make(map[string]string)
	v := reflect.ValueOf(exercise).Elem()
	for _, fieldName := range []string{"Name", "Description", "SetUnit"} {
		validateExerciseField(validate, &v, exercise, fieldName, formattedErrors)
	}

	if len(formattedErrors) > 0 {
		errText := concatErrors(formattedErrors)
		return validation.NewStructValidError(formattedErrors, errText)
	}
	return nil
}

func validateExerciseField(
	validate *validator.Validate,
	v *reflect.Value,
	exercise *usecases.ExerciseInput,
	fieldName string,
	formattedErrors map[string]string) {

	strFld, ok := v.Type().FieldByName(fieldName)
	if !ok {
		formattedErrors["more"] = fmt.Sprintf("Could not parse '%s' field.", fieldName)
		return
	}

	strFldVal := v.FieldByName(fieldName)
	if strFldVal.IsZero() {
		return
	}

	tagVal := strFld.Tag.Get("validate")
	var val interface{}

	switch strFldVal.Kind() {
	case reflect.String:
		val = strFldVal.String()
	case reflect.Int8:
		val = strFldVal.Int()
	default:
		_, ok := formattedErrors["more"]
		if ok {
			formattedErrors["more"] += "\n"
		}
		formattedErrors["more"] += fmt.Sprintf("Could not parse '%s' field.", fieldName)
		return
	}

	err := validate.Var(val, tagVal)
	if err != nil {
		fn, ok := validation.GetFieldJSONTag(exercise, fieldName)
		if ok {
			fieldName = fn
		}
		translateValidationErrs(fieldName, err, formattedErrors)
	}
}

func translateValidationErrs(fieldName string, errs error, formattedErrors map[string]string) {
	validateErrs, ok := errs.(validator.ValidationErrors)
	if !ok {
		// @todo: improve error text
		formattedErrors[fieldName] = errs.Error()
	}

	for _, err := range validateErrs {
		formattedErrors[fieldName] += getErrorTranslation(&err, fieldName)
	}
}

func getErrorTranslation(err *validator.FieldError, fieldName string) string {
	switch (*err).Tag() {
	case "set_unit":
		return fmt.Sprintf("The '%s' is incorrect, allowed values: 1 - 'weight', 2 - 'time'. ", fieldName)
	case "required":
		return fmt.Sprintf("The '%s' field value is required and cannot be empty. ", fieldName)
	case "min":
		var objLengthUnit string
		if (*err).Kind() == reflect.String {
			objLengthUnit = "characters"
		} else {
			objLengthUnit = "elements"
		}

		return fmt.Sprintf("The '%s' has to be at least %s %s long. ", fieldName, (*err).Param(), objLengthUnit)
	case "max":
		var objLengthUnit string
		if (*err).Kind() == reflect.String {
			objLengthUnit = "characters"
		} else {
			objLengthUnit = "elements"
		}

		return fmt.Sprintf("The '%s' has to be at max %s %s long. ", fieldName, (*err).Param(), objLengthUnit)
	case "ex_name_chars":
		return fmt.Sprintf("The '%s' is incorrect, allowed are: letters and numbers. ", fieldName)
	case "printascii":
		return fmt.Sprintf("The '%s' is incorrect, allowed are: printable characters. ", fieldName)
	default:
		return fmt.Sprintf("The '%s' field failed on the '%s' tag validation. ", fieldName, (*err).Tag())
	}
}

func concatErrors(formattedErrors map[string]string) string {
	var s string
	for _, v := range formattedErrors {
		s += v
	}

	return s
}
