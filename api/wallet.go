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
	"strconv"
	"math/big"
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
	Nonce			 uint64
}

type nativeMessage struct {
	Message Message
	CID     Cid
}

func (msg *nativeMessage) ToString() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

func (api *WalletAPI) WalletBalance(address string) (string, error) {
	cmd := exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "wallet", "balance", address)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Errorf(log.Fields{}, "fail to run lotus wallet balance: %v | %v", err, string(stderr.Bytes()))
		return "", xerrors.Errorf("%v: %v", err, string(stderr.Bytes()))
	}
	strs := strings.Split(strings.TrimSpace(string(stdout.Bytes())), " ")
	log.Infof(log.Fields{}, "stderr is: %v.....stderr.Bytes is: %v ......stdout is: %v", string(stdout.Bytes()), strings.TrimSpace(string(stdout.Bytes())), strs)

	return strs[0], nil
}

func (api *WalletAPI) MinerAvailable(MinerID string) (string, error) {
	cmd := exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "state", "miner-info", MinerID)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Errorf(log.Fields{}, "fail to run lotus miner available: %v | %v", err, string(stderr.Bytes()))
		return "", xerrors.Errorf("%v: %v", err, string(stderr.Bytes()))
	}
	strs := strings.Split(strings.TrimSpace(string(stdout.Bytes())), " ")
	log.Infof(log.Fields{}, "stderr is: %v.....stderr.Bytes is: %v ......stdout is: %v", string(stdout.Bytes()), strings.TrimSpace(string(stdout.Bytes())))
	return strs[2], nil;
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

	time.Sleep(2 * time.Second)

	// TODO: Get the message with CID, fill the message
	msgs := []nativeMessage{}
	// cmd = exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "mpool", "pending", "--local", "--from", from, "--to", to)
	cmd = exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "mpool", "pending", "--local")

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
	str = strings.Replace(str, "}{", "},{", -1)
	str = fmt.Sprintf("[%v]", str)

	err = json.Unmarshal([]byte(str), &msgs)
	if err != nil {
		log.Errorf(log.Fields{}, "balance transfer is successful, but fail to marshal pending message: %v [%v]", err, str)
		return msg, nil
	}

	found := false
	for _, lmsg := range msgs {
		if lmsg.CID.Cid == msg.Cid {
			log.Infof(log.Fields{}, "msg '%v' is pending [%v]", msg.Cid, lmsg.ToString())
			msg.GasFeeCap = lmsg.Message.GasFeeCap
			msg.GasLimit = fmt.Sprintf("%v", lmsg.Message.GasLimit)
			msg.GasPremium = lmsg.Message.GasPremium
			found = true
			break
		}
	}

	if !found {
		log.Infof(log.Fields{}, "msg '%v' is not pending", msg.Cid)
	}

	return msg, nil
}

func (api *WalletAPI) WithdrawBalance(miner, owner string, amount string) (types.WithdrawMessage, error) {
	// TODO: Send, get the message with CID, fill the message

	available, err := api.MinerAvailable(miner)
	if err != nil {
		return types.WithdrawMessage{}, xerrors.Errorf("无法获得Miner %v 的余额", miner)
	}
	availableFloat, _ := strconv.ParseFloat(available, 64)
	availableFloatToInt := int64(availableFloat*1000)
	availableBigInt := big.NewInt(availableFloatToInt)
	amountFloat, _ := strconv.ParseFloat(amount, 64)
	amountFloatToInt := int64(amountFloat*1000)
	amountBigInt := big.NewInt(amountFloatToInt)
	bigNum := big.NewInt(1000000000000000)
	amountBigInt.Mul(amountBigInt, bigNum)
	availableBigInt.Mul(availableBigInt, bigNum)
	if amountBigInt.Cmp(availableBigInt) == 1 {
		return types.WithdrawMessage{}, xerrors.Errorf("提现余额大于可用余额！！！\n 可用余额为：%v", available)
	}
	amount = amountBigInt.String()
	// as := fmt.Sprintf("'{%cAmountRequested%c: %c%v%c}'", '"', '"', '"', amount, '"')
	// as = strings.Replace(as, "\\", "", -1)
	// as := "'{\"AmountRequested\": \""+amount+"\"}'"
	// log.Infof(log.Fields{}, "as is: %v", as)
	// cmd := exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "send", "--from", owner, "--method", "16", "--params-json",  as, miner, "0")

	script := fmt.Sprintf(`#!/bin/bash
/usr/local/bin/lotus --repo /opt/chain/lotus send --from %v --method 16 --params-json '{"AmountRequested": "%v"}' %v 0`,
										owner, amount, miner)
	ioutil.WriteFile("/tmp/lotus-withdraw", []byte(script), 0755)
	defer exec.Command("rm", "/tmp/lotus-withdraw").Output()
	cmd := exec.Command("/tmp/lotus-withdraw")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Errorf(log.Fields{}, "fail to run lotus send: %v | %v", err, string(stderr.Bytes()))
		return types.WithdrawMessage{}, xerrors.Errorf("%v: %v", err, string(stderr.Bytes()))
	}

	msg := types.WithdrawMessage{
		Cid: strings.TrimSpace(string(stdout.Bytes())),
	}

	time.Sleep(2 * time.Second)

	msgs := []nativeMessage{}

	cmd = exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "mpool", "pending", "--local")

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
	str = strings.Replace(str, "}{", "},{", -1)
	str = fmt.Sprintf("[%v]", str)

	err = json.Unmarshal([]byte(str), &msgs)
	if err != nil {
		log.Errorf(log.Fields{}, "balance transfer is successful, but fail to marshal pending message: %v [%v]", err, str)
		return msg, nil
	}

	found := false
	for _, lmsg := range msgs {
		if lmsg.CID.Cid == msg.Cid {
			log.Infof(log.Fields{}, "msg '%v' is pending [%v]", msg.Cid, lmsg.ToString())
			msg.GasFeeCap = lmsg.Message.GasFeeCap
			msg.GasLimit = fmt.Sprintf("%v", lmsg.Message.GasLimit)
			msg.GasPremium = lmsg.Message.GasPremium
			found = true
			break
		}
	}

	if !found {
		log.Infof(log.Fields{}, "msg '%v' is not pending", msg.Cid)
	}

	return msg, nil
}

func (api *WalletAPI) HandlingCidExists(cid string) (string, string, error) {
	msgs := []nativeMessage{}
	cmd := exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "mpool", "pending", "--local")
	var stdout1, stderr1 bytes.Buffer
	cmd.Stdout = &stdout1
	cmd.Stderr = &stderr1
	
	_ = cmd.Run()
	str := strings.Replace(strings.TrimSpace(string(stdout1.Bytes())), "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "}{", "},{", -1)
	str = fmt.Sprintf("[%v]", str)

	_ = json.Unmarshal([]byte(str), &msgs)

	nonce := ""
	from := ""
	found := false
	for _, lmsg := range msgs {
		if lmsg.CID.Cid == cid {
			nonce = fmt.Sprintf("%v",lmsg.Message.Nonce)
			from = lmsg.Message.From
			found = true
		}
	}

	if !found {
		return "", "", xerrors.Errorf("你的申请转账已经送达")
	}

	return nonce, from, nil
}

func (api *WalletAPI) HandlingGas(cid, nonce, gas_limit, gas_feecap, gas_premium, from string) (string, error) {
	cmd := exec.Command("/usr/local/bin/lotus", "--repo", "/opt/chain/lotus", "mpool", "replace", "--gas-feecap", gas_feecap, "--gas-premium", gas_premium, "--gas-limit", gas_limit, from, nonce)
	var stdout1, stderr1 bytes.Buffer
	cmd.Stdout = &stdout1
	cmd.Stderr = &stderr1

	err := cmd.Run()
	if err != nil {
		// log.Errorf(log.Fields{}, "some thing just happen", string(stderr1.Bytes()))
		return "", xerrors.Errorf("%v: %v", err, string(stderr1.Bytes()))
	}
	log.Infof(log.Fields{}, "stdout is: %v, stderr is: %v, stderr.string() is: %v, stdout.string() is: %v", stdout1, stderr1, string(stderr1.Bytes()), string(stdout1.Bytes()))
	log.Infof(log.Fields{}, "new cid is: %v", strings.Split(strings.TrimSpace(string(stdout1.Bytes())), "  "))
	newCid := strings.Split(strings.TrimSpace(string(stdout1.Bytes())), "  ")
	return newCid[1], nil
}
