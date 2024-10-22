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
	LastLogin        *time.Time `db:"last_login" json:"last_login"`
	VerificationCode *string    `db:"verification_code" json:"verification_code"`
	Avatar           *string    `db:"avatar" json:"avatar"`
}

type CreateUserModel struct {
	Username string `db:"username" json:"username" form:"username" query:"username" validate:"required"`
	Email    string `db:"email" json:"email" form:"email" query:"email" validate:"required,email"`
	Password string `db:"password" json:"password" form:"password" query:"password" validate:"min=8,max=12"`
}

type UserTokenModel struct {
	User   UserModel  `json:"user"`
	Tokens TokenModel `json:"tokens"`
}

type LoginUserModel struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type ValidateLoginUserModel struct {
	ID       int32  `db:"id" json:"id"`
	Password string `dg:"password" json:"password"`
}
