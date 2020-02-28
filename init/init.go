package main

import (
	m "../models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	db, _ := gorm.Open("sqlite3", "test.db")

	db.AutoMigrate(&m.Asset{})
	db.AutoMigrate(&m.Seed{})
	db.AutoMigrate(&m.PrivateKey{})

	asset_eth := m.Asset{Name: "Ethereum", CoinIndex: 60, Symbol: "ETH", Type: 1, Decimals: 18, Rounding: 6}
	db.Create(&asset_eth)

}
