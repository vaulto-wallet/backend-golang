package models

import "github.com/jinzhu/gorm"

type Confirmation struct {
	gorm.Model
	OrderId uint
	Order   Order
	UserId  uint
	User    User
}

type Confirmations []Confirmation
