package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lexipets/internal/pets"
	_ "lexipets/internal/pets"
	"net/http"
)

func generatePet(gc *gin.Context) {
	session, err := cassandra()
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access Cassandra"})
		return
	}
	defer session.Close()

	var reqJson map[string]string
	err = gc.BindJSON(&reqJson)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind request json"})
		return
	}

	pet, err := pets.New(session, gc, reqJson["name"])

	if err != nil || pet.Name == "" {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate pet. Original error: %v", err)})
		return
	}
	gc.IndentedJSON(http.StatusOK, pet)
}

func savePet(gc *gin.Context) {
	session, err := cassandra()
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access Cassandra"})
		return
	}
	defer session.Close()

	var pet map[string]interface{}
	err = gc.BindJSON(&pet)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to bind JSON. Original error: %v. Failing JSON: %v. Failing request context: %v", err, pet, gc)})
		return
	}

	petJson, err := json.Marshal(pet)

	petId, err := pets.Save(session, gc, petJson)

	if err != nil || petId == "" {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save pet. Original error: %v. Failing JSON: %v", err, pet)})
		return
	}
	gc.IndentedJSON(http.StatusOK, petId)
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))
	router.POST("/pets/generate", generatePet)
	router.POST("/pets", savePet)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
