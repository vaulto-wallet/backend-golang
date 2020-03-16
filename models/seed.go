package models

import "github.com/jinzhu/gorm"

type Seed struct {
	gorm.Model
	Name    string
	Seed    string
	OwnerID int
	Owner   User
}

//Owner   User `gorm:"foreignkey:OwnerID, association_foreignkey:ID"`

type Seeds []Seed

func (s *Seed) Get(db *gorm.DB, seed_id uint) (err interface{}) {
	db.Set("gorm:auto_preload", true).Where("ID = ?", seed_id).First(&s)
	if s.ID == 0 {
		return "Seed not found"
	}
	return nil

}
