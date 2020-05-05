package models

import (
	"github.com/jinzhu/gorm"
	"strings"
)

type TransactionStatus int

const (
	TransactionStatusUnknown   TransactionStatus = 0
	TransactionStatusNew       TransactionStatus = 1
	TransactionStatusSent      TransactionStatus = 2
	TransactionStatusFailed    TransactionStatus = 3
	TransactionStatusPending   TransactionStatus = 4
	TransactionStatusConfirmed TransactionStatus = 5
)

func (status TransactionStatus) String() string {
	names := [...]string{
		"Unknown",
		"New",
		"Sent",
		"Failed",
		"Confirmed",
	}
	if int(status) >= len(names) {
		return names[0]
	}
	return names[status]
}

type TransactionDirection int

const (
	TransactionDirectionUnknown  TransactionDirection = 0
	TransactionDirectionIncoming TransactionDirection = 1
	TransactionDirectionIutgoing TransactionDirection = 2
)

func (direction TransactionDirection) String() string {
	names := [...]string{
		"Unknown",
		"Incoming",
		"Outgoing",
	}
	if int(direction) >= len(names) {
		return names[0]
	}
	return names[direction]
}

type Transaction struct {
	gorm.Model
	Name      string
	AssetId   []uint `gorm:"-"`
	WalletId  []uint `gorm:"-"`
	OrderId   []uint `gorm:"-"`
	AddressId []uint `gorm:"-"`
	TxHash    string
	Tx        string
	TxData    string
	Direction TransactionDirection
	Status    TransactionStatus
	Asset     []*Asset   `json:"-" gorm:"many2many:transaction_assets;"`
	Wallet    []*Wallet  `json:"-" gorm:"many2many:transaction_wallets;"`
	Order     []*Order   `json:"-" gorm:"many2many:transaction_orders;"`
	Address   []*Address `json:"-" gorm:"many2many:transaction_addresses;"`
}

type Transactions []Transaction

func (t *Transactions) FindByHash(txHash string) *Transaction {
	for i, tx := range []Transaction(*t) {
		if strings.ToLower(tx.TxHash) == txHash {
			return &(*t)[i]
		}
	}
	return nil
}

func (t *Transactions) FindByOrder(orderId uint) (res []Transaction) {
	for _, tx := range []Transaction(*t) {
		for _, oid := range tx.OrderId {
			if oid == orderId {
				res = append(res, tx)
				continue
			}
		}
	}
	return res
}
