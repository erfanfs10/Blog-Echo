package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/erfanfs10/Blog-Echo/db"
	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/erfanfs10/Blog-Echo/queries"
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
	result, err := db.DB.Exec(queries.Register, createUser.Username, createUser.Email, hashedPassword)
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
	err = db.DB.Get(&user.User, queries.UserMy, lastInsertID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "user not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// generate JWT for new user
	tokens, err := utils.GenerateJWT(user.User.ID, user.User.IsActive)
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
	err := db.DB.Get(&userDB, queries.Login, loginUser.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusUnauthorized, err, "username or password is wrong")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// check if user is active or not
	if !userDB.IsActive {
		err := fmt.Errorf("%s is not active", userDB.Username)
		return utils.HandleError(c, http.StatusForbidden, err, "User is not active")
	}
	// validate user password
	err = utils.CheckHashedPassword(userDB.Password, loginUser.Password)
	if err != nil {
		return utils.HandleError(c, http.StatusUnauthorized, err, "username or password is wrong")
	}
	// update last login to current data
	_, err = db.DB.Exec(queries.UpdateLastLogin, time.Now(), userDB.ID)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// get tokens for users
	tokens, err := utils.GenerateJWT(userDB.ID, userDB.IsActive)
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
	userID, IsActive, err := utils.ValidateRefreshToken(refreshTokenModel.RefreshToken)
	if err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "invalid refresh token")
	}
	// get last is_active status and set it to token
	err = db.DB.Get(&IsActive, queries.GetIsActiveStatus, userID)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// generate tokens again
	tokens, err := utils.GenerateJWT(userID, IsActive)
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
	emailDB := models.EmailModel{}
	err := db.DB.Get(&emailDB, queries.GetEmail, email.Email)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// check if user is active or not
	if !emailDB.IsActive {
		err := fmt.Errorf("%s is not active", email.Email)
		return utils.HandleError(c, http.StatusForbidden, err, "User is not active")
	}
	// Generate verification code
	verificationCode := rand.Intn(999999)
	// Generate text message
	text := fmt.Sprintf("Your verification code is %d", verificationCode)
	// Update generated verification code to user db
	_, err = db.DB.Exec(queries.UpdateVerificationCode, verificationCode, email.Email)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// Send email to the user
	emailJob := utils.EmailJob{To: email.Email, Subject: "Verificatio Code", Text: text}
	utils.EmailChannel <- emailJob
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
	// get the verification code and is active from db
	verificationCodeDB := models.VerificationCodeModel{}
	err := db.DB.Get(&verificationCodeDB, queries.GetVerificationCode, verifyPassword.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusBadRequest, err, "email is not valid")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// check if user is active or not
	if !verificationCodeDB.IsActive {
		err := fmt.Errorf("%s is not active", verifyPassword.Email)
		return utils.HandleError(c, http.StatusForbidden, err, "User is not active")
	}
	// check if verification codes are the same
	if verificationCodeDB.VerificationCode != verifyPassword.VerificationCode {
		return utils.HandleError(c, http.StatusBadRequest, nil, "validation failed")
	}
	// generate new hashPassword
	newHashedPassword, err := utils.HashPassword(verifyPassword.NewPassword)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// update user password
	result, err := db.DB.Exec(queries.UpdatePassword, newHashedPassword, verifyPassword.Email)
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
