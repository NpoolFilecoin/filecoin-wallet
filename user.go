package main

import (
	"encoding/json"
	log "github.com/EntropyPool/entropy-logger"
	"io/ioutil"
)

type WalletUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type WalletUsers struct {
	Roles []string     `json:"roles"`
	Users []WalletUser `json:"users"`
}

type WalletAuthorizationProxy struct {
	users  WalletUsers
	config string
}

func NewWalletAuthorizationProxy(userCfg string) *WalletAuthorizationProxy {
	proxy := &WalletAuthorizationProxy{
		config: userCfg,
	}

	b, err := ioutil.ReadFile(userCfg)
	if err != nil {
		log.Errorf(log.Fields{}, "fail to read %v: %v", userCfg, err)
		return nil
	}

	err = json.Unmarshal(b, &proxy.users)
	if err != nil {
		log.Errorf(log.Fields{}, "fail to parse %v: %v", userCfg, err)
		return nil
	}

	return proxy
}
