package models

import (
	"github.com/jinzhu/gorm"
)

type FirewallAddress struct {
	gorm.Model
	Address string
	Rule FirewallRule
}

type FirewallRuleSignatures struct {
	Rule FirewallRule
	User User
}

type FirewallRule struct {
	gorm.Model
	PrivateKey PrivateKey
	Amount float64
	SignaturesRequired int
	Period int
}
