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

func (p *WalletAuthorizationProxy) AddUser(newUser types.WalletUser) error {
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
	log.Infof(log.Fields{}, "p.users is: $v, p.users.Users is: %v, b is: %v, p.config is: %v", p.users, p.users.Users, b, p.config)

	return ioutil.WriteFile(p.config, b, 0666)
}

func (p *WalletAuthorizationProxy) ChangeUser(userBefore types.WalletUser, userAfter types.WalletUser) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var changePosition int
	for k, v := range p.users.Users {
		if userBefore.Username ==  v.Username{
			changePosition = k
		}
		if userAfter.Username != userBefore.Username {
			if userAfter.Username == v.Username {
				return xerrors.Errorf("username %v already exists", userAfter.Username)
			}
		}
	}
	p.users.Users[changePosition].Username = userAfter.Username
	p.users.Users[changePosition].Password = userAfter.Password
	p.users.Users[changePosition].Role = userAfter.Role
	b, _ := json.Marshal(p.users)
	log.Infof(log.Fields{}, "after changing, p.users is: $v, p.users.Users is: %v, b is: %v, p.config is: %v", p.users, p.users.Users, b, p.config)
	return ioutil.WriteFile(p.config, b, 0666)
}

func (p *WalletAuthorizationProxy) DeleteUser(user types.WalletUser) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var deletePosition int
	for k, v := range p.users.Users {
		if user.Username ==  v.Username{
			deletePosition = k
			break
		}
	}
		p.users.Users = append(p.users.Users[:deletePosition], p.users.Users[deletePosition + 1:]...)
	b, _ := json.Marshal(p.users)
	log.Infof(log.Fields{}, "after deleting, p.users is: $v, p.users.Users is: %v, b is: %v, p.config is: %v", p.users, p.users.Users, b, p.config)
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

func (p *WalletAuthorizationProxy) ListReviewers() ([]string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	reviewers := []string{}
	for _, user := range p.users.Users {
		if user.Role == "reviewer" {
			reviewers = append(reviewers, user.Username)
		}
	}

	if len(reviewers) == 0 {
		return nil, xerrors.Errorf("no reviewer available")
	}

	return reviewers, nil
}

func (p *WalletAuthorizationProxy) ListRoles() ([]string, error) {
	return p.users.Roles, nil
}

func (p *WalletAuthorizationProxy) ListUsers() ([]types.WalletUser, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.users.Users, nil
}