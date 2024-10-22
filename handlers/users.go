package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/erfanfs10/Blog-Echo/db"
	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/erfanfs10/Blog-Echo/utils"
	"github.com/labstack/echo/v4"
)

func ListUsers(c echo.Context) error {
	users := []models.UserModel{}
	err := db.DB.Select(&users, "SELECT id, username,email FROM users")
	if err != nil {
		return utils.HandleError(c, http.StatusNotFound, err, "no users found")
	}
	return c.JSON(http.StatusOK, users)
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
