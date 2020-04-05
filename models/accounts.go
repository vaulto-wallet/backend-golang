package models

import (
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	UserID     int
	User       User   `json:"-"`
	Name       string `json:"name"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}
