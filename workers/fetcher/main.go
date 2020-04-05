package main

import (
	"../../api/alfaex"
	a "../../api/vaulto"
	"github.com/spf13/viper"
	"log"
	"math/big"
)

func main() {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}

	publicKey := viper.GetString("public_key")
	secretKey := viper.GetString("secret_key")
	externalUrl := viper.GetString("external_url")

	vaulto := a.VaultoAPI{}

	vaultoUrl := viper.GetString("vaulto_url")
	vaultoUser := viper.GetString("vaulto_user")
	vaultoPassword := viper.GetString("vaulto_user")

	vaulto.Init(vaultoUrl)

	vaulto.Login(vaultoUser, vaultoPassword)

	alfa := alfaex.AlfaEXAPI{}
	alfa.Init(publicKey, secretKey, externalUrl)

	operations := alfa.GetOperations()

	for _, o := range operations.Operation {
		amount := big.Float{}
		amount.SetString(o.Arguments.Amount)

		log.Println(o.Id, o.Arguments.Amount, amount.Text('f', 8))
		/*alfa.PutOperation(alfaex.WithdrawOperation{
			OperationID: o.Id,
			State: -5,
		})*/
	}

	alfa.GetOperations()

	wallets, err := vaulto.GetWalletsForAsset("ETH")
	log.Println("Wallets :", err, wallets)

	result, err := vaulto.CreateOrder("ETH", "0xa1894C90D2632850B6c20f217837e626628E5a15", 0.1, "New order")
	log.Println("OrderData :", err, result)

	orders, err := vaulto.GetOrders()
	log.Println("Orders :", err, orders)

	/*

		result, err := vaulto.Clear()
		log.Println("Clear :", err, result)

		result, err = vaulto.Register("user1", "pwd1")
		log.Println("Register : ", err, result)

		result, err = vaulto.Login("user1", "pwd1")
		log.Println("Login : ", err, result)

		result, err = vaulto.CreateAsset( "Ethereum", "ETH", 60, 18, 6)
		log.Println("Create asset : ", err, result)

		assets, err := vaulto.GetAssets()
		log.Println("Assets : ", assets)

		result, err = vaulto.CreateSeed("Seed1", "")
		log.Println("Create seed : ", err, result)

		seeds, err := vaulto.GetSeeds()
		log.Println("Seeds : ", err, seeds)

		result, err = vaulto.CreateWallet( "ETH wallet",  1, 1)
		log.Println("Create wallet : ", err, result)

		wallets, err := vaulto.GetWallets()
		log.Println("Wallets : ", wallets)

		result, err = vaulto.CreateAddress( "New address", 1)
		log.Println("Create address : ", err, result)
	*/
}
