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

type AddAccountInput struct {
	AuthCode        uuid.UUID `json:"auth_code"`
	PrivateKey      string    `json:"private_key"`
	Address         string    `json:"address"`
	WalletType      string    `json:"wallet_type"`
	CustomerName    string    `json:"customer_name"`
	MinerID         string    `json:"miner_id"`
	MinerWalletType string    `json:"miner_wallet_type"`
}

type AddAccountOutput struct {
	Id      uuid.UUID `json:"id"`
	Address string    `json:"address"`
}

type AddCustomerInput struct {
	AuthCode     uuid.UUID `json:"auth_code"`
	CustomerName string    `json:"customer_name"`
}

type AddCustomerOutput struct {
	Id           uuid.UUID `json:"id"`
	CustomerName string    `json:"customer_name"`
}

type AddMinerInput struct {
	AuthCode     uuid.UUID `json:"auth_code"`
	CustomerName string    `json:"customer_name"`
	MinerID      string    `json:"miner_id"`
}

type AddMinerOutput struct {
	Id           uuid.UUID `json:"id"`
	MinerID      string    `json:"miner_id"`
	CustomerName string    `json:"customer_name"`
}

type FilecoinCustomer struct {
	Id           uuid.UUID `gorm:"column:id" json:"id"`
	CustomerName string    `gorm:"column:customer_name" json:"customer_name"`
}

type FilecoinMiner struct {
	Id         uuid.UUID `gorm:"column:id" json:"id"`
	CustomerID uuid.UUID `gorm:"customer_id" json:"customer_id"`
	MinerID    string    `gorm:"miner_id" json:"miner_id"`
}

type ListMinersInput = ListReviewersInput
type ListMinersOutput struct {
	Miners []FilecoinMiner `json:"miners"`
}

type ListCustomersInput = ListReviewersInput
type ListCustomersOutput struct {
	Customers []FilecoinCustomer `json:"miners"`
}

type ListWalletTypesInput = ListReviewersInput
type ListWalletTypesOutput struct {
	WalletTypes []string `json:"wallet_types"`
}

type ListMinerWalletTypesInput = ListReviewersInput
type ListMinerWalletTypesOutput struct {
	MinerWalletTypes []string `json:"miner_wallet_types"`
}

type FilecoinAccount struct {
	Id              uuid.UUID `gorm:"column:id" json:"id"`
	Address         string    `gorm:"column:address" json:"address"`
	WalletType      string    `gorm:"column:wallet_type" json:"wallet_type"`
	CustomerID      uuid.UUID `gorm:"column:customer_id" json:"customer_id"`
	MinerID         string    `gorm:"column:miner_id" json:"miner_id"`
	MinerWalletType string    `gorm:"column:miner_wallet_type" json:"miner_wallet_type"`
}

type ListAccountsInput = ListReviewersInput
type ListAccountsOutput struct {
	Accounts []FilecoinAccount `json:"accounts"`
}

type SetBalanceTransferTargetsInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Address  string    `json:"address"`
	Targets  []string  `json:"targets"`
}

type GetBalanceTransferTargetsInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Address  string    `json:"address"`
}

type GetBalanceTransferTargetsOutput struct {
	Id      uuid.UUID `json:"id"`
	Address string    `json:"address"`
	Targets []string  `json:"targets"`
}

type TransferMessage struct {
	From       string `json:"from"`
	FromOwner  string `json:"from_owner"`
	Creator    string `json:"creator"`
	To         string `json:"to"`
	ToOwner    string `json:"to_owner"`
	Amount     string `json:"amount"`
	Cid        string `json:"cid"`
	GasLimit   string `json:"gas_limit"`
	GasFeeCap  string `json:"gas_feecap"`
	GasPremium string `json:"gas_premium"`
	Reviewer   string `json:"reviewer"`
	Reference  string `json:"reference"`
}

type TransferBalanceInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Amount   string    `json:"amount"`
}

type ConfirmBalanceTransferInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Id       uuid.UUID `json:"id"`
}

type ConfirmBalanceTransferOutput struct {
	Message TransferMessage `json:"message"`
}

type RejectBalanceTransferInput = ConfirmBalanceTransferInput

type AccountInfoInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Address  string    `json:"address"`
}

type AccountInfoOutput struct {
	Account      FilecoinAccount `json:"account"`
	CustomerName string          `json:"customer_name"`
	Balance      string          `json:"balance"`
}
