package middlewares

import (
	"net/http"

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
			userID, IsActive, err := utils.ValidateAccessToken(token[7:])
			if err != nil {
				return utils.HandleError(c, http.StatusUnauthorized, err, "Unauthorized")
			}
			c.Set("userID", userID)
			c.Set("isActive", IsActive)
			return next(c)
		}
	}
}
