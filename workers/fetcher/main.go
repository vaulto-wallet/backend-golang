package main

import (
	"../../api/alfaex"
	a "../../api/vaulto"
	m "../../models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"log"
	"math/big"
	"time"
)

var Vaulto = new(a.VaultoAPI)
var Alfa = new(alfaex.AlfaEXAPI)
var ConnectorDB *gorm.DB

type Connector struct {
	gorm.Model
	ExternalID        string
	ExternalData      string
	ExternalStatus    int
	VaultoOrderID     uint
	VaultoOrderStatus m.OrderStatus
	VaultoTxHash      string
}

type Connectors []Connector

func (c *Connectors) FindByExternalId(externalId string) *Connector {
	for _, c := range []Connector(*c) {
		if c.ExternalID == externalId {
			return &c
		}
	}
	return nil
}

func ProcessOperations() {
	operations := Alfa.GetOperations()
	var connector Connectors

	ConnectorDB.Find(&connector)

	for _, o := range operations.Operation {
		amount := big.Float{}
		amount.SetString(o.Arguments.Amount)
		amountFloat, _ := amount.Float64()

		log.Println(o.Id, o.Arguments.Amount, amount.Text('f', 8))

		connectorEntry := connector.FindByExternalId(o.Id)

		if connectorEntry == nil {

			Alfa.PutOperation(alfaex.WithdrawOperation{
				OperationID: o.Id,
				State:       -5,
			})

			orderId, err := Vaulto.CreateOrder(o.Arguments.Currency, o.Arguments.DestinationWallet, amountFloat, o.Id)
			if err != nil {
				log.Println("Error creating order for ", o.Id)
				continue
			}

			ConnectorDB.Create(&Connector{
				ExternalID:        o.Id,
				ExternalData:      "",
				ExternalStatus:    -5,
				VaultoOrderID:     orderId,
				VaultoOrderStatus: m.OrderStatusNew,
				VaultoTxHash:      "",
			})

		}
	}

	for _, connectorEntry := range connector {
		order, err := Vaulto.GetOrder(connectorEntry.VaultoOrderID)
		if err != nil {
			log.Println("Error fetching order")
			continue
		}
		if connectorEntry.VaultoOrderStatus != m.OrderStatusProcessed && order.Status != connectorEntry.VaultoOrderStatus {
			connectorEntry.VaultoOrderStatus = order.Status
			if order.Status == m.OrderStatusProcessed {
				transactions, err := Vaulto.GetOrderTransactions(connectorEntry.VaultoOrderID)
				if err == nil && len(transactions) > 0 {
					connectorEntry.VaultoTxHash = transactions[0].TxHash
				}

				if err != nil {
					log.Println("Error creating order for ", connectorEntry.ExternalID)
					continue
				}

				Alfa.PutOperation(alfaex.WithdrawOperation{
					OperationID:           connectorEntry.ExternalID,
					State:                 1,
					ExternalTransactionID: transactions[0].TxHash,
				})
			}
			ConnectorDB.Save(&connectorEntry)
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

	ConnectorDB, err = gorm.Open("sqlite3", "connector.db")
	ConnectorDB.AutoMigrate(&Connector{})

	publicKey := viper.GetString("public_key")
	secretKey := viper.GetString("secret_key")
	externalUrl := viper.GetString("external_url")

	vaultoUrl := viper.GetString("vaulto_url")
	vaultoUser := viper.GetString("vaulto_user")
	vaultoPassword := viper.GetString("vaulto_user")

	Vaulto.Init(vaultoUrl)

	Vaulto.Login(vaultoUser, vaultoPassword)

	Alfa.Init(publicKey, secretKey, externalUrl)

	go func() {
		for true {
			ProcessOperations()
			time.Sleep(60 * time.Second)
		}
	}()

	log.Println("Press any key to exit")
	var input string
	fmt.Scanln(&input, "%s")

}
