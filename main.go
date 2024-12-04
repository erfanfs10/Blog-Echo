package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/erfanfs10/Blog-Echo/db"
	"github.com/erfanfs10/Blog-Echo/middlewares"
	"github.com/erfanfs10/Blog-Echo/routes"
	"github.com/erfanfs10/Blog-Echo/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	db.ConnectDb()
	db.CreateTables()
	utils.LoadEnv()
	utils.CreateEmailChannel()
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.Validator = utils.CreateCustomValidator()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	defer db.DB.Close()
	defer close(utils.EmailChannel)

	e.Use(middlewares.SeparateLogs())
	e.Use(middlewares.CustomLogger())
	e.Use(middleware.Recover())

	routes.AuthRoutes(e.Group("/auth/"))
	routes.UserRoutes(e.Group("/user/"))
	routes.PostRoutes(e.Group("/post/"))

	e.Logger.Fatal(e.Start(":8000"))

}
