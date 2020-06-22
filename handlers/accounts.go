package handlers

import (
	m "../models"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func SetAccount(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	var reqAccount m.Account
	if err := json.NewDecoder(req.Body).Decode(&reqAccount); err != nil {
		log.Println(err)
		ReturnError(w, Error(BadRequest))
		return
	}
	if len(reqAccount.Email) > 0 {
		db.Model(&user.Account).Update("email", reqAccount.Email)
	}
	if len(reqAccount.Name) > 0 {
		db.Model(&user.Account).Update("name", reqAccount.Name)
	}

	ReturnResult(w, true)
}

func GetAccount(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)
	ReturnResult(w, user)
}
