package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/erfanfs10/Blog-Echo/configs"
	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/erfanfs10/Blog-Echo/utils"
	"github.com/labstack/echo/v4"
)

func ListUsers(c echo.Context) error {
	fmt.Println("in list")
	users := []models.UserModel{}
	err := configs.DB.Select(&users, "SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotFound, "404")
	}
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	fmt.Println("in get")
	p := c.Param("id")
	user := models.UserModel{}
	err := configs.DB.Get(&user, "SELECT * FROM users WHERE id=?", p)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.String(http.StatusNotFound, "user not found")
		}
		return c.String(http.StatusInternalServerError, "something bad happened")
	}
	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	fmt.Println("in create")
	user := new(models.CreateUserModel)
	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}
	return c.JSON(http.StatusOK, user)
}

func Home(c echo.Context) error {
	userID := c.Get("userID")
	user := models.UserModel{}
	err := configs.DB.Get(&user, "SELECT * FROM users WHERE id=?", userID)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.String(http.StatusNotFound, "user not found")
		}
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, user)

}

func LoginAdmin(c echo.Context) error {
	user := new(models.LoginUserModel)
	if err := c.Bind(user); err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusBadRequest, "bad reqeust")
	}
	var userID int
	err := configs.DB.Get(&userID, "SELECT id FROM users WHERE username=? AND password=?", user.Username, user.Password)
	if err != nil {
		c.Set("err", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return c.String(http.StatusNotFound, "user not found")
		}
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	tokens, err := utils.GenerateJWT(int(userID))
	if err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "error when getting JWT")
	}
	return c.JSON(http.StatusOK, tokens)
}

func RefreshToken(c echo.Context) error {
	// Bind user data to RefreshToken struct
	refreshTokenModel := new(models.RefreshTokenModel)
	if err := c.Bind(refreshTokenModel); err != nil {
		return c.String(http.StatusBadRequest, "bad reqeust")
	}
	// Validating refresh token
	isValid, userID := utils.ValidateRefreshToken(refreshTokenModel.RefreshToken)
	if !isValid {
		return c.String(http.StatusForbidden, "JWT is not valid")
	}
	// generate tokens again
	tokens, err := utils.GenerateJWT(int(userID))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "error when getting JWT")
	}
	return c.JSON(http.StatusOK, tokens)
}
