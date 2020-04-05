package models

import (
	"github.com/jinzhu/gorm"
	"math/big"
)

type AssetType int

const (
	AssetUnknown = 0
	AssetBase    = 1
	AssetERC20   = 2
)

func (a AssetType) String() string {
	asset_type_text := [...]string{
		"Unknown",
		"Basic",
		"ERC20",
	}
	return asset_type_text[a]
}

type Asset struct {
	gorm.Model
	Name      string
	CoinIndex int
	Symbol    string `gorm:"unique_index"`
	Type      AssetType
	Decimals  int
	Rounding  int
}

type Assets []Asset

func (a *Asset) Get(db *gorm.DB, asset_id uint) (err interface{}) {
	db.First(&a, "ID = ?", asset_id)
	if a.ID == 0 {
		return "Asset not found"
	}
	return nil
}

func (ar *Assets) Find(symbol string) (ret *Asset) {
	for i, s := range ([]Asset)(*ar) {
		if s.Symbol == symbol {
			return &([]Asset)(*ar)[i]
		}
	}
	return nil
}

func (a *Asset) ToBigInt(value float64) (ret *big.Int) {
	e := new(big.Int)
	r := new(big.Float)
	e.Exp(big.NewInt(10), big.NewInt(int64(a.Decimals)), nil)
	r.SetFloat64(value).Mul(r, new(big.Float).SetInt(e))
	ret, _ = r.Int(nil)
	return
}

func (a *Asset) ToFloat(value *big.Int) (ret float64) {
	e := new(big.Int)
	r := new(big.Float)
	e.Exp(big.NewInt(10), big.NewInt(int64(a.Decimals)), nil)
	r.SetInt(value)
	ret, _ = new(big.Float).Quo(r, new(big.Float).SetInt(e)).Float64()
	return
}
