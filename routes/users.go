package routes

import (
	"github.com/erfanfs10/Blog-Echo/handlers"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/labstack/echo/v4"
)

func UserRoutes(g *echo.Group) {
	g.GET("list/", handlers.ListUsers, middlewares.Authenticate())
	g.POST("my/", handlers.MyUser, middlewares.Authenticate())
	g.GET(":id/", handlers.GetUser, middlewares.Authenticate())
}
