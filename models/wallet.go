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
	AssetID     int
	Seed        Seed  `json:"-"`
	Asset       Asset `json:"-"`
}

type Wallets []Wallet
