package models

import "github.com/jinzhu/gorm"

type Seed struct {
	Model
	Name    string
	Seed    string
	OwnerId uint
	Owner   *User
}

//Owner   User `gorm:"foreignkey:OwnerId, association_foreignkey:ID"`

type Seeds []Seed

func (s *Seed) Load(db *gorm.DB, seed_id uint) (err interface{}) {
	db.Set("gorm:auto_preload", true).Where("id = ?", seed_id).First(&s)
	if s.ID == 0 {
		return "Seed not found"
	}
	return nil

}
