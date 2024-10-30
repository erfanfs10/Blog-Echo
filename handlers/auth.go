package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
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
	if err := c.Validate(createUser); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "validation failed")
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
	// bind user data to LoginUserModel
	loginUser := new(models.LoginUserModel)
	if err := c.Bind(loginUser); err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// get user form db
	userDB := models.LoginUserModel{}
	err := db.DB.Get(&userDB, "SELECT id, username,password FROM users WHERE username=?", loginUser.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusUnauthorized, err, "username or password is wrong")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// validate user password
	err = utils.CheckHashedPassword(userDB.Password, loginUser.Password)
	if err != nil {
		return utils.HandleError(c, http.StatusUnauthorized, err, "username or password is wrong")
	}
	// get tokens for users
	tokens, err := utils.GenerateJWT(int(userDB.ID))
	if err != nil {
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
	tokens, err := utils.GenerateJWT(userID)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	return c.JSON(http.StatusOK, tokens)
}

func ForgetPassword(c echo.Context) error {
	// Bind forgetPasswordEmail to user input
	email := new(models.EmailModel)
	if err := c.Bind(email); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "bad request")
	}
	// Get the user form db
	err := db.DB.Get(&email.Email, "SELECT email FROM users WHERE email=?", email.Email)
	if err != nil {
		return utils.HandleError(c, http.StatusOK, err, "email sent")
	}
	// Generate verification code
	verificationCode := rand.Intn(999999)
	// Generate text message
	text := fmt.Sprintf("Your verification code is %d", verificationCode)
	// Update generated verification code to user db
	_, err = db.DB.Exec("UPDATE users SET verification_code=? WHERE email=?", verificationCode, email.Email)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// Send email to the user
	err = utils.SendEmail(email.Email, "Verificatio Code", text)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "failed to send email")
	}
	// Return the respose
	return c.JSON(http.StatusOK, map[string]string{"message": "sent"})
}

func VerifyPassword(c echo.Context) error {
	// bind verifyPassword to user input
	verifyPassword := new(models.VerifyPasswordModel)
	if err := c.Bind(verifyPassword); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "bad request")
	}
	// validate user input
	if err := c.Validate(verifyPassword); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "validation failed")
	}
	// get the user from db
	verificationCode := new(models.VerificationCodeModel)
	err := db.DB.Get(&verificationCode.VerificationCode, "SELECT verification_code FROM users WHERE email=?", verifyPassword.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusBadRequest, err, "email is not valid")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// check if verification codes are the same
	if verificationCode.VerificationCode != verifyPassword.VerificationCode {
		return utils.HandleError(c, http.StatusBadRequest, errors.New("validation failed"), "validation failed")
	}
	// generate new hashPassword
	newHashedPassword, err := utils.HashPassword(verifyPassword.NewPassword)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// update user password
	result, err := db.DB.Exec("UPDATE users SET password=?, verification_code='' WHERE email=?", newHashedPassword, verifyPassword.Email)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// get RowsAffected from result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// check if rowsAffected is 0 return error
	if rowsAffected == 0 {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// return respose
	return c.JSON(http.StatusOK, map[string]string{"message": "Your password got updated"})
}
