package handler

import (
	"fmt"
	"my_postgres/app/database"
	"my_postgres/app/model"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func GetPetByID(c *gin.Context) {
	id := c.Param("id")

	type PetResponse struct {
		PetID     int
		PetName   string
		Age       int
		Weight    int
		Breed     string
		Disease   pq.StringArray
		OwnerID   int
		OwnerName string
	}

	var petResponse PetResponse
	if err := database.DB.Table("pets").Joins("JOIN owners ON owners.owner_id=pets.pet_id").Where("pets.pet_id = ?", id).Scan(&petResponse).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to load err as:%v", err)})
		return
	}
	c.JSON(200, petResponse)
}

func AddNewPet(c *gin.Context) {
	var req PetPostWithID
	if err := c.BindJSON(&req); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Invalid JSON payload: %v", err)})
		fmt.Println(err)
		return
	}
	fmt.Println(req)
	pet := model.Pet{
		PetName: req.PetName,
		Age:     uint(req.Age),
		Weight:  req.Weight,
		Breed:   req.Breed,
		Disease: pq.StringArray(req.Disease),
		OwnerID: uint(req.OwnerID),
	}

	if err := database.DB.Create(&pet).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to create: %v", err)})
		return
	}
	c.JSON(201, gin.H{"message": "Successfully to create"})
}

func PutPet(c *gin.Context) {
	type Req struct {
		PetID   uint
		PetName string
		Age     int
		Weight  int
		Breed   string
	}

	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}
	if err := database.DB.Model(&model.Pet{}).
		Where("pet_id = ?", req.PetID).
		Updates(model.Pet{
			PetName: req.PetName,
			Age:     uint(req.Age),
			Weight:  req.Weight,
			Breed:   req.Breed,
		}).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to update: %v", err)})
		return
	}

}

func DelPetByID(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	database.DB.Where("pet_id = ?", id).Delete(&model.Transaction{})
	if err := database.DB.Where("pet_id = ?", id).Delete(&model.Pet{}).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to Delete pet:%v", err)})
		return
	}

	c.Status(204)
}
