package api

import (
	"github.com/NpoolDevOps/fbc-devops-peer/api/lotusapi"
)

type WalletAPIConfig struct {
	Host string `json:"host"`
	Type string `json:"type"`
}

type WalletAPI struct {
	config           WalletAPIConfig
	walletTypes      []string
	minerWalletTypes []string
}

func NewWalletAPI(config WalletAPIConfig) *WalletAPI {
	api := &WalletAPI{
		config:           config,
		walletTypes:      []string{"accounting", "miner"},
		minerWalletTypes: []string{"owner", "worker", "post"},
	}

	return api
}

func (api *WalletAPI) ImportWallet(privateKey string, bearerToken string) (string, error) {
	return lotusapi.ImportWallet(api.config.Host, privateKey, bearerToken)
}

func (api *WalletAPI) WalletTypes() []string {
	return api.walletTypes
}

func (api *WalletAPI) MinerWalletTypes() []string {
	return api.minerWalletTypes
}
