package models

import "time"

type User struct {
	ID               int32      `db:"id" json:"id"`
	Username         string     `db:"username" json:"username"`
	Email            string     `db:"email" json:"email"`
	IsAdmin          bool       `db:"is_admin" json:"is_admin"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	Created          *time.Time `db:"created" json:"created"`
	Updated          *time.Time `db:"updated" json:"updated"`
	LastLogin        *time.Time `db:"last_login" json:"last_login"`
	VerificationCode *string    `db:"verification_code" json:"verification_code"`
	Avatar           *string    `db:"avatar" json:"avatar"`
}

type UserList struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

type UserSearch struct {
	ID        int32   `db:"id" json:"id"`
	Username  string  `db:"username" json:"username"`
	Email     string  `db:"email" json:"email"`
	Avatar    *string `db:"avatar" json:"avatar"`
	PostCount *int    `db:"post_count" json:"post_count"`
}
