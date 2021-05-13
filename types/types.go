package types

import (
	"github.com/google/uuid"
)

type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginOutput struct {
	AuthCode uuid.UUID `json:"auth_code"`
}

type AddUserInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	UserName string    `json:"username"`
	Password string    `json:"password"`
}
