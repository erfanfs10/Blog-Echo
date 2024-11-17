package routes

import (
	"github.com/erfanfs10/Blog-Echo/handlers"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/labstack/echo/v4"
)

func UserRoutes(g *echo.Group) {
	g.GET("list/", handlers.UserList, middlewares.Authenticate(), middlewares.IsActive())
	g.GET("my/", handlers.UserMy, middlewares.Authenticate(), middlewares.IsActive())
	g.GET("search/:username/", handlers.UserSearch, middlewares.Authenticate(), middlewares.IsActive())
}
