package models

import (
	h "../helpers"
	"encoding/json"
)

type FirewallAddressType uint

const (
	FirewallAddressTypeUnknown     = FirewallAddressType(0)
	FirewallAddressTypeInternal    = FirewallAddressType(1)
	FirewallAddressTypeExternal    = FirewallAddressType(2)
	FirewallAddressTypeWhitelisted = FirewallAddressType(3)
)

func (a FirewallAddressType) String() string {
	firewallAddressTypeText := [...]string{
		"Unknown",
		"Internal",
		"External",
		"Whitelisted",
	}
	if uint(a) >= uint(len(firewallAddressTypeText)) {
		return firewallAddressTypeText[0]
	}

	return firewallAddressTypeText[a]
}

type FirewallParticipantsType uint

const (
	FirewallParticipantsTypeUnknown = FirewallParticipantsType(0)
	FirewallParticipantsTypeUsers   = FirewallParticipantsType(1)
	FirewallParticipantsTypeGroup   = FirewallParticipantsType(2)
)

func (a FirewallParticipantsType) String() string {
	firewallParticipantsTypeText := [...]string{
		"Unknown",
		"Users",
		"Groups",
	}
	if uint(a) >= uint(len(firewallParticipantsTypeText)) {
		return firewallParticipantsTypeText[0]
	}
	return firewallParticipantsTypeText[a]
}

type ParticipantsList []uint

func (o *ParticipantsList) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), o)
}

func (o *ParticipantsList) Marshal() (string, error) {
	r, e := json.Marshal(o)
	return string(r), e
}

type FirewallRule struct {
	Model
	WalletId              uint                     `json:"wallet_id"`
	Wallet                Wallet                   `json:"wallet"`
	ParticipantsType      FirewallParticipantsType `json:"participant_type"`
	ParticipantsString    string                   `json:"-"`
	ParticipantIds        ParticipantsList         `json:"participants" gorm:"-"`
	ConfirmationsRequired uint                     `json:"confirmations_required"`
	AddressType           FirewallAddressType      `json:"address_type"`
	Amount                float64                  `json:"amount"`
	Period                uint                     `json:"period"`
}

type FirewallRules []FirewallRule

func (o *FirewallRules) AffectedUsers(wallet *Wallet) (ret []uint) {
	for _, r := range []FirewallRule(*o) {
		participantsList := new(ParticipantsList)
		err := participantsList.Unmarshal(r.ParticipantsString)
		if err != nil {
			continue
		}
		for _, u := range *participantsList {
			switch r.ParticipantsType {
			case FirewallParticipantsTypeUsers:
				ret = h.UintAppendNew(ret, u)
			case FirewallParticipantsTypeGroup:
				ret = h.UintAppendNewArray(ret, wallet.GetWalletUserGroupIds(WalletUserGroup(u)))
			}
		}
	}
	return
}

func (o *FirewallRules) IsConfirmed(confirmations Confirmations) bool {
	return false
}
