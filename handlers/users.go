package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/erfanfs10/Blog-Echo/db"
	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/erfanfs10/Blog-Echo/utils"
	"github.com/labstack/echo/v4"
)

func ListUsers(c echo.Context) error {
	// get and convert page query param from str to int
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	// get and convert page size query param from str to int
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	// calculate offset
	offset := (page - 1) * pageSize
	users := []models.UserModel{}
	err = db.DB.Select(&users, "SELECT id, username,email FROM users LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return utils.HandleError(c, http.StatusNotFound, err, "no users found")
	}
	res := models.ListUsers{
		Count: len(users),
		Users: users,
	}
	return c.JSON(http.StatusOK, res)
}

func GetUser(c echo.Context) error {
	p := c.Param("id")
	user := models.UserModel{}
	err := db.DB.Get(&user, "SELECT id, username, email FROM users WHERE id=?", p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "user not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	return c.JSON(http.StatusOK, user)
}

func MyUser(c echo.Context) error {
	userID := c.Get("userID")
	user := models.UserModel{}
	err := db.DB.Get(&user, "SELECT id,username,email FROM users WHERE id=?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "user not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	return c.JSON(http.StatusOK, user)
}
