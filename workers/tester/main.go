package main

import (
	b "../../api/blockatlas"
	a "../../api/vaulto"
	m "../../models"
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

	asset_id, err := vaulto.CreateAsset(m.AssetTypeBase, "Ethereum", "ETH", 60, 18, 6, "")
	log.Println("Create asset Ethereum: ", err, asset_id)

	asset_erc20_id, err := vaulto.CreateAsset(m.AssetTypeERC20, "ComBox", "CBP", 60, 4, 4, "0x3b423ba3799c26c42f0b389e97e1af8a1cca65e7")
	log.Println("Create asset ERC20: ", err, asset_erc20_id)

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

	wallet_id, err := vaulto.CreateWallet("ETH wallet", seed_id, asset_id)
	log.Println("Create wallet : ", err, wallet_id)

	wallet_erc20_id, err := vaulto.CreateWallet("CBP wallet", seed_id, asset_erc20_id)
	log.Println("Create wallet : ", err, wallet_erc20_id)

	wallets, err := vaulto.GetWallets()
	log.Println("Wallets : ", wallets)

	wallets_eth, err := vaulto.GetWalletsForAsset("ETH")
	log.Println("Wallets for ETH : ", wallets_eth)

	wallets_cbp, err := vaulto.GetWalletsForAsset("CBP")
	log.Println("Wallets for CBP : ", wallets_cbp)

	address_id, err := vaulto.CreateAddress("New address", wallets_eth[0].ID)
	log.Println("Create address : ", err, address_id)

	address_id, err = vaulto.CreateAddress("CBP New address", wallets_cbp[0].ID)
	log.Println("Create address : ", err, address_id)

	addresses, err := vaulto.GetAddressesForWallet(wallet_id)
	log.Println("Addresses ETH: ", err, addresses)

	addresses_erc20, err := vaulto.GetAddressesForWallet(wallet_erc20_id)
	log.Println("Addresses ERC20: ", err, addresses_erc20)

	wallets, err = vaulto_fetcher.GetWalletsForAsset("ETH")
	log.Println("Fetcher Wallets :", err, wallets)

	order_id, err := vaulto_fetcher.CreateOrder("ETH", "0x7d57DA3f1ED1Fd28Fc33338D6D07cf6c13a333c2", 0.01, "New order")
	log.Println("Fetcher OrderData :", err, order_id)

	order_erc20_id, err := vaulto_fetcher.CreateOrder("CBP", "0x7d57DA3f1ED1Fd28Fc33338D6D07cf6c13a333c2", 1, "ERC20 order")
	log.Println("Fetcher OrderData :", err, order_erc20_id)

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
