package main

import (
	"encoding/json"
	log "github.com/EntropyPool/entropy-logger"
	mysqlcli "github.com/NpoolFilecoin/filecoin-wallet/mysql"
	"github.com/NpoolFilecoin/filecoin-wallet/types"
	"github.com/NpoolRD/http-daemon"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

type WalletServerConfig struct {
	Port           int                  `json:"port"`
	UserConfigFile string               `json:"user_config_file"`
	MysqlConfig    mysqlcli.MysqlConfig `json:"mysql"`
}

type WalletServer struct {
	config    WalletServerConfig
	authProxy *WalletAuthorizationProxy
	mysqlCli  *mysqlcli.MysqlCli
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

	return types.UserLoginOutput{
		AuthCode: authCode,
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

	err = s.authProxy.AddUser(input.AuthCode, input.User)
	if err != nil {
		return nil, err.Error(), -4
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

	id := uuid.New()
	err = s.mysqlCli.AddBalanceTransferRequest(mysqlcli.BalanceTransferRequest{
		Id:       id,
		Creator:  user.Username,
		Reviewer: reviewer.Username,
		From:     input.From,
		To:       input.To,
		Amount:   input.Amount,
	})
	if err != nil {
		return nil, err.Error(), -10
	}

	return types.RequestBalanceTransferOutput{
		Id: id,
	}, "", 0
}
