package handlers

import (
	m "../models"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func CreateAccount(db *gorm.DB, w http.ResponseWriter, req *http.Request)  {
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
	ReturnBool(w, true)
}
