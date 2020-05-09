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

	vaulto_user_2 := a.VaultoAPI{}
	vaulto_user_2.Init("http://localhost:8000/api")

	vaulto_auditor_1 := a.VaultoAPI{}
	vaulto_auditor_1.Init("http://localhost:8000/api")

	vaulto_auditor_2 := a.VaultoAPI{}
	vaulto_auditor_2.Init("http://localhost:8000/api")

	vaulto_fetcher := a.VaultoAPI{}
	vaulto_fetcher.Init("http://localhost:8000/api")

	vaulto_worker := a.VaultoAPI{}
	vaulto_worker.Init("http://localhost:8000/api")

	blockatlas := b.BlockAtlasAPI{}
	blockatlas.Init("http://localhost:8420/")

	fmt.Print(vaulto)

	result, err := vaulto.Clear("masterPassword")
	log.Println("Clear :", err, result)

	startResult, err := vaulto.Start("masterPassword")
	log.Println("Start :", err, startResult)

	user_id, err := vaulto.Register("user1", "userpwd1")
	log.Println("Register user 1: ", err, user_id)

	user_2_id, err := vaulto_user_2.Register("user2", "userpwd2")
	log.Println("Register user 2: ", err, user_2_id)

	fetcher_user_id, err := vaulto.Register("fetcher", "fetcher")
	log.Println("Register fetcher : ", err, fetcher_user_id)

	worker_user_id, err := vaulto.Register("worker", "worker")
	log.Println("Register worker : ", err, worker_user_id)

	auditor_1_user_id, err := vaulto.Register("auditor1", "auditorpwd1")
	log.Println("Register auditor 1 : ", err, auditor_1_user_id)

	auditor_2_user_id, err := vaulto.Register("auditor2", "auditorpwd1")
	log.Println("Register auditor 1 : ", err, auditor_2_user_id)

	result, err = vaulto.Login("user1", "userpwd1")
	log.Println("Login Owner : ", err, result)

	result, err = vaulto_user_2.Login("user2", "userpwd2")
	log.Println("Login User 2 : ", err, result)

	result, err = vaulto_auditor_1.Login("auditor1", "auditorpwd1")
	log.Println("Login Auditor 1 : ", err, result)

	result, err = vaulto_auditor_2.Login("auditor2", "auditorpwd2")
	log.Println("Login Auditor 1 : ", err, result)

	result, err = vaulto_fetcher.Login("fetcher", "fetcher")
	log.Println("Login Fetcher : ", err, result)

	result, err = vaulto_worker.Login("worker", "worker")
	log.Println("Login Worker : ", err, result)

	asset_id, err := vaulto.CreateAsset(m.AssetTypeBase, "Ethereum", "ETH", 60, 18, 6, "")
	log.Println("Create asset Ethereum: ", err, asset_id)

	asset_erc20_id, err := vaulto.CreateAsset(m.AssetTypeERC20, "ComBox", "CBP", 60, 4, 4, "0x3b423ba3799c26c42f0b389e97e1af8a1cca65e7")
	log.Println("Create asset ERC20: ", err, asset_erc20_id)

	asset_btct_id, err := vaulto.CreateAsset(m.AssetTypeBase, "BTCT Test", "BTCT", 9, 8, 8, "")
	log.Println("Create asset ERC20: ", err, asset_btct_id)

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

	wallet_btct_id, err := vaulto.CreateWallet("BTCT wallet", seed_id, asset_btct_id)
	log.Println("Create wallet : ", err, wallet_btct_id)

	wallet_btct_shared, err := vaulto.ShareWallet(wallet_btct_id, []uint{user_2_id}, []uint{auditor_1_user_id, auditor_2_user_id})
	log.Println("Share wallet BTCT : ", err, wallet_btct_shared)

	wallets, err := vaulto.GetWallets()
	log.Println("Wallets : ", wallets)

	wallets_fetcher, err := vaulto_fetcher.GetWallets()
	log.Println("Wallets Fetcher: ", wallets_fetcher)

	wallets_worker, err := vaulto_worker.GetWallets()
	log.Println("Wallets Fetcher: ", wallets_worker)

	wallets_eth, err := vaulto.GetWalletsForAsset("ETH")
	log.Println("Wallets for ETH : ", wallets_eth)

	wallets_cbp, err := vaulto.GetWalletsForAsset("CBP")
	log.Println("Wallets for CBP : ", wallets_cbp)

	wallets_btct, err := vaulto.GetWalletsForAsset("BTCT")
	log.Println("Wallets for CBP : ", wallets_btct)

	address_id, err := vaulto.CreateAddress("New address", wallets_eth[0].ID)
	log.Println("Create address : ", err, address_id)

	address_id, err = vaulto.CreateAddress("CBP New address", wallets_cbp[0].ID)
	log.Println("Create address : ", err, address_id)

	address_id, err = vaulto.CreateAddress("BTCT New address #1", wallets_btct[0].ID)
	log.Println("Create address : ", err, address_id)

	address_id, err = vaulto.CreateAddress("BTCT New address #2", wallets_btct[0].ID)
	log.Println("Create address : ", err, address_id)

	addresses, err := vaulto.GetAddressesForWallet(wallet_id)
	log.Println("Addresses ETH: ", err, addresses)

	addresses_erc20, err := vaulto.GetAddressesForWallet(wallet_erc20_id)
	log.Println("Addresses ERC20: ", err, addresses_erc20)

	addresses_btct, err := vaulto.GetAddressesForWallet(wallet_btct_id)
	log.Println("Addresses BTCT: ", err, addresses_btct)

	wallets, err = vaulto_fetcher.GetWalletsForAsset("ETH")
	log.Println("Fetcher ETH Wallets:", err, wallets)
	/*
		order_id, err := vaulto_fetcher.CreateOrder("ETH", "0x7d57DA3f1ED1Fd28Fc33338D6D07cf6c13a333c2", 0.01, "New order")
		log.Println("Fetcher OrderData :", err, order_id)

		order_erc20_id, err := vaulto_fetcher.CreateOrder("CBP", "0x7d57DA3f1ED1Fd28Fc33338D6D07cf6c13a333c2", 1, "ERC20 order")
		log.Println("Fetcher OrderData :", err, order_erc20_id)
	*/

	wallets_btct, err = vaulto_fetcher.GetWalletsForAsset("BTCT")
	log.Println("Fetcher BTCT Wallets:", err, wallets)

	order_btct_id, err := vaulto.CreateOrder(wallets_btct[0].ID, addresses_btct[0].Address, 0.00001, "BTCT order")
	log.Println("Fetcher OrderData :", err, order_btct_id)

	order_btct_confirm, err := vaulto.ConfirmOrder(order_btct_id)
	log.Println("User 1 confirms BTCT Order :", err, order_btct_confirm)

	orders, err := vaulto_fetcher.GetOrders()
	log.Println("Fetcher Orders :", err, orders)

	orders_fetcher_btct, err := vaulto_fetcher.GetOrdersForWallet(wallets_btct[0].ID)
	log.Println("Fetcher Orders for BTCT Wallets :", err, orders_fetcher_btct)

	orders_user_2_btct, err := vaulto_user_2.GetOrdersForWallet(wallets_btct[0].ID)
	log.Println("User 2 Orders for BTCT Wallets :", err, orders_user_2_btct)

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
