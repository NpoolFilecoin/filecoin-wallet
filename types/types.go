package types

import (
	"github.com/google/uuid"
	"time"
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
	Time		 time.Time `gorm:"column:time" json:"time"`
}

type BalanceWithdrawRequest struct {
	Id       uuid.UUID `gorm:"column:id" json:"id"`
	Creator  string    `gorm:"column:creator" json:"creator"`
	Reviewer string    `gorm:"column:reviewer" json:"reviewer"`
	Miner    string    `gorm:"column:owner" json:"miner"`
	Owner    string    `gorm:"column:miner" json:"owner"`
	Amount   float64   `gorm:"column:amount" json:"amount"`
	Status   string    `gorm:"column:status" json:"status"`
	Time		 time.Time `gorm:"column:time" json:"time"`
}

type ListBalanceRequestOutput struct {
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
	Amount     float64 `json:"amount"`
	Cid        string `json:"cid"`
	GasLimit   string `json:"gas_limit"`
	GasFeeCap  string `json:"gas_feecap"`
	GasPremium string `json:"gas_premium"`
	Reviewer   string `json:"reviewer"`
	Reference  string `json:"reference"`
}

type WithdrawMessage struct {
	Miner      string `json:"from"`
	FromOwner  string `json:"from_owner"`
	Creator    string `json:"creator"`
	Owner      string `json:"to"`
	ToOwner    string `json:"to_owner"`
	Amount     float64 `json:"amount"`
	Cid        string `json:"cid"`
	GasLimit   string `json:"gas_limit"`
	GasFeeCap  string `json:"gas_feecap"`
	GasPremium string `json:"gas_premium"`
	Reviewer   string `json:"reviewer"`
	Reference  string `json:"reference"`
}

type ListReviewHistoryInput = ListReviewersInput

type ReviewHistory struct {
	RequestId	 uuid.UUID `gorm:"column:request_id" json:"request_id"`
	From       string `gorm:"column:from" json:"from"`
	FromOwner  string `gorm:"column:from_owner" json:"from_owner"`
	Creator    string `gorm:"column:creator" json:"creator"`
	To         string `gorm:"column:to" json:"to"`
	ToOwner   string `gorm:"column:to_owner" json:"to_owner"`
	Amount     float64 `gorm:"column:amount" json:"amount"`
	Cid        string `gorm:"column:cid" json:"cid"`
	GasLimit   string `gorm:"column:gas_limit" json:"gas_limit"`
	GasFeeCap  string `gorm:"column:gas_feecap" json:"gas_feecap"`
	GasPremium string `gorm:"column:gas_premium" json:"gas_premium"`
	Reviewer   string `gorm:"column:reviewer" json:"reviewer"`
	Status		 string `gorm:"column:status" json:"status"`
	Time			 time.Time `gorm:"column:time" json:"time"`
	Type			 string `gorm:"column:type" json:"type"`
	// Reference  string `gorm:"column:reference" json:"reference"`
}

type ListReviewHistoryOutput struct {
	ReviewListHistorys []ReviewHistory `json:"review_history"`
}

type TransferBalanceInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Amount   string    `json:"amount"`
}

type WithdrawBalanceInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Miner     string    `json:"miner"`
	Owner       string    `json:"owner"`
	Amount   string    `json:"amount"`
}

type ConfirmBalanceTransferInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Id       uuid.UUID `json:"id"`
}

type ConfirmBalanceTransferOutput struct {
	Message TransferMessage `json:"message"`
}

type ConfirmBalanceWithdrawInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Id			 uuid.UUID `json:"id"`
}

type ConfirmBalanceWithdrawOutput struct {
	Message	WithdrawMessage `json:"message"`
}

type RejectBalanceTransferInput = ConfirmBalanceTransferInput

type RejectBalanceWithdrawInput = ConfirmBalanceWithdrawInput

type AccountInfoInput struct {
	AuthCode uuid.UUID `json:"auth_code"`
	Address  string    `json:"address"`
}

type AccountInfoOutput struct {
	Account      FilecoinAccount `json:"account"`
	CustomerName string          `json:"customer_name"`
	Balance      string          `json:"balance"`
}

type MinerInfoInput struct {
	AuthCode		uuid.UUID `json:"auth_code"`
	MinerID			string `json:"miner_id"`
}

type MinerInfoOutput struct {
	Owner					string `json:"owner"`
	CustomerName	string `json:"customer_name"`
	Available			string `json:"available"`
}

type HandlingFeeInput struct {
	AuthCode		uuid.UUID `json:"auth_code"`
	Cid					string `json:"cid"`
	GasFeeCap		string `json:"gas_feecap"`
	GasPremium	string `json:"gas_premium"`
	GasLimit		string `json:"gas_limit"`
}

type HandlingFeeOutput struct {
	Cid			string `json:"cid"`
}

type HandlingInfo struct {
	Nonce				string `json:"nonce"`
	From				string `json:"from"`
	GasLimit		string `json:"gas_limit"`
	GasPremium	string `json:"gas_premium"`
	GasFeecap		string `json:"gas_feecap"`
}

type QueryHandlingStatusInput struct {
	AuthCode		 uuid.UUID `json:"auth_code"`
	Cid					 string `json:"cid"`
}