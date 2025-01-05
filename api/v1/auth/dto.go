package auth

import (
	authEnums "github.com/Npwskp/GymsbroBackend/api/v1/auth/enums"
)

type LoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterDto struct {
	Username      string               `json:"username" validate:"required,min=3,max=20"`
	Email         string               `json:"email" validate:"required,email"`
	Password      string               `json:"password" validate:"required"`
	Age           int                  `json:"age" validate:"required,min=1,max=120"`
	Gender        authEnums.GenderType `json:"gender" validate:"required"`
	OAuthProvider string               `json:"oauth_provider,omitempty" default:""`
	OAuthID       string               `json:"oauth_id,omitempty" default:""`
	Picture       string               `json:"picture,omitempty" default:""`
}

type ReturnToken struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

type GetUserInfo struct {
	Token string `json:"token"`
}
