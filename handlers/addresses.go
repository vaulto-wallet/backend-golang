package handlers

import (
	m "../models"
	hlp "../trusthelpers"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"math/big"
	"net/http"
	"strconv"
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
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	vars := mux.Vars(req)

	walletId, _ := strconv.ParseUint(vars["wallet"], 10, 64)

	var addresses m.Addresses

	db.Find(&addresses, "wallet_id = ?", walletId)

	res, err := json.Marshal(addresses)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, err.Error())
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, addresses)
}

func UpdateAddress(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}

	var r m.Address
	var dbAddress m.Address

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	if r.ID != 0 {
		db.Find(&dbAddress, r.ID)
	} else if len(r.Address) > 0 {
		db.Find(&dbAddress, "address = ?", r.Address)
	}
	if dbAddress.ID == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}

	var dbWallet m.Wallet

	db.Set("gorm:auto_preload", true).Find(&dbWallet, dbAddress.WalletID)

	if dbWallet.ID == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}

	if len(r.BalanceInt) > 0 {
		base := 10
		numString := r.BalanceInt
		if r.BalanceInt[0:2] == "0x" {
			numString = r.BalanceInt[2:]
			base = 16
		}
		if intBalance, isOk := new(big.Int).SetString(numString, base); isOk == true {
			dbAddress.Balance = dbWallet.Asset.ToFloat(intBalance)
			dbAddress.BalanceInt = numString
		}
	}

	if r.Seqno > dbAddress.Seqno {
		dbAddress.Seqno = r.Seqno
	}

	db.Save(&dbAddress)

	ReturnResult(w, true)
}
