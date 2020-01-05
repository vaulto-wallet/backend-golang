package models

import "github.com/jinzhu/gorm"

type AssetType int
const (
	Unknown = 0
	AssetBase = 1
	AssetERC20 = 2
)

func (a AssetType) String() string {
	asset_type_text  := [...]string{
		"Unknown",
		"Basic",
		"ERC20",
	}
	return asset_type_text[a]
}

type Asset struct {
	gorm.Model
	Name string
	CoinIndex int
	Symbol string `gorm:"unique_index"`
	Type AssetType
	Decimals int
	Rounding int
	Address string
}
