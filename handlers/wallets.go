package handlers

import (
	m "../models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func CreateWallet(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	var r m.Wallet
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	if r.SeedId == 0 || r.AssetId == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}

	seed := &m.Seed{}

	if err := seed.Load(db, r.SeedId); err != nil {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusBadRequest, "Seed not found")
		return
	}

	if seed.OwnerId != user.ID {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "No permission")
		return
	}

	asset := &m.Asset{}
	if err := asset.Load(db, r.AssetId); err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, err.(string))
		return
	}

	if len(r.Name) == 0 {
		r.Name = "New wallet"
	}

	newWallet := m.Wallet{
		Name:        r.Name,
		NetworkType: "",
		SeedId:      r.SeedId,
		AssetId:     r.AssetId,
		N:           0,
		ChangeN:     0,
		Coowners:    m.Users{},
		Auditors:    m.Users{},
	}

	db.Create(&newWallet)

	participantsList, _ := (&m.ParticipantsList{user.ID}).Marshal()

	firewallRule := m.FirewallRule{
		WalletId:              newWallet.ID,
		Wallet:                newWallet,
		ParticipantsType:      m.FirewallParticipantsTypeUsers,
		ParticipantsString:    participantsList,
		ConfirmationsRequired: 1,
		AddressType:           m.FirewallAddressTypeExternal,
		Amount:                0,
		Period:                0,
	}
	db.Create(&firewallRule)

	ReturnResult(w, newWallet.ID)
}

func GetWallet(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)
	vars := mux.Vars(req)

	walletId, ok := vars["wallet"]
	if !ok || walletId == "0" {
		ReturnError(w, Error(BadRequest))
		return
	}

	var dbWallet m.Wallet
	db.Preload("FirewallRules").Preload("Coowners").Preload("Auditors").Preload("Seed", "owner_id = ?", user.ID).First(&dbWallet, walletId)
	if dbWallet.ID == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}

	for _, v := range dbWallet.FirewallRules {
		participantsList := new(m.ParticipantsList)

		err := participantsList.Unmarshal(v.ParticipantsString)
		if err != nil {
			continue
		}
		v.ParticipantIds = *participantsList
	}

	ReturnResult(w, dbWallet)
}

func GetWallets(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	var wallets []m.Wallet

	db.Model(user).Preload("Coowners").Preload("Auditors").Related(&wallets, "Wallets")
	db.Preload("Coowners").Preload("Auditors").Preload("Seed", "owner_id = ?", user.ID).Find(&wallets)

	//db.Find(wallets, "")

	ReturnResult(w, wallets)
}

func GetWalletsForAsset(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	//user := req.Context().Value("user")

	vars := mux.Vars(req)

	asset := vars["asset"]

	var wallets m.Wallets

	db.Table("wallets").
		Select("wallets.id, wallets.seed_id, wallets.asset_id, assets.symbol").
		Joins("join assets on assets.id = wallets.asset_id").
		Where("assets.symbol = ?", asset).
		Find(&wallets)

	res, err := json.Marshal(wallets)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, err.Error())
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, wallets)
}

func ShareWallet(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)
	vars := mux.Vars(req)

	dbWallet := new(m.Wallet)

	if wallet_param, ok := vars["wallet"]; !ok {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Cannot get wallet id")
		return
	} else {
		if wallet_id, ok := strconv.Atoi(wallet_param); ok != nil {
			ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Cannot parse wallet id")
			return
		} else {
			db.Preload("Seed").Find(&dbWallet, wallet_id)
		}

	}

	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Wallet not found")
		return
	}

	if !dbWallet.IsOwner(user) {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "Not authorized")
		return
	}

	var r struct {
		Coowners *[]uint `json:"coowners,omitempty"`
		Auditors *[]uint `json:"auditors,omitempty"`
	}

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Invalid shared list")
		return
	}

	coowners := new([]m.User)
	db.Find(&coowners, "id in (?)", *r.Coowners)

	auditors := new([]m.User)
	db.Find(&auditors, "id in (?)", *r.Auditors)

	dbWallet.Coowners = []*m.User{}
	dbWallet.Auditors = []*m.User{}

	for _, o := range *coowners {
		dbWallet.Coowners = append(dbWallet.Coowners, &o)
	}

	for _, o := range *auditors {
		dbWallet.Auditors = append(dbWallet.Auditors, &o)
	}

	db.Save(dbWallet)

	ReturnResult(w, true)
}
