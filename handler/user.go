package handler

import (
	"errors"
	"fmt"

	"github.com/toomtoomluvcat/vet-clinic-backend/app/database"
	"github.com/toomtoomluvcat/vet-clinic-backend/app/model"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type OwnerPost struct {
	Name    string    `json:"Name"`
	Address string    `json:"Address"`
	Tel     string    `json:"Tel"`
	Pets    []PetPost `json:"Pets"`
}

type PetPost struct {
	Name    string   `json:"Name"`
	Age     uint     `json:"Age"`
	Weight  int      `json:"Weight"`
	Disease []string `json:"Disease"`
	Breed   string   `json:"Breed"`
}

type PetPostWithID struct {
	OwnerID int    
	PetName    string  
	Age     int    
	Weight  int     
	Disease []string 
	Breed   string   
}

type OwnerPut struct {
	OwnerID   uint
	OwnerName string
	Address   string
	Tel       string
}

func CreateOwnerWithPet(c *gin.Context) {
	var req OwnerPost
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid JSON payload:%v", err)})
		return
	}

	if req.Name == "" || req.Address == "" || req.Tel == "" {
		c.JSON(400, gin.H{"error": "Request must contain values"})
		return
	}

	owner := model.Owner{
		OwnerName: req.Name,
		Tel:       req.Tel,
		Address:   req.Address,
	}

	for i, p := range req.Pets {
		if p.Name == "" || p.Age == 0 || p.Weight == 0 || p.Breed == "" {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Pet data not complete in pet index:%d", i)})
			return
		}

		pet := model.Pet{
			PetName: p.Name,
			Breed:   p.Breed,
			Age:     p.Age,
			Weight:  p.Weight,
			Disease: pq.StringArray(p.Disease),
		}

		owner.Pets = append(owner.Pets, pet)
	}

	if err := database.DB.Create(&owner).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to insert Owner and Pet: %v", err)})
		return
	}

	c.JSON(201, gin.H{"message": "successfully to create"})
}


func EditOwenr(c *gin.Context) {
	var req OwnerPut

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}
	owner := model.Owner{
		OwnerID:   req.OwnerID,
		OwnerName: req.OwnerName,
		Address:   req.Address,
		Tel:       req.Tel,
	}

	if err := database.DB.Save(&owner).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to Update Owner:%v", err)})
		return
	}

	c.JSON(204, gin.H{"message": "Successfully to update Owner"})

}

func GetOwnerAll(c *gin.Context) {
	var Owners []model.Owner
	err := database.DB.Select("OwnerID", "OwnerName").Find(&Owners).Error
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to GET owner: %v", err)})
		return
	}
	c.JSON(200, Owners)

}
func GetOwnerByID(c *gin.Context) {
	id := c.Param("id")

	type PetResponse struct {
		PetID   uint
		PetName string
		Age     uint
		Weight  int
		Breed   string
		Disease pq.StringArray
		Transactions []model.Transaction `json:"transactions"`
	}

	type OwnerResponse struct {
		OwnerID   uint          `json:"ownerID"`
		OwnerName string        `json:"ownerName"`
		Address   string        `json:"address"`
		Tel       string        `json:"tel"`
		Pets      []PetResponse `json:"pets"`
	}

	var owner model.Owner
	err := database.DB.Where("owner_id = ?", id).First(&owner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "owner not found"})
		} else {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to get owner: %v", err)})
		}
		return
	}
		var pets []model.Pet
		if err := database.DB.Preload("Records").Where("owner_id = ?", id).Find(&pets).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(404, gin.H{"error": "Not found any pet"})
			} else {
				c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to load pets: %v", err)})
			}
			return
		}
		

	response := OwnerResponse{

		OwnerID:   owner.OwnerID,
		OwnerName: owner.OwnerName,
		Address:   owner.Address,
		Tel:       owner.Tel,
		Pets:      make([]PetResponse, 0, len(pets))}

	for _, p := range pets {
		pet := PetResponse{
			PetID:   p.PetID,
			PetName: p.PetName,
			Weight:  p.Weight,
			Age:     p.Age,
			Breed:   p.Breed,
			Disease: pq.StringArray(p.Disease),
			Transactions: p.Records,
		}
		response.Pets = append(response.Pets, pet)
	}

	c.JSON(200, response)
}

func DeleteOwnerByID(c *gin.Context) {
	id := c.Param("id")

	// ดึง Owner เพื่อเช็กว่ามีอยู่ไหม
	var owner model.Owner
	if err := database.DB.First(&owner, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Owner not found"})
		} else {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to find Owner: %v", err)})
		}
		return
	}

	// ดึง Pets ของ Owner
	var pets []model.Pet
	if err := database.DB.Where("owner_id = ?", id).Find(&pets).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to find Pets: %v", err)})
		return
	}

	// ลบ Transactions ของ Pet แต่ละตัว
	for _, pet := range pets {
		if err := database.DB.Where("pet_id = ?", pet.PetID).Delete(&model.Transaction{}).Error; err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to delete Transactions of pet_id %d: %v", pet.PetID, err)})
			return
		}
	}

	// ลบ Pets ทั้งหมดของ Owner
	if err := database.DB.Where("owner_id = ?", id).Delete(&model.Pet{}).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to delete Pets: %v", err)})
		return
	}

	// ลบ Owner
	if err := database.DB.Delete(&model.Owner{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Fail to delete Owner: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully deleted owner and related data"})
}
