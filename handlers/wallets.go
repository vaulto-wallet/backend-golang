package handlers

import (
	v "../api/vaulto"
	m "../models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CreateWallet(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}
	var r v.WalletRequest
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

	if err := seed.Get(db, r.SeedId); err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, err.(string))
		return
	}

	asset := &m.Asset{}
	if err := asset.Get(db, r.AssetId); err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, err.(string))
		return
	}

	if len(r.Name) == 0 {
		r.Name = "New wallet"
	}

	newWallet := m.Wallet{
		Name:        r.Name,
		NetworkType: "",
		SeedID:      r.SeedId,
		AssetID:     r.AssetId,
		N:           0,
		ChangeN:     0,
	}

	db.Create(&newWallet)
	ReturnResult(w, true)
}

func GetWallets(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	var wallets v.WalletsResponse

	db.Table("wallets").
		Select("wallets.id, wallets.seed_id, wallets.asset_id, assets.symbol").
		Joins("JOIN assets ON assets.id = wallets.asset_id").
		Find(&wallets)

	res, err := json.Marshal(wallets)
	if err != nil {
		ReturnResult(w, wallets)
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, wallets)
}

func GetWalletsForAsset(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	vars := mux.Vars(req)

	asset := vars["asset"]

	var wallets v.WalletsResponse

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
