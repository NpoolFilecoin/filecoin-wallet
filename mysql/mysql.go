package walletmysql

import (
	"fmt"

	log "github.com/EntropyPool/entropy-logger"
	"github.com/NpoolFilecoin/filecoin-wallet/types"
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

func (cli *MysqlCli) AddBalanceTransferRequest(request types.BalanceTransferRequest) error {
	request.Status = RequestCreated
	rc := cli.db.Save(&request)
	log.Infof(log.Fields{}, "add balance transfer request %v", request) //
	return rc.Error
}

func (cli *MysqlCli) QueryBalanceTransferRequest(id uuid.UUID) (types.BalanceTransferRequest, error) {
	request := types.BalanceTransferRequest{}
	count := 0

	rc := cli.db.Where("id = ?", id).Find(&request).Count(&count)
	if count == 0 {
		return request, xerrors.Errorf("cannot find request")
	}

	return request, rc.Error
}

func (cli *MysqlCli) QueryBalanceWithdrawRequest(id uuid.UUID) (types.BalanceWithdrawRequest, error) {
	request := types.BalanceWithdrawRequest{}
	count := 0

	rc := cli.db.Where("id = ?", id).Find(&request).Count(&count)
	if count == 0 {
		return request, xerrors.Errorf("cannot find request")
	}

	return request, rc.Error
}

func (cli *MysqlCli) QueryBalanceTransferRequests() ([]types.BalanceTransferRequest, error) {
	requests := []types.BalanceTransferRequest{}
	count := 0

	rc := cli.db.Find(&requests).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any requests")
	}

	return requests, rc.Error
}

func (cli *MysqlCli) QueryReviewHistoryPagination(limit, offset uint64, reviewType string) ([]types.ReviewHistory, error) {
	requests := []types.ReviewHistory{}
	rc := cli.db.Where("type = ?", reviewType).Order("time desc").Limit(limit).Offset(offset).Find(&requests)
	return requests, rc.Error
}

//
func (cli *MysqlCli) QueryReviewHistory() ([]types.ReviewHistory, error) {
	requests := []types.ReviewHistory{}
	count := 0
	rc := cli.db.Find(&requests).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any requests")
	}

	return requests, rc.Error
}

func (cli *MysqlCli) QueryReviewHistoryCid(cid string) error {
	history := types.ReviewHistory{}
	count := 0

	cli.db.Where("cid = ?", cid).Find(&history).Count(&count)
	if 0 >= count {
		return xerrors.Errorf("cannot find cid")
	}
	return nil
}

func (cli *MysqlCli) ConfirmBalanceTransferRequest(request types.BalanceTransferRequest) error {
	request.Status = RequestAccepted
	rc := cli.db.Save(&request)
	log.Infof(log.Fields{}, "confirm transfer request %v", request)
	return rc.Error
}

//
func (cli *MysqlCli) AddReviewHistory(request types.ReviewHistory) error {
	rc := cli.db.Save(&request)
	log.Infof(log.Fields{}, "review history sends to sql %v", request)
	return rc.Error
}

func (cli *MysqlCli) RejectBalanceTransferRequest(request types.BalanceTransferRequest) error {
	request.Status = RequestRejected
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) AddBalanceWithdrawRequest(request types.BalanceWithdrawRequest) error {
	request.Status = RequestCreated
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) QueryBalanceWithdrawRequests() ([]types.BalanceWithdrawRequest, error) {
	requests := []types.BalanceWithdrawRequest{}
	count := 0

	rc := cli.db.Find(&requests).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any requests")
	}

	return requests, rc.Error
}

func (cli *MysqlCli) ConfirmBalanceWithdrawRequest(request types.BalanceWithdrawRequest) error {
	request.Status = RequestAccepted
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) RejectBalanceWithdrawRequest(request types.BalanceWithdrawRequest) error {
	request.Status = RequestRejected
	rc := cli.db.Save(&request)
	return rc.Error
}

func (cli *MysqlCli) AddFilecoinCustomer(customerName string) (uuid.UUID, error) {
	customer := types.FilecoinCustomer{}
	count := 0

	cli.db.Where("customer_name = ?", customerName).Find(&customer).Count(&count)
	if 0 < count {
		return customer.Id, nil
	}

	rc := cli.db.Save(&types.FilecoinCustomer{
		Id:           uuid.New(),
		CustomerName: customerName,
	})
	if rc.Error != nil {
		return uuid.New(), rc.Error
	}

	rc = cli.db.Where("customer_name = ?", customerName).Find(&customer).Count(&count)
	if rc.Error != nil {
		return uuid.New(), rc.Error
	}
	if count == 0 {
		return uuid.New(), xerrors.Errorf("cannot find customer after insterted")
	}
	return customer.Id, nil
}

func (cli *MysqlCli) QueryFilecoinCustomerId(customerName string) (uuid.UUID, error) {
	customer := types.FilecoinCustomer{}
	count := 0

	cli.db.Where("customer_name = ?", customerName).Find(&customer).Count(&count)
	if 0 < count {
		return customer.Id, nil
	}

	return uuid.New(), xerrors.Errorf("cannot find customer '%v'", customerName)
}

func (cli *MysqlCli) QueryFilecoinCustomerName(id uuid.UUID) (string, error) {
	customer := types.FilecoinCustomer{}
	count := 0

	cli.db.Where("id = ?", id).Find(&customer).Count(&count)
	if 0 < count {
		return customer.CustomerName, nil
	}

	return "", xerrors.Errorf("cannot find customer '%v'", id)
}

func (cli *MysqlCli) UpdateHistoryCid(cid, newCid string) error {
	request := types.ReviewHistory{}

	rc := cli.db.Model(&request).Where("cid = ?", cid).UpdateColumn(map[string]interface{}{
		"cid": newCid,
	}).Error

	return rc
}

func (cli *MysqlCli) QueryFilecoinCustomers() ([]types.FilecoinCustomer, error) {
	customers := []types.FilecoinCustomer{}
	count := 0

	_ = cli.db.Find(&customers).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any customer")
	}

	return customers, nil
}

func (cli *MysqlCli) AddFilecoinMiner(minerId string, customerId uuid.UUID) (uuid.UUID, error) {
	miner := types.FilecoinMiner{}
	count := 0

	cli.db.Where("miner_id = ?", minerId).Find(&miner).Count(&count)
	if 0 < count {
		return miner.Id, nil
	}

	rc := cli.db.Save(&types.FilecoinMiner{
		Id:         uuid.New(),
		MinerID:    minerId,
		CustomerID: customerId,
	})
	if rc.Error != nil {
		return uuid.New(), rc.Error
	}

	rc = cli.db.Where("miner_id = ?", minerId).Find(&miner).Count(&count)
	if rc.Error != nil {
		return uuid.New(), rc.Error
	}
	if count == 0 {
		return uuid.New(), xerrors.Errorf("cannot find miner after insterted")
	}
	return miner.Id, nil
}

func (cli *MysqlCli) QueryFilecoinMiner(minerId string) (uuid.UUID, error) {
	miner := types.FilecoinMiner{}
	count := 0

	cli.db.Where("miner_id = ?", minerId).Find(&miner).Count(&count)
	if 0 < count {
		return miner.Id, nil
	}

	return uuid.New(), xerrors.Errorf("cannot find miner '%v'", minerId)
}

func (cli *MysqlCli) QueryFilecoinMiners() ([]types.FilecoinMiner, error) {
	miners := []types.FilecoinMiner{}
	count := 0

	rc := cli.db.Find(&miners).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any customer")
	}

	return miners, rc.Error
}

func (cli *MysqlCli) AddFilecoinAccount(account types.FilecoinAccount) (uuid.UUID, error) {
	accounts := []types.FilecoinAccount{}
	count := 0

	cli.db.Where("address = ?", account.Address).Find(&accounts).Count(&count)
	if 0 < count {
		return accounts[0].Id, nil
	}

	account.Id = uuid.New()
	rc := cli.db.Save(&account)
	if rc.Error != nil {
		return uuid.New(), rc.Error
	}

	return account.Id, nil
}

func (cli *MysqlCli) QueryFilecoinAccount(address string) (types.FilecoinAccount, error) {
	account := types.FilecoinAccount{}
	count := 0

	rc := cli.db.Where("address = ?", address).Find(&account).Count(&count)
	if count == 0 {
		return account, xerrors.Errorf("no filecoin account '%v' available", address)
	}

	return account, rc.Error
}

func (cli *MysqlCli) QueryFilecoinAccounts() ([]types.FilecoinAccount, error) {
	accounts := []types.FilecoinAccount{}
	count := 0

	rc := cli.db.Find(&accounts).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("no filecoin account available")
	}

	return accounts, rc.Error
}

type FilecoinTransferTarget struct {
	Id      uuid.UUID `gorm:"column:id"`
	Address string    `gorm:"column:address"`
	Targets string    `gorm:"column:target_addresses"`
}

func (cli *MysqlCli) SetFilecoinTransferTarget(target FilecoinTransferTarget) error {
	t := FilecoinTransferTarget{}
	count := 0

	cli.db.Where("address = ?", target.Address).Find(&t).Count(&count)
	if 0 < count {
		target.Id = t.Id
	} else {
		target.Id = uuid.New()
	}

	rc := cli.db.Save(&target)
	return rc.Error
}

func (cli *MysqlCli) QueryFilecoinTransferTarget(address string) (*FilecoinTransferTarget, error) {
	target := FilecoinTransferTarget{}
	count := 0

	cli.db.Where("address = ?", address).Find(&target).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("no address targets '%v' available", address)
	}

	return &target, nil
}
