package main

import (
	b "../api/blockatlas"
	a "../api/vaulto"
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

	result, err = vaulto.Register("user1", "pwd1")
	log.Println("Register user : ", err, result)

	result, err = vaulto.Register("fetcher", "fetcher")
	log.Println("Register fetcher : ", err, result)

	result, err = vaulto.Register("worker", "worker")
	log.Println("Register sender : ", err, result)

	result, err = vaulto.Login("user1", "pwd1")
	log.Println("Login : ", err, result)

	result, err = vaulto_fetcher.Login("fetcher", "fetcher")
	log.Println("Login : ", err, result)

	result, err = vaulto.CreateAsset("Ethereum", "ETH", 60, 18, 6)
	log.Println("Create asset : ", err, result)

	assets, err := vaulto.GetAssets()
	log.Println("Assets : ", assets)

	result, err = vaulto.CreateSeed("Seed1", "")
	log.Println("Create seed : ", err, result)

	seeds, err := vaulto.GetSeeds()
	log.Println("Seeds : ", err, seeds)

	result, err = vaulto.CreateWallet("ETH wallet", 1, 1)
	log.Println("Create wallet : ", err, result)

	wallets, err := vaulto.GetWallets()
	log.Println("Wallets : ", wallets)

	wallets_eth, err := vaulto.GetWalletsForAsset("ETH")
	log.Println("Wallets for ETH : ", wallets_eth)

	result, err = vaulto.CreateAddress("New address", 1)
	log.Println("Create address : ", err, result)

	wallets, err = vaulto_fetcher.GetWalletsForAsset("ETH")
	log.Println("Fetcher Wallets :", err, wallets)

	result, err = vaulto_fetcher.CreateOrder("ETH", "0xa1894C90D2632850B6c20f217837e626628E5a15", 0.1, "New order")
	log.Println("Fetcher Order :", err, result)

	orders, err := vaulto_fetcher.GetOrders()
	log.Println("Fetcher Orders :", err, orders)

	txs, err := blockatlas.GetTXs("etherscan", "0x278F5F53156Be78bFE665D5354d40c539ca02ef3")
	log.Println("BlockAtlas TXs :", err, txs[0])

}
