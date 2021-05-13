package main

import (
	"encoding/json"
	log "github.com/EntropyPool/entropy-logger"
	"golang.org/x/xerrors"
	"io/ioutil"
	"sync"
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
	mutex  sync.Mutex
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

func (p *WalletAuthorizationProxy) AddUser(newUser WalletUser) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	validRole := false
	for _, role := range p.users.Roles {
		if newUser.Role == role {
			validRole = true
			break
		}
	}

	if !validRole {
		return xerrors.Errorf("role %v is not in %v", newUser.Role, p.users.Roles)
	}

	for _, user := range p.users.Users {
		if user.Username == newUser.Username {
			return xerrors.Errorf("username %v already exists", user.Username)
		}
	}

	p.users.Users = append(p.users.Users, newUser)
	b, _ := json.Marshal(p.users)

	return ioutil.WriteFile(p.config, b, 0666)
}
