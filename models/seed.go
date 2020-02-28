package models

import "github.com/jinzhu/gorm"

type Seed struct {
	gorm.Model
	Name    string
	Seed    string `json:"-"`
	OwnerID int
	Owner   User `json:"-"`
}

type Seeds []Seed

func (s *Seed) Get(db *gorm.DB, seed_id int) (err interface{}) {
	db.First(&s, "ID = ?", seed_id)
	if s.ID == 0 {
		return "Seed not found"
	}
	return nil
}
