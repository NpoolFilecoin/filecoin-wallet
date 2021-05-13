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

type WalletUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AddUserInput struct {
	AuthCode uuid.UUID  `json:"auth_code"`
	User     WalletUser `json:"user"`
}

type RequestBalanceTransferInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Amount   float64   `json:"amount"`
	Reviewer string    `json:"reviewer"`
}

type RequestBalanceTransferOutput struct {
	Id uuid.UUID `json:"id"`
}

type RequestBalanceWithdrawInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Owner    string    `json:"owner"`
	Miner    string    `json:"miner"`
	Amount   float64   `json:"amount"`
	Reviewer string    `json:"reviewer"`
}

type RequestBalanceWithdrawOutput = RequestBalanceTransferOutput

type UserInfoInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
}

type UserInfoOutput struct {
	Username string `json:"username"`
}
