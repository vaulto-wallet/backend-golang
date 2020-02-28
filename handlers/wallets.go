package handlers

import (
	m "../models"
	"encoding/json"
	"fmt"
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
	var r m.Wallet
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	if r.SeedID == 0 || r.AssetID == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}

	seed := &m.Seed{}

	if err := seed.Get(db, r.SeedID); err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, err.(string))
		return
	}

	asset := &m.Asset{}
	if err := asset.Get(db, r.AssetID); err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, err.(string))
		return
	}

	if len(r.Name) == 0 {
		r.Name = "New wallet"
	}

	if len(r.NetworkType) == 0 {
		r.NetworkType = "main"
	}

	db.Create(&r)
	ReturnResult(w, true)
}

func GetWallets(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)
	var wallets m.Wallets

	//db.Where("owner_id = ? ", dbUser.ID).Find(&wallets)
	db.Find(&wallets)
	res, err := json.Marshal(wallets)
	if err != nil {
		ReturnResult(w, wallets)
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, wallets)
}
