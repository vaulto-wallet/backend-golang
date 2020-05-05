package models

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
)

type OrderStatus int

const (
	OrderStatusUnknown            = 0
	OrderStatusNew                = 1
	OrderStatusProcessing         = 2
	OrderStatusPartiallyProcessed = 3
	OrderStatusProcessed          = 4
)

func (a OrderStatus) String() string {
	orderStatusText := [...]string{
		"Unknown",
		"New",
		"Processing",
		"PartiallyProcessed",
		"Processed",
	}
	return orderStatusText[a]
}

type Order struct {
	gorm.Model
	AssetId   uint    `json:"asset_id,omitempty"`
	Symbol    string  `json:"symbol"`
	WalletId  uint    `json:"wallet_id,omitempty"`
	AddressTo string  `json:"address_to,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Comment   string  `json:"comment,omitempty"`
	//	Tx            string         `json:"tx,omitempty"`
	//	TxHash        string         `json:"tx_hash,omitempty"`
	SubmittedById uint           `json:"-"`
	Status        OrderStatus    `json:"status,omitempty"`
	Wallet        Wallet         `json:"wallet,omitempty"`
	SubmittedBy   User           `json:"-"`
	Asset         Asset          `json:"asset,omitempty"`
	Transaction   []*Transaction `json:"-" gorm:"many2many:transaction_orders;"`
}

func (o Order) String() string {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "{ID:%d Wallet:%d Amount:%f %s To:%s Status:%s}", o.ID, o.WalletId, o.Amount, o.Asset.Symbol, o.AddressTo, o.Status)
	return ret.String()
}

type Orders []Order
