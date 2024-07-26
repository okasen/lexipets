package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lexipets/internal/pets"
	_ "lexipets/internal/pets"
	"net/http"
)

func generatePet(gc *gin.Context) {
	session, err := cassandra()
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	defer session.Close()

	var reqJson map[string]string
	err = gc.BindJSON(&reqJson)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	pet, err := pets.New(session, gc, reqJson["name"])

	if err != nil || pet.Name == "" {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	gc.IndentedJSON(http.StatusOK, pet)
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))
	router.POST("/pets/generate", generatePet)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
