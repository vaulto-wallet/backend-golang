package handlers

import (
	m "../models"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CreateAsset(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}

	var r m.Asset
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		ReturnError(w, Error(BadRequest))
	}

	db.Create(&r)
	ReturnResult(w, r.ID)
}

func GetAssets(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	var assets m.Assets

	db.Find(&assets)

	res, err := json.Marshal(assets)
	if err != nil {
		ReturnResult(w, assets)
		return
	}

	fmt.Println((string)(res))
	ReturnResult(w, assets)
}
