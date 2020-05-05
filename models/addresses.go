package models

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

type Address struct {
	gorm.Model `json:",omitempty"`
	Name       string  `json:"name,omitempty"`
	Address    string  `json:"address,omitempty"`
	PrivateKey string  `json:"private_key,omitempty"`
	WalletID   uint    `json:"wallet_id,omitempty"`
	N          uint32  `json:"n,omitempty"`
	Change     uint32  `json:"change,omitempty"`
	Comment    string  `json:"comment,omitempty"`
	BalanceInt string  `json:"balance_int,omitempty"`
	Balance    float64 `json:"balance,omitempty"`
	Wallet     Wallet  `json:"-"`
	Seqno      uint64  `json:"seqno,omitempty"`
}

type Addresses []Address

func (o Address) String() string {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "{ID:%d Address:%s Wallet:%d Seqno:%d}", o.ID, o.Address, o.WalletID, o.Seqno)
	return ret.String()
}

func (a *Addresses) FindAddress(address string) *Address {
	for _, adr := range []Address(*a) {
		if strings.ToLower(adr.Address) == strings.ToLower(address) {
			return &adr
		}
	}
	return nil
}
