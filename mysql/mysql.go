package walletmysql

import (
	"fmt"
	log "github.com/EntropyPool/entropy-logger"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/xerrors"
)

type MysqlConfig struct {
	Host   string `json:"host"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	DbName string `json:"db"`
}

type MysqlCli struct {
	config MysqlConfig
	url    string
	db     *gorm.DB
}

func NewMysqlCli(config MysqlConfig) *MysqlCli {
	cli := &MysqlCli{
		config: config,
		url: fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
			config.User, config.Passwd, config.Host, config.DbName),
	}

	log.Infof(log.Fields{}, "open mysql db %v", cli.url)
	db, err := gorm.Open("mysql", cli.url)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot open %v: %v", cli.url, err)
		return nil
	}

	log.Infof(log.Fields{}, "successful to create mysql db %v", cli.url)
	db.SingularTable(true)
	cli.db = db

	return cli
}

func (cli *MysqlCli) Delete() {
	cli.db.Close()
}

const (
	RequestCreated  = "created"
	RequestAccepted = "accepted"
	RequestRejected = "rejected"
)

type BalanceTransferRequest struct {
	Id       uuid.UUID `gorm:"column:id"`
	Creator  string    `gorm:"column:creater"`
	Reviewer string    `gorm:"column:reviewer"`
	From     string    `gorm:"column:from"`
	To       string    `gorm:"column:to"`
	Amount   float64   `gorm:"column:amount"`
	Status   string    `gorm:"column:status"`
}

func (cli *MysqlCli) AddBalanceTransferRequest(request BalanceTransferRequest) error {
	request.Status = RequestCreated
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) QueryBalanceTransferRequests() ([]BalanceTransferRequest, error) {
	requests := []BalanceTransferRequest{}
	count := 0

	rc := cli.db.Find(&requests).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any requests")
	}

	return requests, rc.Error
}

func (cli *MysqlCli) ConfirmBalanceTransferRequest(request BalanceTransferRequest) error {
	request.Status = RequestAccepted
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) RejectBalanceTransferRequest(request BalanceTransferRequest) error {
	request.Status = RequestRejected
	rc := cli.db.Save(&request)
	return rc.Error
}

type BalanceWithdrawRequest struct {
	Id       uuid.UUID `gorm:"column:id"`
	Creator  string    `gorm:"column:creater"`
	Reviewer string    `gorm:"column:reviewer"`
	Owner    string    `gorm:"column:owner"`
	Miner    string    `gorm:"column:miner"`
	Amount   float64   `gorm:"column:amount"`
	Status   string    `gorm:"column:status"`
}

func (cli *MysqlCli) AddBalanceWithdrawRequest(request BalanceWithdrawRequest) error {
	request.Status = RequestCreated
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) QueryBalanceWithdrawRequests() ([]BalanceWithdrawRequest, error) {
	requests := []BalanceWithdrawRequest{}
	count := 0

	rc := cli.db.Find(&requests).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any requests")
	}

	return requests, rc.Error
}

func (cli *MysqlCli) ConfirmBalanceWithdrawRequest(request BalanceWithdrawRequest) error {
	request.Status = RequestAccepted
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) RejectBalanceWithdrawRequest(request BalanceWithdrawRequest) error {
	request.Status = RequestRejected
	rc := cli.db.Save(&request)
	return rc.Error
}
