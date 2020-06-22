package models

import (
	"bytes"
	"fmt"
)

type PrivateKeyType int

const (
	PrivateKeyTypeUnknown = 0
	PrivateKeyTypeRoot    = 1
	PrivateKeyTypeSingle  = 2
	PrivateKeyTypeMulti   = 3
)

type WalletUserGroup uint

const (
	WalletUserGroupUnknown  = WalletUserGroup(0)
	WalletUserGroupOwners   = WalletUserGroup(1)
	WalletUserGroupCoowners = WalletUserGroup(2)
	WalletUserGroupAuditors = WalletUserGroup(3)
)

func (a WalletUserGroup) String() string {
	walletUserGroupText := [...]string{
		"Unknown",
		"Owners",
		"Coowners",
		"Auditors",
	}
	if uint(a) >= uint(len(walletUserGroupText)) {
		return walletUserGroupText[0]
	}

	return walletUserGroupText[a]
}

type Wallet struct {
	Model
	Name          string          `json:"name"`
	NetworkType   string          `json:"network_type"`
	SeedId        uint            `json:"seed_id"`
	Seed          Seed            `json:"-"`
	AssetId       uint            `json:"asset_id"`
	Asset         Asset           `json:"asset"`
	N             uint32          `json:"n"`
	ChangeN       uint32          `json:"change_n"`
	Seqno         uint32          `json:"seqno"`
	Balance       float64         `json:"balance"`
	Transaction   []*Transaction  `json:"-" gorm:"many2many:transaction_wallets"`
	Coowners      []*User         `json:"coowners" gorm:"many2many:wallet_owners"`
	Auditors      []*User         `json:"auditors" gorm:"many2many:wallet_auditors"`
	FirewallRules []*FirewallRule `json:"firewall_rules"`
}

func (o Wallet) String() string {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "{ID:%d \"%s\" %s N:%d Change:%d Seqno:%d}", o.ID, o.Asset.Symbol, o.Name, o.N, o.ChangeN, o.Seqno)
	return ret.String()
}

func (o *Wallet) IsOwner(user *User) bool {
	return o.Seed.OwnerId == user.ID
}

func (o *Wallet) HasWritePermission(user *User) bool {
	if o.IsOwner(user) {
		return true
	}
	for _, s := range o.Coowners {
		if s.ID == user.ID {
			return true
		}
	}
	return false
}

func (o *Wallet) HasReadPermission(user *User) bool {
	if o.IsOwner(user) {
		return true
	}
	for _, s := range o.Auditors {
		if s.ID == user.ID {
			return true
		}
	}

	for _, s := range o.Coowners {
		if s.ID == user.ID {
			return true
		}
	}
	return false
}

func (o *Wallet) GetWalletUserGroupIds(group WalletUserGroup) (ids []uint) {
	switch group {
	case WalletUserGroupOwners:
		ids = append(ids, o.Seed.OwnerId)
	case WalletUserGroupCoowners:
		for _, u := range o.Coowners {
			ids = append(ids, u.ID)
		}
	case WalletUserGroupAuditors:
		for _, u := range o.Coowners {
			ids = append(ids, u.ID)
		}
	}
	return
}

type Wallets []Wallet

func (w *Wallets) GetIds() (ret []uint) {
	for _, v := range []Wallet(*w) {
		ret = append(ret, v.ID)
	}
	return
}
