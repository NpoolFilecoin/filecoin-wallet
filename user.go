package main

import (
	"encoding/json"
	log "github.com/EntropyPool/entropy-logger"
	"github.com/NpoolFilecoin/filecoin-wallet/types"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"io/ioutil"
	"sync"
)

type WalletUsers struct {
	Roles []string           `json:"roles"`
	Users []types.WalletUser `json:"users"`
}

type WalletAuthorizationProxy struct {
	users    WalletUsers
	config   string
	authCode map[uuid.UUID]types.WalletUser
	mutex    sync.Mutex
}

func NewWalletAuthorizationProxy(userCfg string) *WalletAuthorizationProxy {
	proxy := &WalletAuthorizationProxy{
		config:   userCfg,
		authCode: map[uuid.UUID]types.WalletUser{},
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

func (p *WalletAuthorizationProxy) AddUser(authCode uuid.UUID, newUser types.WalletUser) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	loginedUser, ok := p.authCode[authCode]
	if !ok {
		return xerrors.Errorf("login firstly to create new user")
	}

	if loginedUser.Role != "admin" {
		return xerrors.Errorf("username %v's role %v cannot create new user", loginedUser.Username, loginedUser.Role)
	}

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

func (p *WalletAuthorizationProxy) Login(username string, password string) (uuid.UUID, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for authCode, user := range p.authCode {
		if user.Username == username {
			return authCode, nil
		}
	}

	for _, user := range p.users.Users {
		if user.Username == username && user.Password == password {
			authCode := uuid.New()
			p.authCode[authCode] = user
			return authCode, nil
		}
	}

	return uuid.New(), xerrors.Errorf("username %v not exists or password wrong", username)
}

func (p *WalletAuthorizationProxy) UserByAuthCode(authCode uuid.UUID) (types.WalletUser, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	user, ok := p.authCode[authCode]
	if !ok {
		return types.WalletUser{}, xerrors.Errorf("auth code is not exists")
	}

	return user, nil
}

func (p *WalletAuthorizationProxy) UserByUsername(username string) (types.WalletUser, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, user := range p.users.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return types.WalletUser{}, xerrors.Errorf("cannot find username %v", username)
}
