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
	SeedID      int
	Seed        Seed `json:"-"`
	AssetID     int
	Asset       Asset `json:"-"`
	N           uint32
	ChangeN     uint32
}

type Wallets []Wallet
