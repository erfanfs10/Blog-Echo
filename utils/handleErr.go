package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, code int, err error, message string) error {
	if err != nil {
		// set the error message into context
		c.Set("err", err.Error())
		// Check if the error is a validation error from `go-playground/validator`
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			return c.JSON(code, map[string]string{"message": message})
		} else {
			return CustomValidationError(code, message, err, c)
		}
	}
	return c.JSON(code, map[string]string{"message": message})
}
