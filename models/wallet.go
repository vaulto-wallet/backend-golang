package models

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
)

type PrivateKeyType int

const (
	PrivateKeyTypeUnknown = 0
	PrivateKeyTypeRoot    = 1
	PrivateKeyTypeSingle  = 2
	PrivateKeyTypeMulti   = 3
)

type Wallet struct {
	gorm.Model
	Name        string
	NetworkType string
	SeedId      uint
	Seed        Seed `json:"-"`
	AssetId     uint
	Asset       Asset
	N           uint32
	ChangeN     uint32
	Seqno       uint32
	Transaction []*Transaction `json:"-" gorm:"many2many:transaction_wallets"`
	Owners      []*User        `json:"-" gorm:"many2many:wallet_owners"`
}

func (o Wallet) String() string {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "{ID:%d \"%s\" %s N:%d Change:%d Seqno:%d}", o.ID, o.Asset.Symbol, o.Name, o.N, o.ChangeN, o.Seqno)
	return ret.String()
}

func (o *Wallet) IsOwner(user *User) bool {
	return o.Seed.OwnerId == user.ID
}

func (o *Wallet) HasPermission(user *User) bool {
	for _, s := range o.Owners {
		if s.ID == user.ID {
			return true
		}
	}
	return false
}

type Wallets []Wallet
