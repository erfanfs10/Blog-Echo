package models

import "time"

type Post struct {
	ID      int32      `db:"id" json:"id"`
	Title   *string    `db:"title" json:"title"`
	Body    *string    `db:"body" json:"body"`
	UserID  *int32     `db:"user_id" json:"user_id"`
	Created *time.Time `db:"created" json:"created"`
	Updated *time.Time `db:"updated" json:"updated"`
}

type PostCreate struct {
	Title *string `db:"title" json:"title" form:"title" query:"title" validate:"required"`
	Body  *string `db:"body" json:"body" form:"body" query:"body" validate:"required"`
}

type PostUpdate struct {
	ID     int32   `db:"id" json:"id" form:"id" query:"id" validate:"required"`
	UserID *int32  `db:"user_id" json:"user_id" `
	Title  *string `db:"title" json:"title" form:"title" query:"title" validate:"required"`
	Body   *string `db:"body" json:"body" form:"body" query:"body" validate:"required"`
}

type PostUser struct {
	ID       *int32  `db:"user_id" json:"id"`
	Username *string `db:"username" json:"username"`
	Email    *string `db:"email" json:"email"`
	Avatar   *string `db:"avatar" json:"avatar"`
}

type PostWithUser struct {
	ID        *int32     `db:"id" json:"id"`
	Title     *string    `db:"title" json:"title"`
	Body      *string    `db:"body" json:"body"`
	Created   *time.Time `db:"created" json:"created"`
	Updated   *time.Time `db:"updated" json:"updated"`
	*PostUser `json:"user"`
}
