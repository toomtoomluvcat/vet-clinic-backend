package main

import (
	"github.com/gin-gonic/gin"
	"github.com/toomtoomluvcat/vet-clinic-backend/app/database"
	"github.com/toomtoomluvcat/vet-clinic-backend/app/handler"
)


func CORSMiddleware() gin.HandlerFunc {

    return func(c *gin.Context) { 
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH") 
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
	r := gin.Default()

	r.Use(CORSMiddleware())

	database.ConnectDB()

	r.POST("/owner", handler.CreateOwnerWithPet)
	r.PUT("/owner", handler.EditOwenr)
	r.GET("/owner", handler.GetOwnerAll)
	r.DELETE("/owner/:id", handler.DeleteOwnerByID)
	r.GET("/owner/:id", handler.GetOwnerByID)
	r.GET("/pet/:id",handler.GetPetByID)
	r.POST("/pet", handler.AddNewPet)
	r.PUT("/pet", handler.PutPet)
	r.DELETE("/pet/:id",handler.DelPetByID)
	r.POST("/transaction", handler.CreateTransaction)
	r.PUT("/transaction", handler.PutTransactionByID)
	r.GET("/transaction",handler.GetAllTransaction)
	r.GET("transaction/:id",handler.GetTransactionByID)
	
	r.Run(":8081")
}
