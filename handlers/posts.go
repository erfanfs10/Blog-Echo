package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/erfanfs10/Blog-Echo/db"
	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/erfanfs10/Blog-Echo/queries"
	"github.com/erfanfs10/Blog-Echo/utils"
	"github.com/labstack/echo/v4"
)

func PostMy(c echo.Context) error {
	// get userID from context
	userID := c.Get("userID")
	// get posts from db
	posts := []models.Post{}
	err := db.DB.Select(&posts, queries.PostMy, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "user not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// return respose
	return c.JSON(http.StatusOK, posts)
}

func PostList(c echo.Context) error {
	posts := []models.PostWithUser{}
	err := db.DB.Select(&posts, queries.PostList)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "user not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// return respose
	return c.JSON(http.StatusOK, posts)
}

func PostCreate(c echo.Context) error {
	// get userID from context
	userID := c.Get("userID")
	// check if user is active
	if isActiveUser := c.Get("isActive"); !isActiveUser.(bool) {
		return echo.ErrForbidden
	}
	// bind user data to PostCreate model
	postCreate := new(models.PostCreate)
	err := c.Bind(postCreate)
	if err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "bad request")
	}
	// validate request data
	if err := c.Validate(postCreate); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "validation failed")
	}
	// create new post
	result, err := db.DB.Exec(queries.PostCreate, postCreate.Title, postCreate.Body, userID)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// get last inserted id from result
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// get the new post from db
	post := models.Post{}
	err = db.DB.Get(&post, queries.PostGet, lastInsertID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "user not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	return c.JSON(http.StatusCreated, post)
}

func PostUpdate(c echo.Context) error {
	// get userID from context
	userID := c.Get("userID")
	// bind user data to PostUpdate model
	postUpdate := new(models.PostUpdate)
	if err := c.Bind(postUpdate); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "bad request")
	}
	// validate request data
	if err := c.Validate(postUpdate); err != nil {
		return utils.HandleError(c, http.StatusBadRequest, err, "validation failed")
	}
	// get the post from db
	postDB := models.PostUpdate{}
	err := db.DB.Get(&postDB, queries.PostGetUpdate, postUpdate.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "post not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// the creator of the post must be the current user
	if userID != *postDB.UserID {
		return utils.HandleError(c, http.StatusForbidden, nil, "forbidden")
	}
	// update the post
	_, err = db.DB.Exec(queries.PostUpdate, postUpdate.Title, postUpdate.Body, postUpdate.ID)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// get the updated post from db
	updatedPost := models.Post{}
	err = db.DB.Get(&updatedPost, queries.PostGet, postDB.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "post not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// return respose
	return c.JSON(http.StatusOK, updatedPost)
}

func PostDelete(c echo.Context) error {
	// get userID from context
	userID := c.Get("userID")
	// get postID from params
	postID := c.Param("post-id")
	// check valid params
	if postID == "" || postID == "0" {
		return c.JSON(http.StatusBadRequest, "invalid params")
	}
	// get the post from db
	var PostUserID int32
	err := db.DB.Get(&PostUserID, queries.PostGetDelete, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(c, http.StatusNotFound, err, "post not found")
		}
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// the creator of the post must be the current user
	if userID != PostUserID {
		return utils.HandleError(c, http.StatusForbidden, nil, "forbidden")
	}
	// delete the post from db
	_, err = db.DB.Exec(queries.PostDelete, postID)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err, "internal server error")
	}
	// return response
	return c.JSON(http.StatusOK, map[string]string{"message": "post deleted"})
}
