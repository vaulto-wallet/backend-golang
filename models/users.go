package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AccountId uint      `json:"-"`
	Account   Account   `json:"account"`
	Wallets   []*Wallet `json:"-" gorm:"many2many:wallet_owners"`
}

type Users []*User
