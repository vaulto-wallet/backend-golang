package handlers

import (
	m "../models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
)

func CreateSeed(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}
	var r interface{}
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
	}
	rm := r.(map[string]interface{})

	seed := make([]byte, 32)
	rand.Read(seed)

	seedName := "New seed"
	seedParam, exists := rm["name"]
	if exists {
		seedName = seedParam.(string)
	}

	seedHex := hex.EncodeToString(seed)

	dbSeed := m.Seed{
		Name:  seedName,
		Seed:  seedHex,
		Owner: dbUser,
	}
	db.Create(&dbSeed)

	ReturnResult(w, true)
}

func GetSeeds(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	var seeds m.Seeds

	db.Where("owner_id = ? ", dbUser.ID).Find(&seeds)

	res, err := json.Marshal(seeds)
	if err != nil {
		ReturnResult(w, seeds)
		return
	}

	fmt.Println((string)(res))

	ReturnResult(w, seeds)
}
