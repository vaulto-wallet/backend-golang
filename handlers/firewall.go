package handlers

import (
	m "../models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func CreateRule(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	var r m.FirewallRule
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil || r.AddressType == m.FirewallAddressTypeUnknown || r.ParticipantsType == m.FirewallParticipantsTypeUnknown || len(r.Participants) == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}

	dbWallet := new(m.Wallet)
	db.Find(dbWallet, r.WalletId)

	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Wallet not found")
		return
	}

	if !dbWallet.IsOwner(user) {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "Not authorized")
		return
	}

	db.Create(&r)
	ReturnResult(w, r.ID)
}

func GetRule(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	vars := mux.Vars(req)
	var rule m.FirewallRule

	if ruleVar, ok := vars["rule"]; ok {
		if ruleId, err := strconv.ParseUint(ruleVar, 10, 64); err != nil {
			ReturnError(w, Error(BadRequest))
			return
		} else {
			db.First(&rule, ruleId)
		}
	}

	if rule.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Firewall rule not found")
	}

	var dbWallet m.Wallet
	db.First(&dbWallet, rule.WalletId)

	if !dbWallet.IsOwner(user) {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "Not authorized")
		return
	}

	ReturnResult(w, rule)
}

func GetRules(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	vars := mux.Vars(req)
	var rules m.FirewallRules
	var dbWallet m.Wallet

	if walletVar, ok := vars["wallet"]; ok {
		if walletId, err := strconv.ParseUint(walletVar, 10, 64); err != nil {
			ReturnError(w, Error(BadRequest))
			return
		} else {
			db.First(&dbWallet, walletId)
		}
	}

	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Wallet not found")
	}

	if !dbWallet.IsOwner(user) {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "Not authorized")
		return
	}

	db.Model(&dbWallet).Related(&rules)

	ReturnResult(w, rules)
}
