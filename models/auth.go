package models

type CreateUserModel struct {
	Username string `db:"username" json:"username" form:"username" query:"username" validate:"required"`
	Email    string `db:"email" json:"email" form:"email" query:"email" validate:"required,email"`
	Password string `db:"password" json:"password" form:"password" query:"password" validate:"required,min=8,max=12"`
}

type UserTokenModel struct {
	User   UserModel  `json:"user"`
	Tokens TokenModel `json:"tokens"`
}

type LoginUserModel struct {
	ID       int32  `json:"id" form:"id" query:"id" db:"id"`
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type EmailModel struct {
	Email string `json:"email" form:"email" query:"email" db:"email"`
}

type VerifyPasswordModel struct {
	Email            string `json:"email" form:"email" query:"email" validate:"required,email"`
	VerificationCode string `json:"verification_code" form:"verification_code" query:"verification_code"`
	NewPassword      string `json:"new_password" form:"new_password" query:"new_password" validate:"min=8,max=12,eqfield=ConfirmPassword"`
	ConfirmPassword  string `json:"confirm_password" form:"confirm_password" query:"confirm_password" validate:"min=8,max=12"`
}

type VerificationCodeModel struct {
	VerificationCode string `json:"verification_code" form:"verification_code" query:"verification_code" db:"verification_code"`
}
