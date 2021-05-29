package main

import (
	"encoding/json"
	"fmt"
	"time"
	log "github.com/EntropyPool/entropy-logger"
	"github.com/NpoolFilecoin/filecoin-wallet/api"
	mysqlcli "github.com/NpoolFilecoin/filecoin-wallet/mysql"
	"github.com/NpoolFilecoin/filecoin-wallet/types"
	"github.com/NpoolRD/http-daemon"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
)

type WalletServerConfig struct {
	Port           int                  `json:"port"`
	UserConfigFile string               `json:"user_config_file"`
	MysqlConfig    mysqlcli.MysqlConfig `json:"mysql"`
	WalletConfig   api.WalletAPIConfig  `json:"wallet"`
}

type WalletServer struct {
	config    WalletServerConfig
	authProxy *WalletAuthorizationProxy
	mysqlCli  *mysqlcli.MysqlCli
	walletAPI *api.WalletAPI
}

func NewWalletServer(cfgFile string) *WalletServer {
	server := &WalletServer{}

	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Errorf(log.Fields{}, "fail to read %v: %v", cfgFile, err)
		return nil
	}

	err = json.Unmarshal(b, &server.config)
	if err != nil {
		log.Errorf(log.Fields{}, "fail to parse %v: %v", cfgFile, err)
		return nil
	}

	authProxy := NewWalletAuthorizationProxy(server.config.UserConfigFile)
	if authProxy == nil {
		log.Errorf(log.Fields{}, "cannot create authorization proxy")
		return nil
	}

	server.authProxy = authProxy

	mysqlCli := mysqlcli.NewMysqlCli(server.config.MysqlConfig)
	if mysqlCli == nil {
		log.Errorf(log.Fields{}, "cannot create mysql client")
		return nil
	}

	server.mysqlCli = mysqlCli

	walletAPI := api.NewWalletAPI(server.config.WalletConfig)
	if walletAPI == nil {
		log.Errorf(log.Fields{}, "cannot create wallet api")
		return nil
	}

	server.walletAPI = walletAPI

	return server
}

func (s *WalletServer) Run() error {
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.UserLoginAPI,
		Handler:  s.UserLoginRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.AddUserAPI,
		Handler:  s.AddUserRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.RequestBalanceTransferAPI,
		Handler:  s.CreateBalanceTransferRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ConfirmBalanceTransferAPI,
		Handler:  s.ConfirmBalanceTransferRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.RejectBalanceTransferAPI,
		Handler:  s.RejectBalanceTransferRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.RequestBalanceWithdrawAPI,
		Handler:  s.CreateBalanceWithdrawRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ConfirmBalanceWithdrawAPI,
		Handler:  s.ConfirmBalanceWithdrawRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.RejectBalanceWithdrawAPI,
		Handler:  s.RejectBalanceWithdrawRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.UserInfoAPI,
		Handler:  s.UserInfoRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListReviewersAPI,
		Handler:  s.ListReviewersRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListRolesAPI,
		Handler:  s.ListRolesRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListBalanceRequestAPI,
		Handler:  s.ListBalanceRequestRequest,
		Method:   "POST",
	})
	//zpl
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListReviewHistoryAPI,
		Handler: 	s.ListReviewHistoryRequest,
		Method: 	"POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.AddCustomerAPI,
		Handler:  s.AddCustomerRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.AddMinerAPI,
		Handler:  s.AddMinerRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListMinersAPI,
		Handler:  s.ListMinersRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListCustomersAPI,
		Handler:  s.ListCustomersRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListWalletTypesAPI,
		Handler:  s.ListWalletTypesRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListMinerWalletTypesAPI,
		Handler:  s.ListMinerWalletTypesRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.AddAccountAPI,
		Handler:  s.AddAccountRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListAccountsAPI,
		Handler:  s.ListAccountsRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.SetBalanceTransferTargetsAPI,
		Handler:  s.SetBalanceTransferTargetsRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.GetBalanceTransferTargetsAPI,
		Handler:  s.GetBalanceTransferTargetsRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.TransferBalanceAPI,
		Handler:  s.TransferBalanceRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.WithdrawBalanceAPI,
		Handler:  s.WithdrawBalanceRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.AccountInfoAPI,
		Handler:  s.AccountInfoRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.MinerInfoAPI,
		Handler:  s.MinerInfoReqeust,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.HandlingFeeAPI,
		Handler:  s.HandlingFeeRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.QueryHandlingStatusAPI,
		Handler:  s.QueryHandlingStatusRequest,
		Method:   "POST",
	})
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.ListUsersAPI,
		Handler:  s.ListUsersRequest,
		Method:   "POST",
	})

	

	httpdaemon.Run(s.config.Port)
	return nil
}

func (s *WalletServer) UserLoginRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.UserLoginInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.Username == "" {
		return nil, "invalid username", -3
	}

	authCode, err := s.authProxy.Login(input.Username, input.Password)
	if err != nil {
		return nil, err.Error(), -4
	}

	user, err := s.authProxy.UserByUsername(input.Username)
	if err != nil {
		return nil, err.Error(), -5
	}
	return types.UserLoginOutput{
		AuthCode: authCode,
		Username: user.Username,
		Role:     user.Role,
	}, "", 0
}

func (s *WalletServer) AddUserRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.AddUserInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.User.Username == "" || input.User.Password == "" {
		return nil, "invalid username or password", -3
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -4
	}

	if user.Role != "admin" {
		return nil, "only admin can add user", -5
	}

	err = s.authProxy.AddUser(input.User)
	if err != nil {
		return nil, err.Error(), -6
	}

	return nil, "", 0
}

func (s *WalletServer) CreateBalanceTransferRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.RequestBalanceTransferInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.From == "" || input.To == "" {
		return nil, "empty from or to is not allowed", -3
	}

	if input.Amount <= 0 {
		return nil, "invalid amount to be transfered", -4
	}

	if input.Reviewer == "" {
		return nil, "reviewer is must", -5
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -6
	}

	if user.Role != "accounter" {
		return nil, "only role 'accounter' can transfer balance", -7
	}

	reviewer, err := s.authProxy.UserByUsername(input.Reviewer)
	if err != nil {
		return nil, err.Error(), -8
	}

	if reviewer.Role != "reviewer" {
		return nil, "reviewer do not have role 'reviewer'", -9
	}

	createTime := time.Unix(time.Now().Unix(), 0)
	id := uuid.New()
	err = s.mysqlCli.AddBalanceTransferRequest(types.BalanceTransferRequest{
		Id:       id,
		Creator:  user.Username,
		Reviewer: reviewer.Username,
		From:     input.From,
		To:       input.To,
		Amount:   input.Amount,
		Time: 		createTime,
	})
	if err != nil {
		return nil, err.Error(), -10
	}

	return types.RequestBalanceTransferOutput{
		Id: id,
	}, "", 0
}

func (s *WalletServer) ConfirmBalanceTransferRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ConfirmBalanceTransferInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "reviewer" {
		return nil, "only role 'reviewer' can confirm transfer balance", -4
	}

	request, err := s.mysqlCli.QueryBalanceTransferRequest(input.Id)
	if err != nil {
		return nil, err.Error(), -5
	}

	if request.Reviewer != user.Username {
		return nil, "not right reviewer", -6
	}

	msg, err := s.walletAPI.TransferBalance(request.From, request.To, fmt.Sprintf("%v", request.Amount))
	if err != nil {
		return nil, err.Error(), -7
	}

	s.mysqlCli.ConfirmBalanceTransferRequest(request)
	confirmTime := time.Unix(time.Now().Unix(), 0)

		// TODO: Query account info then fill to msg
		msg.From = request.From
		msg.To = request.To
		msg.Amount = request.Amount
		msg.Creator = request.Creator
		msg.Reviewer = request.Reviewer

		//
		accountFrom, _ := s.mysqlCli.QueryFilecoinAccount(request.From)
		msg.FromOwner, _ = s.mysqlCli.QueryFilecoinCustomerName(accountFrom.CustomerID)

		accountTo, _ := s.mysqlCli.QueryFilecoinAccount(request.To)
		msg.ToOwner, _ = s.mysqlCli.QueryFilecoinCustomerName(accountTo.CustomerID)
	//
	s.mysqlCli.AddReviewHistory(types.ReviewHistory{
		RequestId: request.Id,
		From: msg.From,
		FromOwner: msg.FromOwner,
		Creator: msg.Creator,
		To: msg.To,
		ToOwner: msg.ToOwner,
		Amount: msg.Amount,
		Cid: msg.Cid,
		GasLimit: msg.GasLimit,
		GasFeeCap: msg.GasFeeCap,
		GasPremium: msg.GasPremium,
		Reviewer: msg.Reviewer,
		Status: "accepted",
		Time: confirmTime,
		Type: "Transfer",
	})
	return types.ConfirmBalanceTransferOutput{
		Message: msg,
	}, "", 0
}

func (s *WalletServer) RejectBalanceTransferRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.RejectBalanceTransferInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "reviewer" {
		return nil, "only role 'reviewer' can reject transfer balance", -4
	}

	request, err := s.mysqlCli.QueryBalanceTransferRequest(input.Id)
	if err != nil {
		return nil, err.Error(), -5
	}

	if request.Reviewer != user.Username {
		return nil, "not right reviewer", -6
	}

	msg, err := s.walletAPI.TransferBalance(request.From, request.To, fmt.Sprintf("%v", request.Amount))
	if err != nil {
		return nil, err.Error(), -7
	}

	s.mysqlCli.RejectBalanceTransferRequest(request)
	rejectTime := time.Unix(time.Now().Unix(), 0)

	msg.From = request.From
	msg.To = request.To
	msg.Amount = request.Amount
	msg.Creator = request.Creator
	msg.Reviewer = request.Reviewer

	//
	accountFrom, _ := s.mysqlCli.QueryFilecoinAccount(request.From)
	msg.FromOwner, _ = s.mysqlCli.QueryFilecoinCustomerName(accountFrom.CustomerID)

	accountTo, _ := s.mysqlCli.QueryFilecoinAccount(request.To)
	msg.ToOwner, _ = s.mysqlCli.QueryFilecoinCustomerName(accountTo.CustomerID)
	//
	s.mysqlCli.AddReviewHistory(types.ReviewHistory{
		RequestId: request.Id,
		From: msg.From,
		FromOwner: msg.FromOwner,
		Creator: msg.Creator,
		To: msg.To,
		ToOwner: msg.ToOwner,
		Amount: msg.Amount,
		Cid: msg.Cid,
		GasLimit: msg.GasLimit,
		GasFeeCap: msg.GasFeeCap,
		GasPremium: msg.GasPremium,
		Reviewer: msg.Reviewer,
		Status: "rejected",
		Time: rejectTime,
		Type: "Transfer",
	})
	return nil, "", 0
}

func (s *WalletServer) CreateBalanceWithdrawRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.RequestBalanceWithdrawInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.Owner == "" || input.Miner == "" {
		return nil, "empty owner or miner is not allowed", -3
	}

	if input.Amount <= 0 {
		return nil, "invalid amount to be transfered", -4
	}

	if input.Reviewer == "" {
		return nil, "reviewer is must", -5
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -6
	}

	if user.Role != "accounter" {
		return nil, "only role 'accounter' can withdraw balance", -7
	}

	reviewer, err := s.authProxy.UserByUsername(input.Reviewer)
	if err != nil {
		return nil, err.Error(), -8
	}

	if reviewer.Role != "reviewer" {
		return nil, "reviewer do not have role 'reviewer'", -9
	}

	createTime := time.Unix(time.Now().Unix(), 0)
	id := uuid.New()
	err = s.mysqlCli.AddBalanceWithdrawRequest(types.BalanceWithdrawRequest{
		Id:       id,
		Creator:  user.Username,
		Reviewer: reviewer.Username,
		Owner:    input.Owner,
		Miner:    input.Miner,
		Amount:   input.Amount,
		Time:			createTime,
	})
	if err != nil {
		return nil, err.Error(), -10
	}

	return types.RequestBalanceWithdrawOutput{
		Id: id,
	}, "", 0
}

func (s *WalletServer) ConfirmBalanceWithdrawRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}
	log.Infof(log.Fields{}, "b is: %v", string(b))

	input := types.ConfirmBalanceWithdrawInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	log.Infof(log.Fields{}, "authcode is: %v, input is : %v, input.authcode is: %v", input.AuthCode, input, input.AuthCode)

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "reviewer" {
		return nil, "only role 'reviewer' can confirm withdraw balance", -4
	}

	request, err := s.mysqlCli.QueryBalanceWithdrawRequest(input.Id)
	if err != nil {
		return nil, err.Error(), -5
	}

	if request.Reviewer != user.Username {
		return nil, "not right reviewer", -6
	}

	msg, err := s.walletAPI.WithdrawBalance(request.Miner, request.Owner, fmt.Sprintf("%v", request.Amount))
	if err != nil {
		return nil, err.Error(), -7
	}

	s.mysqlCli.ConfirmBalanceWithdrawRequest(request)
	confirmTime := time.Unix(time.Now().Unix(), 0)
	customerId, _ := s.mysqlCli.QueryFilecoinAccount(request.Owner)
	msg.FromOwner, _ = s.mysqlCli.QueryFilecoinCustomerName(customerId.CustomerID)
	msg.ToOwner = msg.FromOwner

	msg.Miner = request.Miner
	msg.Owner = request.Owner
	msg.Amount = request.Amount
	msg.Creator = request.Creator
	msg.Reviewer = request.Reviewer

	s.mysqlCli.AddReviewHistory(types.ReviewHistory{
		RequestId: request.Id,
		From: msg.Miner,
		FromOwner: msg.FromOwner,
		Creator: msg.Creator,
		To: msg.Owner,
		ToOwner: msg.ToOwner,
		Amount: msg.Amount,
		Cid: msg.Cid,
		GasLimit: msg.GasLimit,
		GasFeeCap: msg.GasFeeCap,
		GasPremium: msg.GasPremium,
		Reviewer: msg.Reviewer,
		Status: "accepted",
		Time: confirmTime,
		Type: "Withdraw",
	})
	return types.ConfirmBalanceWithdrawOutput{
		Message: msg,
	}, "", 0
}

func (s *WalletServer) RejectBalanceWithdrawRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input :=  types.RejectBalanceWithdrawInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "reviewer"{
		return nil, "only role 'reviewer' can reject withdraw balance", -4
	}

	request, err := s.mysqlCli.QueryBalanceWithdrawRequest(input.Id)
	if err != nil {
		return nil, err.Error(), -5
	}

	if request.Reviewer != user.Username {
		return nil, "not right reviewer", -6
	}

	msg, err := s.walletAPI.WithdrawBalance(request.Miner, request.Owner, fmt.Sprintf("%v", request.Amount))
	if err != nil {
		return nil, err.Error(), -7
	}

	s.mysqlCli.RejectBalanceWithdrawRequest(request)
	rejectTime := time.Unix(time.Now().Unix(), 0)

	customerId, _ := s.mysqlCli.QueryFilecoinAccount(request.Owner)
	msg.FromOwner, _ = s.mysqlCli.QueryFilecoinCustomerName(customerId.CustomerID)
	msg.ToOwner = msg.FromOwner

	msg.Miner = request.Miner
	msg.Owner = request.Owner
	msg.Amount = request.Amount
	msg.Creator = request.Creator
	msg.Reviewer = request.Reviewer
	//
	s.mysqlCli.AddReviewHistory(types.ReviewHistory{
		RequestId: request.Id,
		From: msg.Miner,
		FromOwner: msg.FromOwner,
		Creator: msg.Creator,
		To: msg.Owner,
		ToOwner: msg.ToOwner,
		Amount: msg.Amount,
		Cid: msg.Cid,
		GasLimit: msg.GasLimit,
		GasFeeCap: msg.GasFeeCap,
		GasPremium: msg.GasPremium,
		Reviewer: msg.Reviewer,
		Status: "rejected",
		Time: rejectTime,
		Type: "Withdraw",
	})
	return nil, "", 0
}

func (s *WalletServer) UserInfoRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.UserInfoInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	return types.UserInfoOutput{
		Username: user.Username,
		Role:     user.Role,
	}, "", 0
}

func (s *WalletServer) ListReviewersRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.UserInfoInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	reviewers, err := s.authProxy.ListReviewers()
	if err != nil {
		return nil, err.Error(), -4
	}

	return types.ListReviewersOutput{
		Reviewers: reviewers,
	}, "", 0
}

func (s *WalletServer) ListRolesRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.UserInfoInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	roles, err := s.authProxy.ListRoles()
	if err != nil {
		return nil, err.Error(), -4
	}

	return types.ListRolesOutput{
		Roles: roles,
	}, "", 0
}

func (s *WalletServer) ListBalanceRequestRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListBalanceRequestInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	transferReqs, _ := s.mysqlCli.QueryBalanceTransferRequests()
	withdrawReqs, _ := s.mysqlCli.QueryBalanceWithdrawRequests()

	return types.ListBalanceRequestOutput{
		TransferRequests: transferReqs,
		WithdrawRequests: withdrawReqs,
	}, "", 0
}

//
func (s *WalletServer) ListReviewHistoryRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListReviewHistoryInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	reviewHis, _ := s.mysqlCli.QueryReviewHistory()

	return types.ListReviewHistoryOutput{
		ReviewListHistorys: reviewHis,
	}, "", 0
}

func (s *WalletServer) AddAccountRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.AddAccountInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "admin" {
		return nil, "only admin can add account", -4
	}

	validType := false
	walletTypes := s.walletAPI.WalletTypes()
	for _, walletType := range walletTypes {
		if input.WalletType == walletType {
			validType = true
			break
		}
	}

	if !validType {
		return nil, "invalid wallet type", -5
	}

	if input.WalletType == "miner" {
		validType = false
		minerWalletTypes := s.walletAPI.MinerWalletTypes()
		for _, walletType := range minerWalletTypes {
			if input.MinerWalletType == walletType {
				validType = true
				break
			}
		}

		if !validType {
			return nil, "invalid miner wallet type", -6
		}

		_, err = s.mysqlCli.QueryFilecoinMiner(input.MinerID)
		if err != nil {
			return nil, err.Error(), -7
		}
	}

	customerId, err := s.mysqlCli.QueryFilecoinCustomerId(input.CustomerName)
	if err != nil {
		return nil, err.Error(), -8
	}

	if input.Address == "" {
		return nil, "address should be specified", -9
	}

	addr := input.Address
	havePrivateKey := true
	exists, _ := s.walletAPI.WalletExists(input.Address)
	if !exists {
		if input.PrivateKey == "" {
			addr = input.Address
			havePrivateKey = false
		}else {
			addr, err = s.walletAPI.ImportWallet(input.PrivateKey)
			if err != nil {
				return nil, err.Error(), -9
			}

			if addr == "null" {
				return nil, "key is already imported", -10
			}

			addr = strings.Replace(addr, "\"", "", -1)
			if addr != input.Address {
				return nil, "input address is not what you imported", -11
			}
		}
	} else {
		log.Infof(log.Fields{}, "address '%v' exists, just update database", input.Address)
		addr = input.Address
	}

	id, err := s.mysqlCli.AddFilecoinAccount(types.FilecoinAccount{
		Address:         addr,
		WalletType:      input.WalletType,
		CustomerID:      customerId,
		MinerID:         input.MinerID,
		MinerWalletType: input.MinerWalletType,
		HavePrivateKey:  havePrivateKey,
	})

	return types.AddAccountOutput{
		Address: addr,
		Id:      id,
	}, "", 0
}

func (s *WalletServer) AddCustomerRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.AddCustomerInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "admin" {
		return nil, "only admin can add customer", -4
	}

	id, err := s.mysqlCli.AddFilecoinCustomer(input.CustomerName)
	if err != nil {
		return nil, err.Error(), -5
	}

	return types.AddCustomerOutput{
		Id:           id,
		CustomerName: input.CustomerName,
	}, "", 0
}

func (s *WalletServer) AddMinerRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.AddMinerInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "admin" {
		return nil, "only admin can add customer", -4
	}

	customerId, err := s.mysqlCli.QueryFilecoinCustomerId(input.CustomerName)
	if err != nil {
		return nil, err.Error(), -5
	}

	id, err := s.mysqlCli.AddFilecoinMiner(input.MinerID, customerId)
	if err != nil {
		return nil, err.Error(), -6
	}

	return types.AddMinerOutput{
		Id:           id,
		CustomerName: input.CustomerName,
		MinerID:      input.MinerID,
	}, "", 0
}

func (s *WalletServer) ListMinersRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListMinersInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	miners, err := s.mysqlCli.QueryFilecoinMiners()
	if err != nil {
		return nil, err.Error(), -4
	}

	return types.ListMinersOutput{
		Miners: miners,
	}, "", 0
}

func (s *WalletServer) ListCustomersRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListCustomersInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	customers, err := s.mysqlCli.QueryFilecoinCustomers()
	if err != nil {
		return nil, err.Error(), -4
	}

	return types.ListCustomersOutput{
		Customers: customers,
	}, "", 0
}

func (s *WalletServer) ListWalletTypesRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListWalletTypesInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	walletTypes := s.walletAPI.WalletTypes()

	return types.ListWalletTypesOutput{
		WalletTypes: walletTypes,
	}, "", 0
}

func (s *WalletServer) ListMinerWalletTypesRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListCustomersInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	minerWalletTypes := s.walletAPI.MinerWalletTypes()

	return types.ListMinerWalletTypesOutput{
		MinerWalletTypes: minerWalletTypes,
	}, "", 0
}

func (s *WalletServer) ListAccountsRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListAccountsInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	accounts, err := s.mysqlCli.QueryFilecoinAccounts()
	if err != nil {
		return nil, err.Error(), -4
	}

	return types.ListAccountsOutput{
		Accounts: accounts,
	}, "", 0
}

func (s *WalletServer) SetBalanceTransferTargetsRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.SetBalanceTransferTargetsInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	_, err = s.mysqlCli.QueryFilecoinAccount(input.Address)
	if err != nil {
		return nil, err.Error(), -4
	}

	if len(input.Targets) == 0 {
		return nil, "targets cannot be empty", -5
	}

	for _, addr := range input.Targets {
		_, err := s.mysqlCli.QueryFilecoinAccount(addr)
		if err != nil {
			return nil, err.Error(), -6
		}
	}

	err = s.mysqlCli.SetFilecoinTransferTarget(mysqlcli.FilecoinTransferTarget{
		Address: input.Address,
		Targets: strings.Join(input.Targets, ","),
	})
	if err != nil {
		return nil, err.Error(), -7
	}

	return nil, "", 0
}

func (s *WalletServer) GetBalanceTransferTargetsRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.GetBalanceTransferTargetsInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	_, err = s.mysqlCli.QueryFilecoinAccount(input.Address)
	if err != nil {
		return nil, err.Error(), -4
	}

	target, err := s.mysqlCli.QueryFilecoinTransferTarget(input.Address)
	if err != nil {
		return nil, err.Error(), -5
	}

	return types.GetBalanceTransferTargetsOutput{
		Id:      target.Id,
		Address: input.Address,
		Targets: strings.Split(target.Targets, ","),
	}, "", 0
}

func (s *WalletServer) TransferBalanceRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.TransferBalanceInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	_, err = s.mysqlCli.QueryFilecoinAccount(input.From)
	if err != nil {
		return nil, err.Error(), -4
	}

	_, err = s.mysqlCli.QueryFilecoinAccount(input.To)
	if err != nil {
		return nil, err.Error(), -5
	}

	msg, err := s.walletAPI.TransferBalance(input.From, input.To, input.Amount)
	if err != nil {
		return nil, err.Error(), -6
	}

	return msg, "", 0
}

func (s *WalletServer) WithdrawBalanceRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int){
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.WithdrawBalanceInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	_, err = s.mysqlCli.QueryFilecoinAccount(input.Owner)
	if err != nil {
		return nil, err.Error(), -4
	}

	msg, err := s.walletAPI.WithdrawBalance(input.Miner, input.Owner, input.Amount)
	if err != nil {
		return nil, err.Error(), -5
	}

	return msg, "", 0
}

func (s *WalletServer) AccountInfoRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.AccountInfoInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	account, err := s.mysqlCli.QueryFilecoinAccount(input.Address)
	if err != nil {
		return nil, err.Error(), -4
	}

	customerName, err := s.mysqlCli.QueryFilecoinCustomerName(account.CustomerID)
	if err != nil {
		return nil, err.Error(), -5
	}

	balance, err := s.walletAPI.WalletBalance(input.Address)
	if err != nil {
		return nil, err.Error(), -6
	}

	return types.AccountInfoOutput{
		Account:      account,
		CustomerName: customerName,
		Balance:      balance,
	}, "", 0
}

func (s *WalletServer) MinerInfoReqeust(w http.ResponseWriter, req *http.Request) (interface{}, string, int){
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.MinerInfoInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	accounts, err := s.mysqlCli.QueryFilecoinAccounts()
	if err != nil {
		return nil, err.Error(), -4
	}

	var customerID uuid.UUID
	var address string
	count := 0

	for _, value := range accounts {
		if input.MinerID == value.MinerID && value.MinerWalletType == "owner"{
			count++
			address = value.Address
			customerID = value.CustomerID
		}
	}
	if count > 1 {
		return nil, "miner只能对应一个owner", -5
	}

	customerName, err := s.mysqlCli.QueryFilecoinCustomerName(customerID)
	if err != nil {
		return nil, err.Error(), -6
	}

	available, err := s.walletAPI.MinerAvailable(input.MinerID)
	if err != nil {
		return nil, err.Error(), -7
	}

	return types.MinerInfoOutput{
		Available: available,
		CustomerName: customerName,
		Owner: address,
	}, "", 0
}

func (s *WalletServer) HandlingFeeRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int){
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.HandlingFeeInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	err = s.mysqlCli.QueryReviewHistoryCid(input.Cid)
	if err != nil {
		return nil, err.Error(), -4
	}

	nonce, from, err := s.walletAPI.HandlingCidExists(input.Cid)
	if err != nil {
		return nil, err.Error(), -5
	}

	cid, err := s.walletAPI.HandlingGas(input.Cid, nonce, input.GasLimit, input.GasFeeCap, input.GasPremium, from)
	if err != nil {
		return nil, err.Error(),  -6
	}

	err = s.mysqlCli.UpdateHistoryCid(input.Cid, cid)
	if err != nil {
		return nil, err.Error(), -7
	}
	return types.HandlingFeeOutput{
		Cid: cid,
	}, "", 0
}

func (s *WalletServer) QueryHandlingStatusRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.QueryHandlingStatusInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	_, err = s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	_, _, err = s.walletAPI.HandlingCidExists(input.Cid)
	if err != nil {
		return nil, err.Error(), 0
	}

	return nil, "消息还在堵塞，请继续加手续费", -4
}

func (s *WalletServer) ListUsersRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.ListUsersInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	user, err := s.authProxy.UserByAuthCode(input.AuthCode)
	if err != nil {
		return nil, err.Error(), -3
	}

	if user.Role != "admin" {
		return nil, "Only Admin can See user lists!!!", -4
	}

	userLists, err := s.authProxy.ListUsers()
	
	return types.ListUsersOutput{
		Users: userLists,
	}, "", 0
}