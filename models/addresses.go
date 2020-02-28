package models

import (
	"github.com/jinzhu/gorm"
	"math/big"
)

type Address struct {
	gorm.Model
	Address string
	Addon   string
	Wallet  Wallet
	Type    int
	Balance big.Int
	N       int
}
