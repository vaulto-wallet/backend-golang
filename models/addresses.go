package models

import (
	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model `json:",omitempty"`
	Name       string  `json:"name,omitempty"`
	Address    string  `json:"address,omitempty"`
	PrivateKey string  `json:"private_key,omitempty"`
	WalletID   int     `json:"wallet_id,omitempty"`
	N          uint32  `json:"n,omitempty"`
	Change     uint32  `json:"change,omitempty"`
	Comment    string  `json:"comment,omitempty"`
	BalanceInt string  `json:"balance_int,omitempty"`
	Balance    float64 `json:"balance,omitempty"`
	Wallet     Wallet  `json:"-"`
	Seqno      uint64  `json:"seqno,omitempty"`
}

type Addresses []Address
