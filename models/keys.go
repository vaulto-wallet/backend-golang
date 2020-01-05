package models

import "github.com/jinzhu/gorm"

type Seed struct {
	gorm.Model
	Name string
	Seed string
}

type PrivateKeyType int

const (
	Unknown = 0
	Root = 1
	Single = 2
	Multi = 3
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