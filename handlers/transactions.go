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

func CreateTransaction(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}
	var r m.Transaction
	err := json.NewDecoder(req.Body).Decode(&r)

	orders := new([]*m.Order)
	db.Find(orders, r.OrderIds)

	addresses := new([]*m.Address)
	db.Find(addresses, r.AddressIds)

	r.Order = *orders
	r.Address = *addresses

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.Create(&r)
	ReturnResult(w, r.ID)
}

func UpdateTransaction(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := new(m.User)
	dbTransaction := new(m.Transaction)
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}

	var r m.Transaction
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.First(&dbTransaction, r.ID)

	if dbTransaction.Status != r.Status && r.Status != m.TransactionStatus(m.TransactionStatusUnknown) {
		dbTransaction.Status = r.Status
	}

	if dbTransaction.Tx != r.Tx && len(r.Tx) > 0 {
		dbTransaction.Tx = r.Tx
	}

	if dbTransaction.TxHash != r.TxHash && len(r.TxHash) > 0 {
		dbTransaction.TxHash = r.TxHash
	}

	if dbTransaction.TxData != r.TxData && len(r.TxData) > 0 {
		dbTransaction.TxData = r.TxData
	}

	db.Save(&dbTransaction)
	ReturnResult(w, true)
}

func GetTransactions(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	vars := mux.Vars(req)

	var transactions []m.Transaction

	if order, ok := vars["order"]; ok && order != "0" {
		if orderId, err := strconv.ParseUint(order, 10, 64); err != nil {
			ReturnError(w, Error(BadRequest))
		} else {
			var order m.Order
			db.First(&order, orderId)

			db.Model(&order).Related(&transactions, "Transaction")
		}
	} else {
		db.Find(&transactions)
	}

	res, err := json.Marshal(transactions)
	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, transactions)
}

func GetTransaction(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)
	var transaction m.Transaction

	vars := mux.Vars(req)

	transactionId, _ := strconv.ParseUint(vars["transaction"], 10, 64)

	db.First(&transaction, transactionId)
	if transactionId == 0 {
		ReturnError(w, Error(BadRequest))
	}

	ReturnResult(w, transaction)
}
