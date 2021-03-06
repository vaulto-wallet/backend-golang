package models

import (
	h "../helpers"
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
	Model
	Name       string    `json:"name,omitempty"`
	Email      string    `json:"email,omitempty"`
	PublicKey  string    `json:"public_key,omitempty"`
	PrivateKey string    `json:"private_key,omitempty"`
	OTPKey     string    `json:"otp_key,omitempty"`
	OTPStatus  OTPStatus `json:"otp_type,omitempty"`
}

type Accounts []Account

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
