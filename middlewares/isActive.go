package middlewares

import (
	"errors"
	"net/http"

	"github.com/erfanfs10/Blog-Echo/utils"
	"github.com/labstack/echo/v4"
)

func IsActive() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isActive := c.Get("isActive")
			if isActive.(bool) {
				return next(c)
			}
			err := errors.New("user is not active")
			return utils.HandleError(c, http.StatusForbidden, err, "User is not active")
		}
	}
}
