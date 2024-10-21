package models

import "time"

type UserModel struct {
	ID               int32      `db:"id" json:"id"`
	Username         string     `db:"username" json:"username"`
	Email            string     `db:"email" json:"email"`
	IsAdmin          *bool      `db:"is_admin" json:"is_admin"`
	IsActive         *bool      `db:"is_active" json:"is_active"`
	Created          *time.Time `db:"created" json:"created"`
	Updated          *time.Time `db:"updated" json:"updated"`
	Firstname        *string    `db:"firstname" json:"firstname"`
	Lastname         *string    `db:"lastname" json:"lastname"`
	Gender           *string    `db:"gender" json:"gender"`
	Phonenumber      *string    `db:"phonenumber" json:"phonenumber"`
	LoggedIN         *bool      `db:"logged_in" json:"logged_in"`
	LastLogin        *time.Time `db:"last_login" json:"last_login"`
	Birthday         *time.Time `db:"birthday" json:"birthday"`
	Access           *string    `db:"access" json:"access"`
	Language         *string    `db:"language" json:"language"`
	Theme            *string    `db:"theme" json:"theme"`
	About            *string    `db:"about" json:"about"`
	VerificationCode *string    `db:"verification_code" json:"verification_code"`
	Avatar           *string    `db:"avatar" json:"avatar"`
	Password         string     `db:"password" json:"password"`
}

type CreateUserModel struct {
	Username string `db:"username" json:"username" form:"username"`
	Email    string `db:"email" json:"email" form:"email"`
	IsAdmin  *bool  `db:"is_admin" json:"is_admin" form:"is_admin"`
	Password string `db:"password" json:"password" form:"password"`
}

type LoginUserModel struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}
