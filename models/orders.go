package models

import "github.com/jinzhu/gorm"

type OrderStatus int

const (
	OrderStatusUnknown    = 0
	OrderStatusNew        = 1
	OrderStatusProcessing = 2
	OrderStatusProcessed  = 3
)

func (a OrderStatus) String() string {
	orderStatusText := [...]string{
		"Unknown",
		"New",
		"Processing",
		"Processed",
	}
	return orderStatusText[a]
}

type Order struct {
	gorm.Model
	Amount        float64
	To            string
	WalletID      int
	SubmittedByID int
	Comment       string
	Status        OrderStatus
	Wallet        Wallet `json:"-"`
	SubmittedBy   User   `json:"-"`
}

type Orders []Order
