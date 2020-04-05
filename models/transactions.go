package models

import "github.com/jinzhu/gorm"

type TransactionStatus int

const (
	TransactionStatusUnknown   TransactionStatus = 0
	TransactionStatusNew       TransactionStatus = 1
	TransactionStatusSent      TransactionStatus = 2
	TransactionStatusFailed    TransactionStatus = 3
	TransactionStatusConfirmed TransactionStatus = 4
)

func (status TransactionStatus) String() string {
	names := [...]string{
		"Unknown",
		"New",
		"Sent",
		"Failed",
		"Confirmed",
	}
	return names[status]
}

type Transaction struct {
	gorm.Model
	Name      string
	AssetId   uint
	WalletId  uint
	OrderId   uint
	AddressId []uint
	TxHash    string
	Tx        string
	TxData    string
	Status    TransactionStatus
	Asset     Asset     `json:"-"`
	Wallet    Wallet    `json:"-"`
	Order     Order     `json:"-"`
	Address   []Address `json:"-"`
}

type Transactions []Transaction
