package models

import (
	h "../helpers"
	"bytes"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type OrderStatus int

const (
	OrderStatusUnknown            = OrderStatus(0)
	OrderStatusNew                = OrderStatus(1)
	OrderStatusConfirmed          = OrderStatus(2)
	OrderStatusProcessing         = OrderStatus(3)
	OrderStatusPartiallyProcessed = OrderStatus(4)
	OrderStatusProcessed          = OrderStatus(5)
)

func (a OrderStatus) String() string {
	orderStatusText := [...]string{
		"Unknown",
		"New",
		"Confirmed",
		"Processing",
		"PartiallyProcessed",
		"Processed",
	}
	return orderStatusText[a]
}

type OrderDestination struct {
	OrderId   uint    `json:"-"`
	AddressTo string  `json:"address_to,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
}

func (o *OrderDestination) String() string {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "{To:%s Amount:%f}", o.AddressTo, o.Amount)
	return ret.String()

}

type OrderDestinations []OrderDestination

type Order struct {
	gorm.Model
	AssetId       uint                `json:"asset_id,omitempty"`
	Symbol        string              `json:"symbol"`
	WalletId      uint                `json:"wallet_id,omitempty"`
	Comment       string              `json:"comment,omitempty"`
	Destinations  []*OrderDestination `json:"destinations,omitempty"`
	SubmittedById uint                `json:"-"`
	Status        OrderStatus         `json:"status,omitempty"`
	Wallet        Wallet              `json:"wallet,omitempty"`
	SubmittedBy   User                `json:"-"`
	Asset         Asset               `json:"asset,omitempty"`
	Transactions  []*Transaction      `json:"-" gorm:"many2many:transaction_orders;"`
	Confirmations []*Confirmation     `json:"-"`
	RuleId        uint                `json:"rule_id,omitempty"`
	Rule          FirewallRule        `json:"-"`
}

func (o *Order) Sum() (ret float64) {
	ret = 0
	for _, d := range o.Destinations {
		ret += d.Amount
	}
	return
}

func (o *Order) FindRule(db *gorm.DB, confirmation Confirmation) (*FirewallRule, bool, error) {
	targetAddress := FirewallAddressTypeInternal
	for _, d := range o.Destinations {
		var address Address
		db.First(&address, "address = ?", d.AddressTo)
		if address.ID == 0 {
			targetAddress = FirewallAddressTypeExternal
			break
		}
	}
	var dbWallet Wallet
	db.First(&dbWallet, o.WalletId)
	if dbWallet.ID == 0 {
		return nil, false, errors.New("Wallet not found")
	}

	db.Preload("FirewallRules").Preload("Coowners").Preload("Auditors").First(&dbWallet, o.WalletId)

	var groupsAffected [][]uint

	var activeRule *FirewallRule
	activeRule = nil
	isConfirmable := false

	for _, r := range dbWallet.FirewallRules {
		confirmationsCounter := uint(0)
		var orders Orders
		// Check if order addresses suit current rule
		if r.AddressType < targetAddress {
			continue
		}

		// Check if order amount doesn't exceed one time allowed amount
		if r.Period == 0 && r.Amount != 0 && o.Sum() > r.Amount {
			continue
		}

		// Check if order amount doesn't exceed allowed amount is specified period
		if r.Period > 0 {
			db.Find(&orders, "wallet_id = ? AND rule_id = ? AND created_at >= ?", o.WalletId, o.RuleId, time.Now().Add(-time.Second*time.Duration(r.Period)))
			sum := orders.Sum()
			if o.Sum()+sum > r.Amount {
				continue
			}
		}

		participantsList := new(ParticipantsList)
		err := participantsList.Unmarshal(r.Participants)
		if err != nil {
			continue
		}

		// split participants list to array of []uint
		switch r.ParticipantsType {
		case FirewallParticipantsTypeUsers:
			for _, u := range *participantsList {
				groupsAffected = append(groupsAffected, []uint{u})
			}
		case FirewallParticipantsTypeGroup:
			for _, g := range *participantsList {
				groupsAffected = append(groupsAffected, dbWallet.GetWalletUserGroupIds(WalletUserGroup(g)))
			}
		}
		for _, c := range o.Confirmations {
			for i, g := range groupsAffected {
				if h.UintFind(g, c.UserId) != -1 {
					var newGroupsAffected [][]uint
					for ig, gg := range groupsAffected {
						if ig == i { // skip group
							continue
						}
						newGroupsAffected = append(newGroupsAffected, h.Remove(gg, c.UserId)) // remove user ID from the confirmation from other user groups
					}
					groupsAffected = newGroupsAffected
					confirmationsCounter += 1
					break
				}
			}
		}
		// Order is already confirmed
		if confirmationsCounter >= r.ConfirmationsRequired {
			return nil, false, errors.New("Order is already confirmed")
		}

		for _, g := range groupsAffected {
			if h.UintFind(g, confirmation.UserId) != -1 {
				isConfirmable = true
				confirmationsCounter += 1
				if confirmationsCounter == r.ConfirmationsRequired && (activeRule == nil || activeRule.ConfirmationsRequired < r.ConfirmationsRequired) {
					activeRule = r
				}
			}
		}

	}
	return activeRule, isConfirmable, nil
}

func (o Order) String() string {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "{ID:%d Wallet:%d Asset: %s %s Status:%s}", o.ID, o.WalletId, o.Asset.Symbol, o.Destinations, o.Status)
	return ret.String()
}

type Orders []Order

func (o *Orders) Sum() (ret float64) {
	ret = 0
	for _, c := range []Order(*o) {
		for _, d := range c.Destinations {
			ret += d.Amount
		}
	}
	return
}
