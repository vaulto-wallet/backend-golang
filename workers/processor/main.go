package main

import (
	"../../api/blockatlas"
	"../../api/vaulto"
	m "../../models"
	"./builder"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"math/big"
	"strings"
	"time"
)

var Assets m.Assets
var Vaulto *vaulto.VaultoAPI
var Blockatlas *blockatlas.BlockAtlasAPI

func processOrders() {
	orders, err := Vaulto.GetOrders()
	log.Println("Orders :", err, orders)
	for i, a := range Assets {
		log.Println("Assets ", i, a.Symbol)
	}
	for i, o := range orders {
		log.Println("Order ", i, o.Amount, o.Status)
		if o.Status == m.OrderStatusNew {
			log.Println("Processing order")

			gasPrice, _ := Blockatlas.GasPrice(Assets.Get(o.AssetId).Symbol)

			asset := Assets.Find(o.Symbol)
			if asset == nil {
				log.Println("Asset not found ", o.Symbol)
				continue
			}

			wallets, err := Vaulto.GetWalletsForAsset(o.Symbol)
			if err != nil {
				log.Println("Error fetching wallets for ", o.Symbol)
				continue
			}
			if len(wallets) == 0 {
				log.Println("No wallets found for ", o.Symbol)
				continue
			}
			wallet := wallets[0]

			addresses, err := Vaulto.GetAddressesForWallet(wallet.ID)
			if err != nil {
				log.Println("Error fetching addresses for ", o.Symbol, " wallet ", wallet.ID)
				continue
			}
			if len(wallets) == 0 {
				log.Println("No addresses found for ", o.Symbol, " wallet ", wallet.ID)
				continue
			}

			address := addresses[0]

			tx := builder.BuildEthereum([]byte(address.PrivateKey),
				o.AddressTo,
				*asset.ToBigInt(o.Amount),
				*new(big.Int).SetInt64(21000),
				*new(big.Int).SetInt64(gasPrice),
				*new(big.Int).SetUint64(address.Seqno))

			log.Println("Transaction built", tx)

			txId, err := Vaulto.CreateTransaction(asset.ID, wallet.ID, []uint{o.ID}, []uint{address.ID}, "", tx, "")
			log.Println("Transaction saved. ID : ", txId)

			Vaulto.UpdateOrder(o.ID, m.OrderStatusProcessing)
			Vaulto.UpdateAddress(address.ID, "", "", address.Seqno+1)
		}

		if o.Status == m.OrderStatusProcessing || o.Status == m.OrderStatusPartiallyProcessed {
			log.Println("Order in processing state", o.ID)
			transactions, err := Vaulto.GetTransactions()
			if err != nil {
				continue
			}
			transactionsDone := 0
			transactionsProgress := 0

			for _, t := range transactions {
				if t.Status == m.TransactionStatusConfirmed {
					transactionsDone += 1
				} else {
					transactionsProgress += 1
				}
			}
			if transactionsProgress == 0 && transactionsDone > 0 {
				Vaulto.UpdateOrder(o.ID, m.OrderStatusProcessed)
			} else if transactionsProgress > 0 && transactionsDone > 0 && o.Status != m.OrderStatusPartiallyProcessed {
				Vaulto.UpdateOrder(o.ID, m.OrderStatusPartiallyProcessed)
			}

		}

	}

}

func processTransactions() {
	transactions, _ := Vaulto.GetTransactions()
	for i, t := range transactions {
		log.Println("Transaction ", i, t.ID)

		if t.Status == m.TransactionStatusNew {
			asset := Assets.Get(t.AssetId)
			if len(t.Tx) == 0 {
				log.Println("No transaction data")
				continue
			}
			tx := struct {
				Tx string `json:"tx"`
			}{t.Tx}

			txHash, err := Blockatlas.SendTransaction(asset.Symbol, tx)

			if err != nil {
				Vaulto.UpdateTransaction(t.ID, m.TransactionStatusFailed, "", "", err.Error())
				log.Println("Problem sending transaction")
				continue
			}
			if len(txHash) == 0 {
				Vaulto.UpdateTransaction(t.ID, m.TransactionStatusPending, "", "", "")

			}
			Vaulto.UpdateTransaction(t.ID, m.TransactionStatusSent, "", txHash, "")
		}
	}
}

func scanTransactions() {
	wallets, err := Vaulto.GetWallets()
	if err != nil {
		return
	}

	transactions, err := Vaulto.GetTransactions()
	var txss = m.Transactions(transactions)

	if err != nil {
		return
	}
	for _, w := range wallets {
		addresses, err := Vaulto.GetAddressesForWallet(w.ID)
		if err != nil {
			log.Println("Error fetching addresses for wallet ", w.ID)
		}
		for _, a := range addresses {
			txs, err := Blockatlas.GetTXs(Assets.Get(w.AssetId).Symbol, a.Address)
			if err != nil {
				log.Println("Error fetching transactions for address ", a.Address)
			}
			newSeqno := uint64(0)
			for _, t := range txs {
				tx := txss.FindByHash(t.ID)
				if tx != nil {
					if tx.Status == m.TransactionStatusSent {
						txStatus := m.TransactionStatusPending
						if t.Status == "completed" {
							txStatus = m.TransactionStatusConfirmed
						}
						buf := new(bytes.Buffer)
						json.NewEncoder(buf).Encode(&t)
						Vaulto.UpdateTransaction(tx.ID, txStatus, "", "", buf.String())
					}
				} else {
					//TODO: No transaction found
				}
				if newSeqno <= t.Sequence && strings.ToLower(t.From) == strings.ToLower(a.Address) {
					newSeqno = t.Sequence + 1
				}
			}
			if a.Seqno < newSeqno {
				Vaulto.UpdateAddress(a.ID, a.Address, "", newSeqno)
			}

		}

	}
}

func main() {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}

	Vaulto = new(vaulto.VaultoAPI)
	Blockatlas = new(blockatlas.BlockAtlasAPI)

	Blockatlas.Init("http://localhost:8420/")

	vaultoUrl := viper.GetString("vaulto_url")
	vaultoUser := viper.GetString("vaulto_user")
	vaultoPassword := viper.GetString("vaulto_user")

	Vaulto.Init(vaultoUrl)

	Vaulto.Login(vaultoUser, vaultoPassword)

	Assets, _ = Vaulto.GetAssets()

	scanTransactions()

	go func() {
		for true {
			processOrders()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for true {
			processTransactions()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for true {
			time.Sleep(60 * time.Second)
			scanTransactions()
		}
	}()

	log.Println("Press any key to exit")
	var input string
	fmt.Scanln(&input, "%s")

}
