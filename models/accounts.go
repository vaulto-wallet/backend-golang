package models

import (
	h "../helpers"
	"github.com/jinzhu/gorm"
	"github.com/xlzd/gotp"
	"time"
)

type OTPStatus int

const (
	OTPStatusUnknown  = 0
	OTPStatusNew      = 1
	OTPStatusDisabled = 2
	OTPStatusActive   = 3
)

type Account struct {
	gorm.Model
	Name       string    `json:"name"`
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private_key"`
	OTPKey     string    `json:"otp_key"`
	OTPStatus  OTPStatus `json:"otp_type"`
}

func (a *Account) VerifyTOTP(totpCode string) bool {
	if len(a.OTPKey) == 0 || a.OTPStatus != OTPStatusActive {
		return false
	}
	totp := gotp.NewDefaultTOTP(a.OTPKey)
	return totp.Verify(totpCode, int(time.Now().UTC().Unix()))
}

func (a *Account) Decrypt(password []byte, ciphertext []byte) ([]byte, error) {
	return h.DecryptWithRSA(password, []byte(a.PrivateKey), ciphertext)
}

func (a *Account) Encrypt(plaintext []byte) ([]byte, error) {
	return h.EncryptWithRSA([]byte(a.PublicKey), plaintext)
}
