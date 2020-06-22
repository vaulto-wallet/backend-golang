package models

import (
	"encoding/hex"
	"github.com/jinzhu/gorm"
)

type ConfigRecord struct {
	Model
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c ConfigRecord) Get(db *gorm.DB) ConfigRecord {
	db.First(&c, "name = ?", c.Name)
	return c
}

func (c ConfigRecord) GetHex(db *gorm.DB) []byte {
	db.First(&c, "name = ?", c.Name)
	if c.ID == 0 {
		return nil
	}
	data, err := hex.DecodeString(c.Value)
	if err != nil {
		return nil
	}
	return data
}

func (c ConfigRecord) Set(db *gorm.DB) {
	record := new(ConfigRecord)
	db.First(&record, "name = ?", c.Name)
	if record.ID == 0 {
		record.Name = c.Name
		record.Value = c.Value
		db.Save(record)
	} else {
		record.Value = c.Value
		db.Update(record)
	}
}
