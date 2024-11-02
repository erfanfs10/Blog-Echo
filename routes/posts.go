package routes

import (
	"github.com/erfanfs10/Blog-Echo/handlers"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/labstack/echo/v4"
)

func PostRoutes(g *echo.Group) {
	g.GET("my/", handlers.PostMy, middlewares.Authenticate())
	g.POST("create/", handlers.PostCreate, middlewares.Authenticate())
	g.PATCH("update/", handlers.PostUpdate, middlewares.Authenticate())
	g.DELETE("delete/:post-id/", handlers.PostDelete, middlewares.Authenticate())
	g.GET("list/", handlers.PostList, middlewares.Authenticate())
}
