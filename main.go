package main

import (
	"github.com/erfanfs10/Blog-Echo/db"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/erfanfs10/Blog-Echo/routes"
	"github.com/erfanfs10/Blog-Echo/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func init() {
	db.ConnectDb()
	utils.LoadEnv()
}

func main() {

	e := echo.New()
	e.Validator = utils.CreateCustomValidator()
	defer db.DB.Close()

	e.Use(middlewares.SeparateLogs())
	e.Use(middlewares.CustomLogger())
	// e.Use(middleware.Recover())

	routes.AuthRoutes(e.Group("/auth/"))
	routes.UserRoutes(e.Group("/user/"))
	routes.PostRoutes(e.Group("/post/"))

	e.Logger.Fatal(e.Start(":1323"))

}
