package api

import (
	"bytes"
	log "github.com/EntropyPool/entropy-logger"
	"github.com/NpoolDevOps/fbc-devops-peer/api/lotusapi"
	"github.com/NpoolFilecoin/filecoin-wallet/types"
	"io/ioutil"
	"os/exec"
)

type WalletAPIConfig struct {
	Host string `json:"host"`
	Type string `json:"type"`
}

type WalletAPI struct {
	config           WalletAPIConfig
	walletTypes      []string
	minerWalletTypes []string
	bearerToken      string
}

func NewWalletAPI(config WalletAPIConfig) *WalletAPI {
	api := &WalletAPI{
		config:           config,
		walletTypes:      []string{"accounting", "miner"},
		minerWalletTypes: []string{"owner", "worker", "post"},
	}

	bearerToken, err := ioutil.ReadFile("/opt/chain/lotus/token")
	if err != nil {
		log.Errorf(log.Fields{}, "cannot read token file")
		bearerToken = []byte("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.EhlHl0JkXpI-1JYuyPHECkif7TyZEMRnADoBgbd2PBw")
	}

	api.bearerToken = string(bearerToken)

	return api
}

func (api *WalletAPI) ImportWallet(privateKey string) (string, error) {
	return lotusapi.ImportWallet(api.config.Host, privateKey, api.bearerToken)
}

func (api *WalletAPI) WalletTypes() []string {
	return api.walletTypes
}

func (api *WalletAPI) MinerWalletTypes() []string {
	return api.minerWalletTypes
}

func (api *WalletAPI) WalletExists(address string) (bool, error) {
	return lotusapi.WalletExists(api.config.Host, address, api.bearerToken)
}

func (api *WalletAPI) TransferBalance(from, to string, amount string) (types.TransferMessage, error) {
	cmd := exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "send", "--from", from, to, amount)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Errorf(log.Fields{}, "fail to run lotus send: %v | %v", err, string(stderr.Bytes()))
		return types.TransferMessage{}, err
	}

	msg := types.TransferMessage{
		Cid: string(stdout.Bytes()),
	}

	// TODO: Get the message with CID, fill the message

	return msg, nil
}

func (api *WalletAPI) WithdrawBalance(minerId, owner string, amount string) (types.TransferMessage, error) {
	// TODO: Send, get the message with CID, fill the message
	return types.TransferMessage{}, nil
}
