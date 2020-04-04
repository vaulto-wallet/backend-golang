package handlers

import (
	m "../models"
	hlp "../trusthelpers"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CreateAddress(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}

	var r m.Address
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	if r.WalletID == 0 {
		ReturnError(w, Error(BadRequest))
	}

	dbWallet := m.Wallet{}
	db.Set("gorm:auto_preload", true).Find(&dbWallet, r.WalletID)

	if dbWallet.ID == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}
	private_key, address := hlp.GenerateAddress(dbWallet.Asset.Symbol, dbWallet.Seed.Seed, dbWallet.ChangeN, dbWallet.N)

	modelAddress := m.Address{
		Address:    address,
		PrivateKey: private_key,
		WalletID:   r.WalletID,
		N:          dbWallet.N,
		Change:     0,
		Comment:    "",
		BalanceInt: "0",
		Balance:    0,
	}
	dbWallet.N++
	db.Save(&dbWallet)
	db.Create(&modelAddress)

	ReturnResult(w, modelAddress.ID)

}

func GetAddress(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	ReturnError(w, Error(NotImplemented))
}

func UpdateAddress(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	ReturnError(w, Error(NotImplemented))
}
