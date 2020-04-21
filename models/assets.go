package models

import (
	"github.com/jinzhu/gorm"
	"math/big"
)

type AssetType int

const (
	AssetTypeUnknown = 0
	AssetTypeBase    = 1
	AssetTypeERC20   = 2
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
	Address   string
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

func (ar *Assets) GetBasicAsset(asset_id uint) (ret *Asset) {
	a := ar.Get(asset_id)
	if a == nil {
		return nil
	}
	if a.Type == AssetTypeERC20 {
		ba := ar.Find("ETH")
		return ba
	}
	return a
}

func (ar *Assets) Find(symbol string) (ret *Asset) {
	for _, s := range ([]Asset)(*ar) {
		if s.Symbol == symbol {
			return &s
		}
	}
	return nil
}

func (ar *Assets) Get(asset_id uint) (ret *Asset) {
	for _, s := range ([]Asset)(*ar) {
		if s.ID == asset_id {
			return &s
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
