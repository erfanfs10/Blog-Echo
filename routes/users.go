package routes

import (
	"github.com/erfanfs10/Blog-Echo/handlers"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/labstack/echo/v4"
)

func UserRoutes(g *echo.Group) {
	g.GET("", handlers.ListUsers)
	g.POST("create/", handlers.CreateUser)
	g.GET(":id/", handlers.GetUser)
}

func UserRoutesAdmin(g *echo.Group) {
	g.GET("", handlers.Home, middlewares.Authenticate())
	g.POST("login/", handlers.LoginAdmin)
	g.POST("refresh/", handlers.RefreshToken)
}
