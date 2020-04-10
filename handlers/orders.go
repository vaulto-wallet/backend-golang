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
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}
	var r m.Order
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	if r.AssetId == 0 {
		if len(r.Symbol) == 0 {
			ReturnError(w, Error(BadRequest))
			return
		}
		var asset m.Asset
		db.Where("symbol = ?", r.Symbol).First(&asset)
		r.AssetId = asset.ID
	}

	newOrder := m.Order{
		Amount:        r.Amount,
		AddressTo:     r.AddressTo,
		AssetId:       r.AssetId,
		WalletId:      r.WalletId,
		SubmittedById: dbUser.ID,
		Comment:       r.Comment,
		Status:        m.OrderStatusNew,
	}

	db.Create(&newOrder)
	ReturnResult(w, newOrder.ID)
}

func UpdateOrder(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	dbOrder := m.Order{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}

	var r m.Order
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.First(&dbOrder, r.ID)

	if dbOrder.Status != (m.OrderStatus)(r.Status) {
		dbOrder.Status = (m.OrderStatus)(r.Status)
	}

	if dbOrder.Comment != r.Comment && len(r.Comment) > 0 {
		dbOrder.Comment = r.Comment
	}

	db.Save(&dbOrder)
	ReturnResult(w, true)
}

func GetOrders(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)
	var orders m.Orders

	db.Table("orders").Select("orders.id, orders.address_to, orders.asset_id, orders.amount, orders.comment, orders.status, orders.submitted_by_id, orders.wallet_id, assets.symbol").Joins("LEFT JOIN assets ON assets.id = orders.asset_id").Find(&orders)
	res, err := json.Marshal(orders)
	if err != nil {
		ReturnResult(w, orders)
		return
	}

	fmt.Println((string)(res))
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
