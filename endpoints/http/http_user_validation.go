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
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := errs.(*validator.InvalidValidationError); ok {
			fmt.Println(errs)
			return errs
		}

		validErrs := errs.(validator.ValidationErrors)

		formatedErrors := make(map[string]string, len(validErrs))
		var errText, txt string

		for _, err := range validErrs {

			// fmt.Println("Namespace: ", err.Namespace()) // can differ when a custom TagNameFunc is registered or
			// fmt.Println("Field: ", err.Field())         // by passing alt name to ReportError like below
			// fmt.Println("StructNamespace: ", err.StructNamespace())
			// fmt.Println("StructField: ", err.StructField())
			// fmt.Println("Tag: ", err.Tag())
			// fmt.Println("ActualTag: ", err.ActualTag())
			// fmt.Println("Kind: ", err.Kind())
			// fmt.Println("Type: ", err.Type())
			// fmt.Println("Value: ", err.Value())
			// fmt.Println("Param: ", err.Param())
			// fmt.Println()

			switch err.Tag() {
			// user json tags
			case "pwd":
				txt = fmt.Sprint("The 'password' is not strong enough")
			case "email":
				txt = fmt.Sprintf("The '%s' is not a valid email address", err.StructField())
			case "required":
				txt = fmt.Sprintf("The '%s' field value is required and cannot be empty", err.StructField())
			case "min":
				var objLengthUnit string
				if err.Kind() == reflect.String {
					objLengthUnit = " characters"
				} else {
					objLengthUnit = " elements"
				}

				txt = fmt.Sprintf("The '%s' has to be at least %d%s long", err.StructField(), 5, objLengthUnit)
			case "max":
				var objLengthUnit string
				if err.Kind() == reflect.String {
					objLengthUnit = " characters"
				} else {
					objLengthUnit = " elements"
				}

				txt = fmt.Sprintf("The '%s' has to be at max %d%s long", err.StructField(), 5, objLengthUnit)
			default:
				txt = fmt.Sprintf("The '%s' field failed on the '%s' tag validation", err.StructField(), err.Tag())
			}

			errText += txt
			formatedErrors[err.StructField()] = txt
		}

		return validation.NewStructValidError(formatedErrors, errText)
	}
	return nil
}
