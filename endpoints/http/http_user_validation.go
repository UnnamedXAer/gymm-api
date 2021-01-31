package http

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func validateUserInput(validate *validator.Validate, u *usecases.UserInput) error {
	errs := validate.Struct(u)

	if errs != nil {
		if _, ok := errs.(*validator.InvalidValidationError); ok {
			fmt.Println(errs)
			return errs
		}

		validErrs := errs.(validator.ValidationErrors)

		formatedErrors := make(map[string]string, len(validErrs))
		var errText, txt string

		for _, err := range validErrs {

			fieldName, found := validation.GetFieldJSONTag(u, err.StructField())
			if found == false {
				fieldName = err.StructField()
			}

			switch err.Tag() {
			// user json tags
			case "pwd":
				txt = fmt.Sprintf("The '%s' is not strong enough", fieldName)
			case "email":
				txt = fmt.Sprintf("The '%s' is not a valid email address", fieldName)
			case "required":
				txt = fmt.Sprintf("The '%s' field value is required and cannot be empty", fieldName)
			case "min":
				var objLengthUnit string
				if err.Kind() == reflect.String {
					objLengthUnit = "characters"
				} else {
					objLengthUnit = "elements"
				}

				txt = fmt.Sprintf("The '%s' has to be at least %s %s long", fieldName, err.Param(), objLengthUnit)
			case "max":
				var objLengthUnit string
				if err.Kind() == reflect.String {
					objLengthUnit = "characters"
				} else {
					objLengthUnit = "elements"
				}

				txt = fmt.Sprintf("The '%s' has to be at max %s %s long", fieldName, err.Param(), objLengthUnit)
			default:
				txt = fmt.Sprintf("The '%s' field failed on the '%s' tag validation", fieldName, err.Tag())
			}

			errText += txt
			formatedErrors[fieldName] = txt
		}

		return validation.NewStructValidError(formatedErrors, errText)
	}

	return nil
}
