package handlers

import (
	h "../helpers"
	m "../models"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/xlzd/gotp"
	"log"
	"net/http"
)

type AuthToken struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func UserLogin(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	var user m.User
	var dbUser m.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
	}

	db.Preloads("Accounts").First(&dbUser, "Username = ?", user.Username)
	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}

	hash := sha256.New()
	hash.Write([]byte(user.Password))
	if hex.EncodeToString(hash.Sum(nil)) != dbUser.Password {
		ReturnError(w, Error(IncorrectPassword))
		return
	}

	authToken := &AuthToken{user.Username,
		jwt.StandardClaims{
			ExpiresAt: 0,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authToken)

	tokenString, error := token.SignedString([]byte("cryptosecret"))
	if error != nil {
		fmt.Println(error)
	}
	ReturnResult(w, tokenString)
}

func UserRegister(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	var user m.User
	var dbUser m.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		log.Println(err)
		ReturnError(w, Error(BadRequest))
		return
	}

	db.First(&dbUser, "username = ?", user.Username)
	log.Println(dbUser)
	if dbUser.ID != 0 {
		ReturnError(w, Error(AlreadyRegistered))
		return
	}

	passwordHash := sha256.Sum256([]byte(user.Password))

	priv, pub := h.GenerateRSAKey(passwordHash[:])

	a := m.Account{
		Name:       "",
		PublicKey:  hex.EncodeToString(pub),
		PrivateKey: hex.EncodeToString(priv),
		OTPKey:     gotp.RandomSecret(12),
		OTPStatus:  m.OTPStatusNew,
	}

	db.Create(&a)

	dbUser = m.User{Username: user.Username, Password: hex.EncodeToString(passwordHash[:]), Account: a}
	db.Create(&dbUser)

	ReturnResult(w, dbUser.ID)
}

func Start(db *gorm.DB, w http.ResponseWriter, req *http.Request) []byte {
	InitReq := struct {
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(req.Body).Decode(&InitReq); err != nil {
		log.Println(err)
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Cannot decode request")
		return []byte{}
	}

	testString := []byte("TestStringMessage")

	privkey := m.ConfigRecord{Name: "PrivateKey"}.Get(db)
	pubkey := m.ConfigRecord{Name: "PublicKey"}.Get(db)

	encryptedString, err := h.EncryptWithRSA([]byte(pubkey.Value), testString)
	if err != nil {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Cannot initialize protected storage ")
		return nil
	}
	decryptedString, err := h.DecryptWithRSA([]byte(InitReq.Password), []byte(privkey.Value), encryptedString)

	if string(decryptedString) != string(testString) {
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Incorrect password")
		return nil
	}
	ReturnResult(w, true)

	return []byte(InitReq.Password)
}

func Clear(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	db.DropTableIfExists("assets")
	db.DropTableIfExists("users")
	db.DropTableIfExists("seeds")
	db.DropTableIfExists("wallets")
	db.DropTableIfExists("orders")
	db.DropTableIfExists("order_destinations")
	db.DropTableIfExists("orders_dbs")
	db.DropTableIfExists("addresses")
	db.DropTableIfExists("accounts")
	db.DropTableIfExists("config_records")
	db.DropTableIfExists("transactions")
	db.DropTableIfExists("transaction_addresses")
	db.DropTableIfExists("transaction_assets")
	db.DropTableIfExists("transaction_orders")
	db.DropTableIfExists("transaction_wallets")
	db.DropTableIfExists("wallet_owners")
	db.DropTableIfExists("wallet_auditors")
	db.DropTableIfExists("firewall_rules")
	db.DropTableIfExists("confirmations")
	db.AutoMigrate(&m.Asset{})
	db.AutoMigrate(&m.User{})
	db.AutoMigrate(&m.Account{})
	db.AutoMigrate(&m.Seed{})
	db.AutoMigrate(&m.Wallet{})
	db.AutoMigrate(&m.Address{})
	db.AutoMigrate(&m.Order{})
	db.AutoMigrate(&m.OrderDestination{})
	db.AutoMigrate(&m.Transaction{})
	db.AutoMigrate(&m.ConfigRecord{})
	db.AutoMigrate(&m.FirewallRule{})
	db.AutoMigrate(&m.Confirmation{})

	InitReq := struct {
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(req.Body).Decode(&InitReq); err != nil {
		log.Println(err)
		ReturnErrorWithStatusString(w, Error(BadRequest), 400, "Cannot decode request")
		return
	}

	priv, pub := h.GenerateRSAKey([]byte(InitReq.Password))

	m.ConfigRecord{Name: "PrivateKey", Value: string(priv)}.Set(db)
	m.ConfigRecord{Name: "PublicKey", Value: string(pub)}.Set(db)

	ReturnResult(w, true)
}
