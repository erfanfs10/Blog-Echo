package main

import (
	"github.com/erfanfs10/Blog-Echo/configs"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/erfanfs10/Blog-Echo/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func init() {
	configs.ConnectDb()
}

func main() {

	e := echo.New()

	defer configs.DB.Close()

	e.Use(middlewares.SeparateLogs())
	e.Use(middlewares.CustomLogger())
	// e.Use(middleware.Recover())

	routes.UserRoutes(e.Group("/user/"))
	routes.UserRoutesAdmin(e.Group("/admin/"))

	e.Logger.Fatal(e.Start(":1323"))

}
