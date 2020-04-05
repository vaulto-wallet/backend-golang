package models

import "github.com/jinzhu/gorm"

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
	Asset       Asset `json:"-"`
	N           uint32
	ChangeN     uint32
	Seqno       uint32
}

type Wallets []Wallet
