package main

import (
	"../../api/vaulto"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

func processOrders(api *vaulto.VaultoAPI) {
	orders, err := api.GetOrders()
	log.Println("Orders :", err, orders)
	for i, o := range orders {
		log.Println("Order ", i, o.Amount)
	}
	for i, a := range Assets {
		log.Println("Assets ", i, a.Symbol)
	}

}

var Assets vaulto.AssetsResponse

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
