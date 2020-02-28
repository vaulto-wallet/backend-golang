package handlers

import (
	m "../models"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func CreateAccount(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
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

	ReturnResult(w, true)
}
