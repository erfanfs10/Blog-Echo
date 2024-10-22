package utils

import "github.com/labstack/echo/v4"

func HandleError(c echo.Context, code int, err error, message string) error {
	c.Set("err", err.Error())
	return c.String(code, message)
}
