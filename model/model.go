package model

import (
	"time"

	"github.com/lib/pq"
)

type Owner struct {
	OwnerID   uint `gorm:"primaryKey"`
	OwnerName string
	Address   string
	Tel       string
	Pets      []Pet `gorm:"foreignKey:OwnerID"`
}

type Pet struct {
	OwnerID uint
	Owner   Owner
	PetID   uint `gorm:"primaryKey"`
	PetName string
	Age     uint
	Weight  int
	Breed   string
	Disease pq.StringArray `gorm:"type:text[]"`
	Records []Transaction  `gorm:"foreignKey:PetID"`
}
type Transaction struct {
	TransactionsID uint `gorm:"primaryKey"`
	PetID          uint
	Pet            Pet
	Symptom        string
	Comment        string
	EventTime      time.Time
	Expenses       int
	Medicine       pq.StringArray `gorm:"type:text[]"`
}
