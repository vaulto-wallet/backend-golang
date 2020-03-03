package handlers

import (
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
	var r m.Order
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	if r.WalletID == 0 {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.Create(&r)
	ReturnResult(w, true)
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

	db.First(&dbOrder, "ID = ?", r.ID)

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
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)
	var orders m.Orders

	db.Find(&orders)
	res, err := json.Marshal(orders)
	if err != nil {
		ReturnResult(w, orders)
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, orders)
}
