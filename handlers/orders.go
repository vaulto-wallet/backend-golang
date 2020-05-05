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

func CreateOrder(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	var r m.Order
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	dbWallet := new(m.Wallet)
	if r.WalletId == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Invalid wallet")
	}

	db.Find(dbWallet, r.WalletId)

	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Wallet not found")
	}

	newOrder := m.Order{
		Amount:        r.Amount,
		AddressTo:     r.AddressTo,
		AssetId:       dbWallet.AssetId,
		WalletId:      r.WalletId,
		SubmittedById: user.ID,
		Comment:       r.Comment,
		Status:        m.OrderStatusNew,
	}

	db.Create(&newOrder)
	ReturnResult(w, newOrder.ID)
}

func UpdateOrder(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	//user := req.Context().Value("user")
	dbOrder := m.Order{}

	var r m.Order

	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.First(&dbOrder, r.ID)

	if dbOrder.Status != r.Status {
		dbOrder.Status = r.Status
	}

	if dbOrder.Comment != r.Comment && len(r.Comment) > 0 {
		dbOrder.Comment = r.Comment
	}

	db.Save(&dbOrder)
	ReturnResult(w, true)
}

func GetOrders(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	orders := new(m.Orders)

	vars := mux.Vars(req)

	if walletId, ok := vars["wallet"]; ok && walletId != "0" {
		wId, err := strconv.Atoi(walletId)
		if err != nil {
			ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Invalid wallet ID")
		}
		db.Preload("Asset").Find(orders, "wallet_id = ?", wId)
	} else {
		db.Preload("Asset").Find(orders)
	}
	fmt.Println(orders)
	ReturnResult(w, orders)
}

func GetOrder(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)
	vars := mux.Vars(req)
	var order m.Order

	if orderId, ok := vars["order"]; ok {
		if orderId, err := strconv.ParseUint(orderId, 10, 64); err != nil {
			ReturnError(w, Error(BadRequest))
		} else {
			db.Find(&order, orderId)
		}
	}

	ReturnResult(w, order)
}
