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
	Username string    `json:"username"`
	Role     string    `json:"role"`
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
	Role     string `json:"role"`
}

type ListReviewersInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
}

type ListReviewersOutput struct {
	Reviewers []string `json:"reviewers"`
}

type ListRolesInput = ListReviewersInput
type ListRolesOutput struct {
	Roles []string `json:"roles"`
}

type ListBalanceRequestInput = ListReviewersInput

type BalanceTransferRequest struct {
	Id       uuid.UUID `gorm:"column:id" json:"id"`
	Creator  string    `gorm:"column:creator" json:"creator"`
	Reviewer string    `gorm:"column:reviewer" json:"reviewer"`
	From     string    `gorm:"column:from" json:"from"`
	To       string    `gorm:"column:to" json:"to"`
	Amount   float64   `gorm:"column:amount" json:"amount"`
	Status   string    `gorm:"column:status" json:"status"`
}

type BalanceWithdrawRequest struct {
	Id       uuid.UUID `gorm:"column:id" json:"id"`
	Creator  string    `gorm:"column:creator" json:"creator"`
	Reviewer string    `gorm:"column:reviewer" json:"reviewer"`
	Owner    string    `gorm:"column:owner" json:"owner"`
	Miner    string    `gorm:"column:miner" json:"miner"`
	Amount   float64   `gorm:"column:amount" json:"amount"`
	Status   string    `gorm:"column:status" json:"status"`
}

type ListBanalceRequestOutput struct {
	TransferRequests []BalanceTransferRequest `json:"transfer_requests"`
	WithdrawRequests []BalanceWithdrawRequest `json:"withdraw_requests"`
}
