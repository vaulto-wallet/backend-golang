package main

import (
	a "../vaultoapi"
	"fmt"
	"log"
)

func main() {
	vaulto := a.VaultoAPI{}
	vaulto.Init("http://localhost:8000/api")
	fmt.Print(vaulto)

	result, err := vaulto.Clear()
	log.Println("Clear :", err, result)

	result, err = vaulto.Register("user1", "pwd1")
	log.Println("Register : ", err, result)

	result, err = vaulto.Login("user1", "pwd1")
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

	result, err = vaulto.CreateAddress("New address", 1)
	log.Println("Create address : ", err, result)

}
