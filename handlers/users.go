package handlers

import (
	m "../models"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
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

	db.First(&dbUser, "Username = ?", user.Username)
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

	h := sha256.New()
	h.Write([]byte(user.Password))
	dbUser = m.User{Username: user.Username, Password: hex.EncodeToString(h.Sum(nil))}
	db.Create(&dbUser)
	ReturnResult(w, true)
}

func Clear(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	db.DropTableIfExists("assets")
	db.DropTableIfExists("users")
	db.DropTableIfExists("seeds")
	db.DropTableIfExists("wallets")
	db.DropTableIfExists("orders")
	db.DropTableIfExists("addresses")
	db.AutoMigrate(&m.Asset{})
	db.AutoMigrate(&m.User{})
	db.AutoMigrate(&m.Seed{}).AddForeignKey("owner", "users(id)", "RESTRICT", "RESTRICT")
	db.AutoMigrate(&m.Wallet{})
	db.AutoMigrate(&m.Address{})
	db.AutoMigrate(&m.Order{})
	//db.Model(&m.Seed{}).AddForeignKey("owner", "users(id)", "RESTRICT", "RESTRICT")

	ReturnResult(w, true)
}
