package routes

import (
	"github.com/erfanfs10/Blog-Echo/handlers"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(g *echo.Group) {
	g.POST("register/", handlers.Register)
	g.POST("login/", handlers.Login)
	g.POST("refresh/", handlers.RefreshToken)
	g.POST("forget_password/", handlers.ForgetPassword)
	g.POST("verify_password/", handlers.VerifyPassword)
}
