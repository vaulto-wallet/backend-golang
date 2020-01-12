package models

import "github.com/jinzhu/gorm"

type Seed struct {
	gorm.Model
	Name string
	Seed string
}

type PrivateKeyType int

const (
	PrivateKeyTypeUnknown = 0
	PrivateKeyTypeRoot = 1
	PrivateKeyTypeSingle = 2
	PrivateKeyTypeMulti = 3
)


type PrivateKey struct{
	gorm.Model
	Name string
	Seed Seed
	PublicKey string
	Type PrivateKeyType
	Owner User
	Asset Asset
	NetworkType string
}