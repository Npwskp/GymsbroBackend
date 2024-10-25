package auth

type LoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterDto struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Age      int    `json:"age" validate:"required,min=1,max=120"`
	Gender   string `json:"gender" validate:"required"`
}

type ReturnToken struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

type GetUserInfo struct {
	Token string `json:"token"`
}
