package handlers

import (
	v "../api/vaulto"
	m "../models"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CreateOrder(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}
	var r v.OrderRequest
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
		AssetID:       r.AssetId,
		WalletID:      r.WalletId,
		SubmittedByID: dbUser.ID,
		Comment:       r.Comment,
		Status:        0,
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

	var r v.OrderRequest
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.First(&dbOrder, r.Id)

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
	var orders v.OrdersResponse

	db.Table("orders").Select("orders.id, orders.address_to, orders.amount, orders.comment, orders.status, assets.symbol").Joins("LEFT JOIN assets ON assets.id = orders.asset_id").Find(&orders)
	res, err := json.Marshal(orders)
	if err != nil {
		ReturnResult(w, orders)
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, orders)
}
