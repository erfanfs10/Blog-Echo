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
	users := []models.UserModel{}
	err := db.DB.Select(&users, "SELECT id, username,email FROM users")
	if err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusNotFound, "404")
	}
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	p := c.Param("id")
	user := models.UserModel{}
	err := db.DB.Get(&user, "SELECT id, username, email FROM users WHERE id=?", p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Set("err", err.Error())
			return c.String(http.StatusNotFound, "user not found")
		}
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "something bad happened")
	}
	return c.JSON(http.StatusOK, user)
}

func MyUser(c echo.Context) error {
	userID := c.Get("userID")
	user := models.UserModel{}
	err := db.DB.Get(&user, "SELECT id,username,email FROM users WHERE id=?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Set("err", err.Error())
			return c.String(http.StatusNotFound, "user not found")
		}
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	createUser := new(models.CreateUserModel)
	err := c.Bind(createUser)
	if err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err = c.Validate(createUser); err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusBadRequest, "your data is not valid")
	}
	hashedPassword, err := utils.HashPassword(createUser.Password)
	if err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "failed generating password")
	}
	result, err := db.DB.Exec("INSERT INTO users(username,email,password) values(?,?,?)", createUser.Username, createUser.Email, hashedPassword)
	if err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusConflict, "can not create user")
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "failed generating password")
	}
	user := models.UserTokenModel{}
	err = db.DB.Get(&user.User, "SELECT id,username,email FROM users WHERE id=?", strconv.Itoa(int(lastInsertID)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Set("err", err.Error())
			return c.String(http.StatusNotFound, "user not found")
		}
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	tokens, err := utils.GenerateJWT(int(user.User.ID))
	if err != nil {
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "error when getting JWT")
	}
	user.Tokens = tokens
	return c.JSON(http.StatusCreated, user)
}

func Home(c echo.Context) error {
	userID := c.Get("userID")
	user := models.UserModel{}
	err := db.DB.Get(&user, "SELECT * FROM users WHERE id=?", userID)
	if err != nil {
		c.Set("err", err.Error())
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
	// TODO check hashed password
	err := db.DB.Get(&userID, "SELECT id FROM users WHERE username=? AND password=?", user.Username, user.Password)
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
		c.Set("err", err.Error())
		return c.String(http.StatusInternalServerError, "error when getting JWT")
	}
	return c.JSON(http.StatusOK, tokens)
}
