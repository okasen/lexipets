package main

import (
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

	pet, err := pets.New(session, gc)

	if err != nil || pet.Name == "" {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	gc.IndentedJSON(http.StatusOK, pet)
}

func main() {
	router := gin.Default()
	router.GET("/pets/generate", generatePet)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
