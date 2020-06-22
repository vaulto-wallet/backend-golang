package models

type User struct {
	Model
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AccountId uint      `json:"-"`
	Account   Account   `json:"account"`
	Wallets   []*Wallet `json:"-" gorm:"many2many:wallet_owners"`
}

type Users []*User
