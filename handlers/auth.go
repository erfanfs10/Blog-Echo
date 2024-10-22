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

func Register(c echo.Context) error {
	// bind request data to CreateUserModel struct
	createUser := new(models.CreateUserModel)
	err := c.Bind(createUser)
	if err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "bad request")
	}
	// validate request data
	if err = c.Validate(createUser); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "invalid data")
	}
	// hashing the password
	hashedPassword, err := utils.HashPassword(createUser.Password)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// create new user
	result, err := db.DB.Exec("INSERT INTO users(username,email,password) values(?,?,?)", createUser.Username, createUser.Email, hashedPassword)
	if err != nil {
		return utils.HandleError(c, http.StatusConflict, err, "user already exists")
	}
	// get last inserted id from result
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// get the new user from db
	user := models.UserTokenModel{}
	err = db.DB.Get(&user.User, "SELECT id,username,email FROM users WHERE id=?", lastInsertID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "user not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// generate JWT for new user
	tokens, err := utils.GenerateJWT(int(user.User.ID))
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// assign tokens and return the response
	user.Tokens = tokens
	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	// bind request data to LoginUserModel struct
	user := new(models.LoginUserModel)
	if err := c.Bind(user); err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// TODO check hashed password
	validateLoginUser := models.ValidateLoginUserModel{}
	err := db.DB.Get(&validateLoginUser, "SELECT id, password FROM users WHERE username=?", user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusUnauthorized, err, "username or password is wrong")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// validate user password
	err = utils.CheckHashedPassword(validateLoginUser.Password, user.Password)
	if err != nil {
		return utils.HandleError(c, http.StatusUnauthorized, err, "username or password is wrong")
	}
	// get tokens for users
	tokens, err := utils.GenerateJWT(int(validateLoginUser.ID))
	if err != nil {
		c.Set("err", err.Error())
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// return response
	return c.JSON(http.StatusOK, tokens)
}

func RefreshToken(c echo.Context) error {
	// Bind user data to RefreshToken struct
	refreshTokenModel := new(models.RefreshTokenModel)
	if err := c.Bind(refreshTokenModel); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "bad request")
	}
	// Validating refresh token
	userID, err := utils.ValidateRefreshToken(refreshTokenModel.RefreshToken)
	if err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "invalid refresh token")
	}
	// generate tokens again
	tokens, err := utils.GenerateJWT(int(userID))
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	return c.JSON(http.StatusOK, tokens)
}
