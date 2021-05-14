package api

import (
	"github.com/NpoolDevOps/fbc-devops-peer/api/lotusapi"
)

type WalletAPIConfig struct {
	Host string `json:"host"`
	Type string `json:"type"`
}

type WalletAPI struct {
	config WalletAPIConfig
}

func NewWalletAPI(config WalletAPIConfig) *WalletAPI {
	api := &WalletAPI{
		config: config,
	}

	return api
}

func (api *WalletAPI) ImportWallet(privateKey string) (string, error) {
	return lotusapi.ImportWallet(api.config.Host, privateKey)
}
