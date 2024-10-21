package main

import (
	"github.com/erfanfs10/Blog-Echo/db"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/erfanfs10/Blog-Echo/routes"
	"github.com/erfanfs10/Blog-Echo/utils"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func init() {
	db.ConnectDb()
}

func main() {

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	defer db.DB.Close()

	e.Use(middlewares.SeparateLogs())
	e.Use(middlewares.CustomLogger())
	// e.Use(middleware.Recover())

	routes.UserRoutes(e.Group("/user/"))
	routes.UserRoutesAdmin(e.Group("/admin/"))

	e.Logger.Fatal(e.Start(":1323"))

}
