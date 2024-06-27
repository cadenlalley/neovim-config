package web

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

// Valid content types.
const ContentTypeMultipartFormData = "multipart/form-data"
const ContentTypeApplicationJSON = "application/json"

func init() {
	validate = validator.New()

	// Use the json tag names in error output instead of struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")
		formTag := fld.Tag.Get("form")

		name := jsonTag
		if formTag != "" {
			name = formTag
		}

		value := strings.SplitN(name, ",", 2)[0]
		if value == "-" {
			return ""
		}
		return value
	})
}

// ValidateRequest checks will bind a request to an interface
// and validate that the request has the required fields.
func ValidateRequest(c echo.Context, contentType string, v interface{}) error {
	if !strings.Contains(c.Request().Header.Get("Content-Type"), contentType) {
		return fmt.Errorf("invalid Content-Type for request, expected: '%s'", contentType)
	}

	if err := c.Bind(v); err != nil {
		return fmt.Errorf("could not bind request '%s'", err.Error())
	}

	if err := validate.Struct(v); err != nil {
		for _, e := range err.(validator.ValidationErrors) {

			switch e.ActualTag() {
			case "required":
				return fmt.Errorf("missing required field '%s'", e.Field())
			case "email":
				return fmt.Errorf("invalid email address '%s' supplied for field '%s'", e.Value(), e.Field())
			default:
				return fmt.Errorf("invalid value '%s' supplied for field '%s'", e.Value(), e.Field())
			}
		}
	}

	return nil
}
