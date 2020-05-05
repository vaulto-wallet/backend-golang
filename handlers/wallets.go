package handlers

import (
	h "../helpers"
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
		Owners:      m.Users{user},
	}

	db.Create(&newWallet)
	ReturnResult(w, newWallet.ID)
}

func GetWallets(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	var wallets []m.Wallet

	db.Model(user).Preload("Owners").Related(&wallets, "Wallets")

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
		Owners []uint `json:"owners"`
	}

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Invalid shared list")
		return
	}

	if !h.UintInArray(r.Owners, user.ID) {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Owner is not is list")
		return
	}

	owners := new([]m.User)
	db.Find(&owners, "id in (?)", r.Owners)
	dbWallet.Owners = []*m.User{}
	for _, o := range *owners {
		dbWallet.Owners = append(dbWallet.Owners, &o)
	}

	db.Save(dbWallet)

	ReturnResult(w, true)
}
