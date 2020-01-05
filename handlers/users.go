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
	UserName string
	jwt.Claims
}

func UserLogin(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	var user m.User
	var dbUser m.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		ThrowError(w, Error(BadRequest))
		return
	}

	db.First(&dbUser, "Username = ?", user.Username)
	if dbUser.ID == 0 {
		ThrowError(w, Error(NoUser))
		return
	}

	hash := sha256.New()
	hash.Write([]byte(user.Password))
	if hex.EncodeToString( hash.Sum(nil) ) != dbUser.Password {
		ThrowError(w, Error(IncorrectPassword))
		return
	}

	authToken := &AuthToken{UserName:user.Username}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authToken)

	tokenString, error := token.SignedString([]byte("cryptosecret"))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

func UserRegister(db *gorm.DB, w http.ResponseWriter, req *http.Request)  {
	var user m.User
	var dbUser m.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		log.Println(err)
		ThrowError(w, Error(BadRequest))
		return
	}

	db.First(&dbUser, "username = ?", user.Username)
	log.Println(dbUser)
	if dbUser.ID != 0 {
		ThrowError(w, Error(AlreadyRegistered))
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	dbUser = m.User{Username:user.Username, Password:hex.EncodeToString(h.Sum(nil))}
	db.Create(&dbUser)
	ReturnSuccess(w, Error(Success))
}

