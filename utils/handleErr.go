package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, code int, err error, message string) error {
	c.Set("err", err.Error())
	if _, ok := err.(validator.ValidationErrors); ok {
		return CustomValidationError(code, message, err, c)
	}
	return c.JSON(code, map[string]string{"message": message})
}
