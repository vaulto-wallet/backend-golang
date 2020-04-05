package handlers

import (
	m "../models"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
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
	var transactions m.Transactions

	db.Find(&transactions)
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
	var transactions m.Transactions

	db.Find(&transactions)
	res, err := json.Marshal(transactions)
	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, transactions)
}
