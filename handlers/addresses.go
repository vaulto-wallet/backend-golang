package handlers

import (
	h "../helpers"
	hlp "../helpers"
	m "../models"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"math/big"
	"net/http"
	"strconv"
)

func CreateAddress(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	//user := req.Context().Value("user").(*m.User)
	masterPassword := req.Context().Value("masterPassword")

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
	db.Preload("Asset").Preload("Seed").First(&dbWallet, r.WalletID)

	if dbWallet.Asset.Type != m.AssetTypeBase {
		seedId := dbWallet.SeedId
		assets := new(m.Assets)
		db.Find(&assets)
		basicAsset := assets.GetBasicAsset(dbWallet.AssetId)
		if basicAsset.ID == 0 {
			ReturnError(w, Error(BadRequest))
		}

		db.First(&dbWallet, "asset_id = ? AND seed_id = ?", basicAsset.ID, seedId)
	}

	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Wallet not found")
		return
	}

	systemPrivKey := m.ConfigRecord{Name: "PrivateKey"}.Get(db)
	systemPublicKey := m.ConfigRecord{Name: "PublicKey"}.Get(db)

	encryptedSeed, err := hex.DecodeString(dbWallet.Seed.Seed)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Cannot decode seed")
		return
	}

	decryptedSeed, err := h.DecryptWithRSA(masterPassword.([]byte), []byte(systemPrivKey.Value), encryptedSeed)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Cannot decrypt seed code : "+err.Error())
		return
	}

	privateKey, address := hlp.GenerateAddress(dbWallet.Asset.Symbol, hex.EncodeToString(decryptedSeed), dbWallet.ChangeN, dbWallet.N)

	encryptedPrivateKey, err := h.EncryptWithRSA([]byte(systemPublicKey.Value), privateKey)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Cannot encrypt seed code")
		return
	}

	modelAddress := m.Address{
		Address:    address,
		PrivateKey: hex.EncodeToString(encryptedPrivateKey),
		WalletID:   dbWallet.ID,
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
	//user := req.Context().Value("user")

	vars := mux.Vars(req)
	var addresses m.Addresses

	if walletIdReq, ok := vars["wallet"]; !ok || walletIdReq == "0" {
		db.Find(&addresses)
		ReturnResult(w, addresses)
		return
	}

	walletId, _ := strconv.ParseUint(vars["wallet"], 10, 64)

	dbWallet := new(m.Wallet)
	db.Preload("Asset").First(&dbWallet, walletId)
	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Cannot find wallet")
		return
	}
	if dbWallet.Asset.Type != m.AssetTypeBase {
		seedId := dbWallet.SeedId
		assets := new(m.Assets)
		db.Find(assets)
		basicAsset := assets.GetBasicAsset(dbWallet.AssetId)
		dbWallet = new(m.Wallet)
		db.First(&dbWallet, "asset_id = ? AND seed_id = ?", basicAsset.ID, seedId)
	}

	db.Find(&addresses, "wallet_id = ?", dbWallet.ID)

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

	db.Preload("Asset").Find(&dbWallet, dbAddress.WalletID)

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

	db.Model(&dbAddress).Updates(struct{ Seqno uint64 }{dbAddress.Seqno})

	ReturnResult(w, true)
}
