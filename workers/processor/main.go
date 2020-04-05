package main

import (
	"../../api/vaulto"
	m "../../models"
	"./builder"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"math/big"
	"time"
)

func processOrders(api *vaulto.VaultoAPI) {
	orders, err := api.GetOrders()
	log.Println("Orders :", err, orders)
	for i, a := range Assets {
		log.Println("Assets ", i, a.Symbol)
	}
	for i, o := range orders {
		log.Println("OrderData ", i, o.Amount, o.Status)
		if o.Status == m.OrderStatusNew {
			log.Println("Processing order")

			asset := Assets.Find(o.Symbol)
			if asset == nil {
				log.Println("Asset not found ", o.Symbol)
				continue
			}

			wallets, err := api.GetWalletsForAsset(o.Symbol)
			if err != nil {
				log.Println("Error fetching wallets for ", o.Symbol)
				continue
			}
			if len(wallets) == 0 {
				log.Println("No wallets found for ", o.Symbol)
				continue
			}
			wallet := wallets[0]

			addresses, err := api.GetAddressesForWallet(wallet.ID)
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
				*new(big.Int).SetInt64(10000000),
				*new(big.Int).SetUint64(address.Seqno))

			log.Println("Transaction built", tx)
			//api.UpdateOrder(o.ID, m.OrderStatusProcessing)
		}

	}

}

var Assets m.Assets

func main() {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}

	vaulto := vaulto.VaultoAPI{}

	vaultoUrl := viper.GetString("vaulto_url")
	vaultoUser := viper.GetString("vaulto_user")
	vaultoPassword := viper.GetString("vaulto_user")

	vaulto.Init(vaultoUrl)

	vaulto.Login(vaultoUser, vaultoPassword)

	Assets, _ = vaulto.GetAssets()

	go func() {
		for true {
			processOrders(&vaulto)
			time.Sleep(2 * time.Second)
		}
	}()

	log.Println("Press any key to exit")
	var input string
	fmt.Scanln(&input, "%s")

}
