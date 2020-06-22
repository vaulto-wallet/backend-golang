package handlers

import (
	m "../models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

func CreateTransaction(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	var r m.Transaction
	err := json.NewDecoder(req.Body).Decode(&r)

	orders := new([]*m.Order)
	db.Find(orders, "id IN (?)", r.OrderId)

	addresses := new([]*m.Address)
	db.Find(addresses, "id IN (?)", r.AddressId)

	assets := new([]*m.Asset)
	db.Find(assets, "id IN (?)", r.AssetId)

	wallets := new([]*m.Wallet)
	db.Find(wallets, "id IN (?)", r.WalletId)

	r.Order = *orders
	r.Address = *addresses
	r.Asset = *assets
	r.Wallet = *wallets

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.Create(&r)
	ReturnResult(w, r.ID)
}

func UpdateTransaction(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	dbTransaction := new(m.Transaction)

	var r m.Transaction
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.First(&dbTransaction, r.ID)

	if dbTransaction.ID == 0 {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Transactions not found")
		return
	}

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

	db.Model(&dbTransaction).Updates(struct {
		Status m.TransactionStatus
		Tx     string
		TxData string
		TxHash string
	}{
		dbTransaction.Status, dbTransaction.Tx, dbTransaction.TxData, dbTransaction.TxHash})
	ReturnResult(w, true)
}

func GetTransactions(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var transactions []m.Transaction

	if order, ok := vars["order"]; ok && order != "0" {
		if orderId, err := strconv.ParseUint(order, 10, 64); err != nil {
			ReturnError(w, Error(BadRequest))
			return
		} else {
			var order m.Order
			db.First(&order, orderId)
			db.Model(&order).Related(&transactions, "Transactions")
		}
	} else if wallet, ok := vars["wallet"]; ok && order != "0" {
		if walletId, err := strconv.ParseUint(wallet, 10, 64); err != nil {
			ReturnError(w, Error(BadRequest))
		} else {
			var wallet m.Wallet
			db.First(&wallet, walletId)
			db.Model(&wallet).Related(&transactions, "Transactions")
		}
	} else {
		db.Find(&transactions)
	}

	ReturnResult(w, transactions)
}

func GetTransaction(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	transaction := new(m.Transaction)

	vars := mux.Vars(req)

	if txValue, ok := vars["id"]; ok && txValue != "0" {
		transactionId, _ := strconv.ParseUint(txValue, 10, 64)
		db.Set("gorm:auto_preload", true).First(transaction, transactionId)
	} else if txHashValue, ok := vars["txhash"]; ok && len(txHashValue) > 0 {
		txHash := strings.ToLower(txHashValue)
		db.Set("gorm:auto_preload", true).First(transaction, "tx_hash = ?", txHash)
	} else {
		ReturnError(w, Error(BadRequest))
	}

	ReturnResult(w, transaction)
}
