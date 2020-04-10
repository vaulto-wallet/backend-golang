package main

import (
	b "../../api/blockatlas"
	a "../../api/vaulto"
	"fmt"
	"log"
)

func main() {
	vaulto := a.VaultoAPI{}
	vaulto.Init("http://localhost:8000/api")

	vaulto_fetcher := a.VaultoAPI{}
	vaulto_fetcher.Init("http://localhost:8000/api")

	blockatlas := b.BlockAtlasAPI{}
	blockatlas.Init("http://localhost:8420/")

	fmt.Print(vaulto)

	result, err := vaulto.Clear()
	log.Println("Clear :", err, result)

	user_id, err := vaulto.Register("user1", "pwd1")
	log.Println("Register user : ", err, user_id)

	user_id, err = vaulto.Register("fetcher", "fetcher")
	log.Println("Register fetcher : ", err, user_id)

	user_id, err = vaulto.Register("worker", "worker")
	log.Println("Register sender : ", err, user_id)

	result, err = vaulto.Login("user1", "pwd1")
	log.Println("Login : ", err, result)

	result, err = vaulto_fetcher.Login("fetcher", "fetcher")
	log.Println("Login : ", err, result)

	asset_id, err := vaulto.CreateAsset("Ethereum", "ETH", 60, 18, 6)
	log.Println("Create asset : ", err, asset_id)

	assets, err := vaulto.GetAssets()
	log.Println("Assets : ", assets)

	eth := assets.Find("ETH")
	intEth := eth.ToBigInt(0.1)
	log.Println("ETH Float to BigInt : ", intEth)
	floatEth := eth.ToFloat(intEth)
	log.Println("ETH BigInt to Float: ", floatEth)

	seed_id, err := vaulto.CreateSeed("Seed1", "orange okay much equip pond cushion ask hover bar shove ceiling have")
	log.Println("Create seed : ", err, seed_id)

	seeds, err := vaulto.GetSeeds()
	log.Println("Seeds : ", err, seeds)

	wallet_id, err := vaulto.CreateWallet("ETH wallet", 1, 1)
	log.Println("Create wallet : ", err, wallet_id)

	wallets, err := vaulto.GetWallets()
	log.Println("Wallets : ", wallets)

	wallets_eth, err := vaulto.GetWalletsForAsset("ETH")
	log.Println("Wallets for ETH : ", wallets_eth)

	address_id, err := vaulto.CreateAddress("New address", 1)
	log.Println("Create address : ", err, address_id)

	address_id, err = vaulto.CreateAddress("New address", 1)
	log.Println("Create address : ", err, address_id)

	addresses, err := vaulto.GetAddressesForWallet(1)
	log.Println("Addresses : ", err, addresses)

	wallets, err = vaulto_fetcher.GetWalletsForAsset("ETH")
	log.Println("Fetcher Wallets :", err, wallets)

	order_id, err := vaulto_fetcher.CreateOrder("ETH", "0x7d57DA3f1ED1Fd28Fc33338D6D07cf6c13a333c2", 0.01, "New order")
	log.Println("Fetcher OrderData :", err, order_id)

	orders, err := vaulto_fetcher.GetOrders()
	log.Println("Fetcher Orders :", err, orders)

	txs, err := blockatlas.GetTXs("ETH", "0x278F5F53156Be78bFE665D5354d40c539ca02ef3")
	log.Println("BlockAtlas TXs :", err, txs)

	price, err := blockatlas.GasPrice("ETH")
	log.Println("BlockAtlas GasPrice :", err, price)

	//var input e.SigningInput

	input := map[string]interface{}{}

	input["to"] = "0x278F5F53156Be78bFE665D5354d40c539ca02ef3"
	input["value"] = 96000000000

	gasRequired, err := blockatlas.EstimateGas("etherscan", input)
	log.Println("BlockAtlas EstimatedGas :", err, gasRequired)

}
