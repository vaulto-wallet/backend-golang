package main

import (
	a "../vaultoapi"
	alfaex "./alfaex"
	"log"
	"math/big"
)

func main() {
	public_key := "7e041f7c036cbe256001b886fa95780f2e8f8c046d909cbdbdcbb25242709188"
	secret_key := "20ca2ad92dbf8db21f1cc9871ce5938f4c72a474142c7c5b95f5a47b18eeda3b9ac60b8f87e415f84c285ff252a4e74b1a71e1c0abcea97e40d4d369acde2acb"
	url := "http://167.99.243.191"
	vaulto := a.VaultoAPI{}
	vaulto.Init("http://localhost:8000/api")

	vaulto.Login("fetcher", "fetcher")

	alfa := alfaex.AlfaEXAPI{}
	alfa.Init(public_key, secret_key, url)

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
	log.Println("Order :", err, result)

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
