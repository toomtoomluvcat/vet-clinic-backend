package handler

import (
	"fmt"
	"my_postgres/app/database"
	"my_postgres/app/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type TransactionPost struct {
	PetID    string   `json:"PetID"`
	Symptom  string   `json:"Symptom"`
	Comment  string   `json:"Comment"`
	Expenses int      `json:"Expenses"`
	Medicine []string `json:"Medicine"`
}

type TransactionCommentPut struct {
	TransactionsID uint   `json:"TransactionsID"`
	Comment        string `json:"Comment"`
}

type TransactionResponse struct {
	TransactionsID uint
	OwnerName      string
	PetName        string
	Breed          string
	Symptom        string
	Comment        string
	EventTime      time.Time
	Expenses       int
	Medicine       pq.StringArray `gorm:"type:text[]"`
}

func CreateTransaction(c *gin.Context) {
	var req TransactionPost

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid JSON payload: %v", err)})
		fmt.Println(err)
		return
	}

	if req.PetID == "" || req.Comment == "" || req.Expenses == 0  {
		c.JSON(400, gin.H{"error": "payload must contain values"})
		return
	}
	PetIDInt, err := strconv.Atoi(req.PetID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Fail to convert Petid"})
	}
	Transaction := model.Transaction{
		PetID:     uint(PetIDInt),
		Symptom:   req.Symptom,
		Comment:   req.Comment,
		EventTime: time.Now(),
		Expenses:  req.Expenses,
		Medicine:  pq.StringArray(req.Medicine),
	}

	if err := database.DB.Create(&Transaction).Error; err != nil {
		c.JSON(500, gin.H{"error": "Fail to Insert Datatype"})
		return
	}

	c.JSON(201, gin.H{"message": "success fully to create Transaction"})
}

func GetAllTransaction(c *gin.Context) {
	var transaction []TransactionResponse

	err := database.DB.Table("transactions").Select("transactions.transactions_id,owners.owner_name,pets.pet_name,transactions.symptom,transactions.event_time,transactions.expenses,transactions.medicine").
		Joins("JOIN pets ON transactions.pet_id = pets.pet_id").Joins("JOIN owners ON pets.owner_id = owners.owner_id").Find(&transaction).Error

	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to load:%v", err)})
		return
	}
	c.JSON(200, gin.H{"transaction": transaction})
}

func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	type Transaction struct {
		TransactionsID uint
		Symptom        string
		Comment        string
		EventTime      time.Time
		Expenses       int
		Medicine       pq.StringArray `gorm:"type:text[]"`
	}

	type OwnerAndPet struct {
		OwnerName   string
		PetName     string
		PetID       int
		Breed       string
		Transaction []Transaction
	}

	type tempInfo struct {
		OwnerName string
		PetID     int
		PetName   string
		Breed     string
	}

	var Temp tempInfo
	if err := database.DB.Table("pets").
		Select("owners.owner_name,pets.pet_id, pets.pet_name, pets.breed").
		Joins("JOIN owners ON pets.owner_id = owners.owner_id").
		Where("pets.pet_id = ?", id).
		Scan(&Temp).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("error when get owner: %v", err)})
		return
	}

	var transactions []Transaction
	if err := database.DB.Table("transactions").
		Select("transactions.transactions_id, transactions.symptom, transactions.comment, transactions.event_time, transactions.expenses, transactions.medicine").
		Where("transactions.pet_id = ?", id).
		Scan(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("error when get transactions: %v", err)})
		return
	}

	ownerAndPet := OwnerAndPet{
		OwnerName:   Temp.OwnerName,
		PetName:     Temp.PetName,
		PetID:       Temp.PetID,
		Breed:       Temp.Breed,
		Transaction: transactions,
	}

	c.JSON(200, ownerAndPet)
}

func PutTransactionByID(c *gin.Context) {
	type Req struct {
		PetID          uint 
		TransactionsID uint `gorm:"primaryKey"`
		Symptom        string
		Comment        string
		EventTime	  time.Time
		Expenses       int
		Medicine       pq.StringArray `gorm:"type:text[]"`
	}

	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid JSON payload:%v",err)})
		fmt.Println(err)
		return
	}

	fmt.Println(req.TransactionsID)

	if err := database.DB.Save(&model.Transaction{TransactionsID: req.TransactionsID,EventTime: req.EventTime, Symptom: req.Symptom,
		Comment: req.Comment, Expenses: req.Expenses, Medicine: req.Medicine,PetID: req.PetID,}).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to update transaction at: %v", err)})
		return
	}
	c.Status(204)

}

func DelTransactionByID(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&model.Transaction{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to Delete transaction:%v", err)})
		return
	}
	c.Status(204)
}
