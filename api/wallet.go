package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/EntropyPool/entropy-logger"
	"github.com/NpoolDevOps/fbc-devops-peer/api/lotusapi"
	"github.com/NpoolFilecoin/filecoin-wallet/types"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
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

	api.bearerToken = fmt.Sprintf("Bearer %v", string(bearerToken))

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

type Cid struct {
	Cid string `json:"/"`
}

type Message struct {
	From       string
	To         string
	GasFeeCap  string
	GasLimit   uint64
	GasPremium string
}

type nativeMessage struct {
	Message Message
	CID     Cid
}

func (msg *nativeMessage) ToString() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

func (api *WalletAPI) TransferBalance(from, to string, amount string) (types.TransferMessage, error) {
	cmd := exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "send", "--from", from, to, amount)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Errorf(log.Fields{}, "fail to run lotus send: %v | %v", err, string(stderr.Bytes()))
		return types.TransferMessage{}, xerrors.Errorf("%v: %v", err, string(stderr.Bytes()))
	}

	msg := types.TransferMessage{
		Cid: strings.TrimSpace(string(stdout.Bytes())),
	}

	time.Sleep(5 * time.Second)

	// TODO: Get the message with CID, fill the message
	msgs := []nativeMessage{}
	cmd = exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "mpool", "pending", "--local", "--from", from, "--to", to)

	var stdout1, stderr1 bytes.Buffer
	cmd.Stdout = &stdout1
	cmd.Stderr = &stderr1

	err = cmd.Run()
	if err != nil {
		log.Errorf(log.Fields{}, "balance transfer is successful, but fail to get pending message: %v", string(stderr1.Bytes()))
		return msg, nil
	}

	str := strings.Replace(strings.TrimSpace(string(stdout1.Bytes())), "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\\\"", "\"", -1)

	maps := map[string]interface{}{}
	err = json.Unmarshal([]byte(str), &maps)
	if err != nil {
		log.Errorf(log.Fields{}, "balance transfer is successful, but fail to marshal pending message: %v [%v]", err, str)
		return msg, nil
	}

	for k, v := range maps {
		b, _ := json.Marshal(v)
		log.Infof(log.Fields{}, "%v: %v", k, string(b))
	}

	err = json.Unmarshal([]byte(str), &msgs)
	if err != nil {
		log.Errorf(log.Fields{}, "balance transfer is successful, but fail to marshal pending message: %v [%v]", err, str)
		return msg, nil
	}

	found := false
	for _, lmsg := range msgs {
		if lmsg.CID.Cid == msg.Cid {
			log.Infof(log.Fields{}, "msg '%v' is pending [%v]", msg.Cid, lmsg.ToString())
			found = true
			break
		}
	}

	if !found {
		log.Infof(log.Fields{}, "msg '%v' is not pending", msg.Cid)
	}

	return msg, nil
}

func (api *WalletAPI) WithdrawBalance(minerId, owner string, amount string) (types.TransferMessage, error) {
	// TODO: Send, get the message with CID, fill the message
	return types.TransferMessage{}, nil
}
