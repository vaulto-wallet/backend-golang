package models

type Confirmation struct {
	Model
	OrderId uint  `json:"order_id"`
	Order   Order `json:"order"`
	UserId  uint  `json:"user_id"`
	User    User  `json:"user"`
}

type Confirmations []Confirmation
