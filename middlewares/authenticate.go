package middlewares

import (
	"github.com/erfanfs10/Blog-Echo/utils"
	"github.com/labstack/echo/v4"
)

func Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.ErrUnauthorized
			}
			isValid, userID := utils.ValidateAccessToken(token[7:])
			if !isValid {
				return echo.ErrUnauthorized
			}
			c.Set("userID", userID)
			return next(c)
		}
	}
}
