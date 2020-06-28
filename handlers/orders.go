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

	if r.WalletId == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Invalid wallet")
		return
	}

	dbWallet := new(m.Wallet)
	db.Find(dbWallet, r.WalletId)

	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Wallet not found")
		return
	}

	newOrder := m.Order{
		Destinations:  r.Destinations,
		AssetId:       dbWallet.AssetId,
		WalletId:      r.WalletId,
		SubmittedById: user.ID,
		Comment:       r.Comment,
		Status:        m.OrderStatusNew,
	}
	db.Create(&newOrder)

	confirmation := m.Confirmation{
		OrderId: newOrder.ID,
		UserId:  user.ID,
	}

	_, confirmable, err := newOrder.FindRule(db, confirmation)
	if err != nil {
		db.Model(newOrder).Update(struct {
			Status m.OrderStatus
		}{m.OrderStatusRejected})
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, err.Error())
		return
	}
	if !confirmable {
		db.Model(newOrder).Update(struct {
			Status m.OrderStatus
		}{m.OrderStatusRejected})

		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "No rules found")
		return
	}

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
			return
		}

		var wallet m.Wallet
		db.Preload("Auditors").Preload("Coowners").Preload("Seed").Find(&wallet, walletId)

		if !wallet.HasReadPermission(user) {
			ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "No permission")
			return
		}

		db.Preload("Asset").Preload("Destinations").Find(orders, "wallet_id in (?)", []int{wId})
	} else {
		var wallets m.Wallets
		//db.Model(user).Preload("Coowners").Related(&wallets, "Wallets")
		db.Find(&wallets)
		db.Preload("Asset").Preload("Destinations").Preload("Confirmations").Find(orders, "wallet_id in (?)", wallets.GetIds())
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

func ConfirmOrder(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)
	vars := mux.Vars(req)

	var dbOrder m.Order

	if orderId, ok := vars["order"]; ok && orderId != "0" {
		oId, err := strconv.Atoi(orderId)

		if err != nil {
			ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Invalid order ID")
			return
		}
		db.Preload("Confirmations").First(&dbOrder, oId)
	}
	if dbOrder.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Order not found")
		return
	}

	var dbWallet m.Wallet
	db.Preload("Seed").Preload("FirewallRules").Preload("Coowners").Preload("Auditors").First(&dbWallet, dbOrder.WalletId)

	if dbWallet.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, "Wallet not found")
		return
	}

	if !dbWallet.HasWritePermission(user) {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "No permission")
		return
	}

	confirmation := m.Confirmation{
		OrderId: dbOrder.ID,
		UserId:  user.ID,
	}

	rule, confirmable, err := dbOrder.FindRule(db, confirmation)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), http.StatusBadRequest, err.Error())
		return
	}
	if !confirmable {
		ReturnErrorWithStatusString(w, Error(NotAuthorized), http.StatusForbidden, "No rules found")
		return
	}
	db.Save(&confirmation)

	if rule != nil && rule.ID != 0 {
		db.Model(dbOrder).Update(struct {
			RuleId uint
			Status m.OrderStatus
		}{rule.ID, m.OrderStatusConfirmed})
	}

	ReturnResult(w, true)
}
