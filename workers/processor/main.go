package main

import "C"
import (
	"../../api/blockatlas"
	"../../api/vaulto"
	h "../../helpers"
	m "../../models"
	"./builder"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

var Assets m.Assets
var Vaulto *vaulto.VaultoAPI
var Blockatlas *blockatlas.BlockAtlasAPI

/*
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

			gasPrice, _ := Blockatlas.GasPrice(Assets.Load(Assets.GetBasicAsset(o.AssetId).ID).Symbol)

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

			if asset.Symbol == "ETH" {
				tx := builder.BuildEthereum([]byte(address.PrivateKey),
					o.AddressTo,
					*asset.ToBigInt(o.Amount),
					*new(big.Int).SetInt64(21000),
					*new(big.Int).SetInt64(gasPrice),
					*new(big.Int).SetUint64(address.Seqno),
					[]byte{})

				log.Println("ETH Transaction built", tx)
				txId, err := Vaulto.CreateTransaction([]uint{asset.ID}, []uint{wallet.ID}, []uint{o.ID}, []uint{address.ID}, "", tx, "")
				log.Println("Transaction saved. ID : ", txId, err)
				Vaulto.UpdateOrder(o.ID, m.OrderStatusProcessing)
				Vaulto.UpdateAddress(address.ID, "", "", address.Seqno+1)
			} else if asset.Symbol == "BTCT" {

			} else if asset.Type == m.AssetTypeERC20 {
				payload := builder.BuildERC20Transfer(o.AddressTo, *asset.ToBigInt(o.Amount))
				tx := builder.BuildEthereum([]byte(address.PrivateKey),
					asset.Address,
					*new(big.Int).SetInt64(0),
					*new(big.Int).SetInt64(200000),
					*new(big.Int).SetInt64(gasPrice),
					*new(big.Int).SetUint64(address.Seqno),
					payload)

				txId, err := Vaulto.CreateTransaction([]uint{Assets.GetBasicAsset(asset.ID).ID}, []uint{wallet.ID}, []uint{o.ID}, []uint{address.ID}, "", tx, "")
				log.Println("Transaction saved. ID : ", txId, err)

				log.Println("ERC20 Transaction built", tx)
				Vaulto.UpdateOrder(o.ID, m.OrderStatusProcessing)
				Vaulto.UpdateAddress(address.ID, "", "", address.Seqno+1)
			} else {

			}
		}

		if o.Status == m.OrderStatusProcessing || o.Status == m.OrderStatusPartiallyProcessed {
			log.Println("Order in processing state", o.ID)
			transactions, err := Vaulto.GetOrderTransactions(o.ID)
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

}*/

func processOrders() {
	orders, err := Vaulto.GetOrders()
	log.Println("Orders :", err, orders)
	for i, a := range Assets {
		log.Println("Assets ", i, a.Symbol)
	}

	wallets, err := Vaulto.GetWallets()
	if err != nil {
		return
	}

	for _, w := range wallets {
		orders, err := Vaulto.GetOrdersForWallet(w.ID)
		if err != nil || len(orders) == 0 {
			log.Println("No orders for Wallet", w.ID)
			continue
		}
		addresses, err := Vaulto.GetAddressesForWallet(w.ID)
		if err != nil || len(orders) == 0 {
			log.Println("No addresses for Wallet", w.ID)
			continue
		}

		transactions, err := Vaulto.GetTransactionsForWallet(w.ID)
		if w.Asset.Symbol == "BTCT" {
			builder.BuiltBitcoinTransactions(orders, addresses, transactions)
		}
	}

}

func processTransactions() {
	transactions, _ := Vaulto.GetTransactions()
	for i, t := range transactions {
		log.Println("Transaction ", i, t.ID)

		if t.Status == m.TransactionStatusNew {
			asset := Assets.Get(t.AssetId[0])
			if len(t.Tx) == 0 {
				log.Println("No transaction data")
				continue
			}
			tx := struct {
				Tx string `json:"tx"`
			}{t.Tx}

			txHash, err := Blockatlas.SendTransaction(Assets.GetBasicAsset(asset.ID).Symbol, tx)

			if err != nil {
				Vaulto.UpdateTransaction(t.ID, m.TransactionStatusFailed, "", "", err.Error())
				log.Println("Problem sending transaction")
				continue
			}
			if len(txHash) == 0 {
				Vaulto.UpdateTransaction(t.ID, m.TransactionStatusFailed, "", "", "")

			} else {
				Vaulto.UpdateTransaction(t.ID, m.TransactionStatusSent, "", txHash, "")
			}
		}
	}
}

func scanTransactions() {
	wallets, err := Vaulto.GetWallets()
	if err != nil {
		return
	}

	allAddrArray, err := Vaulto.GetAddressesForWallet(0)
	var allAddresses = m.Addresses(allAddrArray)

	//transactions, err := Vaulto.GetTransactions()
	//var txss = m.Transactions(transactions)

	if err != nil {
		return
	}
	for _, w := range wallets {
		if Assets.Get(w.AssetId).Type != m.AssetTypeBase {
			continue
		}
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
				tx, err := Vaulto.GetTransactionByHash(t.ID)
				if err != nil {
					log.Println("Error fetching transaction : ", tx.ID)
					continue
				}

				if tx != nil && tx.ID != 0 {
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
					log.Println("New transaction : ", t)
					affectedAddresses := *new([]uint)
					affectedWallets := *new([]uint)
					affectedAssets := *new([]uint)
					if addr := allAddresses.FindAddress(t.From); addr != nil {
						affectedAddresses = h.UintAppendNew(affectedAddresses, addr.ID)
						affectedWallets = h.UintAppendNew(affectedWallets, addr.WalletID)
						affectedAssets = h.UintAppendNew(affectedAssets, w.AssetId)

					}
					if addr := allAddresses.FindAddress(t.To); addr != nil {
						affectedAddresses = h.UintAppendNew(affectedAddresses, addr.ID)
						affectedWallets = h.UintAppendNew(affectedWallets, addr.WalletID)
						affectedAssets = h.UintAppendNew(affectedAssets, w.AssetId)
					}

					for _, input := range t.Inputs {
						if addr := allAddresses.FindAddress(input.Address); addr != nil {
							affectedAddresses = h.UintAppendNew(affectedAddresses, addr.ID)
							affectedWallets = h.UintAppendNew(affectedWallets, addr.WalletID)
							affectedAssets = h.UintAppendNew(affectedAssets, w.AssetId)
						}
					}

					for _, output := range t.Outputs {
						if addr := allAddresses.FindAddress(output.Address); addr != nil {
							affectedAddresses = h.UintAppendNew(affectedAddresses, addr.ID)
							affectedWallets = h.UintAppendNew(affectedWallets, addr.WalletID)
							affectedAssets = h.UintAppendNew(affectedAssets, w.AssetId)
						}
					}
					buf := new(bytes.Buffer)
					json.NewEncoder(buf).Encode(&t)

					txId, err := Vaulto.CreateTransaction(affectedAssets, affectedWallets, []uint{}, affectedAddresses,
						t.ID, "", buf.String())
					log.Println("Error", err, "Tx ID :", txId, "TxHash :", t.ID, "Addresses : ", affectedAddresses, "Wallets : ", affectedWallets, "Assets :", affectedAssets)
					Vaulto.UpdateTransaction(txId, m.TransactionStatusConfirmed, "", "", "")

				}
				if w.Asset.Symbol == "ETH" && newSeqno <= t.Sequence && strings.ToLower(t.From) == strings.ToLower(a.Address) {
					newSeqno = t.Sequence + 1
				}
			}
			if w.Asset.Symbol == "ETH" && a.Seqno < newSeqno {
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
			break
		}
	}()

	go func() {
		for true {
			//processTransactions()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for true {
			time.Sleep(60 * time.Second)
			//scanTransactions()
		}
	}()

	log.Println("Press any key to exit")
	var input string
	fmt.Scanln(&input, "%s")

}
