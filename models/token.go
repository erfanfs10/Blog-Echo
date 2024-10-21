package models

type TokenModel struct {
	RefreshToken string `json:"refresh"`
	AccessToken  string `json:"access"`
}

type RefreshTokenModel struct {
	RefreshToken string `json:"refresh" form:"refresh" query:"refresh"`
}
