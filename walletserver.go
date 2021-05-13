package main

import (
	"encoding/json"
	log "github.com/EntropyPool/entropy-logger"
	"github.com/NpoolFilecoin/filecoin-wallet/types"
	"github.com/NpoolRD/http-daemon"
	"io/ioutil"
	"net/http"
)

type WalletServerConfig struct {
	Port           int    `json:"port"`
	UserConfigFile string `json:"user_config_file"`
}

type WalletServer struct {
	config WalletServerConfig
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

	httpdaemon.Run(s.config.Port)
	return nil
}

func (s *WalletServer) UserLoginRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	return nil, "", 0
}

func (s *WalletServer) AddUserRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	return nil, "", 0
}
