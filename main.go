package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"lexipets/internal/pets"
	_ "lexipets/internal/pets"
	"lexipets/internal/users"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type CassandraSession struct {
	db *gocql.Session
}

func (session *CassandraSession) createUser(gc *gin.Context) {
	var auth users.New

	if err := gc.ShouldBindJSON(&auth); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userExists, reason, err := users.Exists(auth.Username, auth.Email, session.db)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userExists {
		gc.JSON(http.StatusBadRequest, gin.H{"error": reason})
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := users.User{
		Username:  auth.Username,
		Password:  string(passHash),
		Email:     auth.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	user, err = users.Create(session.db, user)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.IndentedJSON(http.StatusOK, user)
}

func (session *CassandraSession) login(gc *gin.Context) {
	var login users.Login

	if err := gc.ShouldBindJSON(&login); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := users.Authenticate(login, session.db)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Failed to generate token") {
			gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gc.Set("currentUser", user)
	gc.Set("token", token)
	gc.IndentedJSON(http.StatusOK, token)
}

func getOwnUser(gc *gin.Context) {
	user, exists := gc.Get("currentUser")
	if !exists {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to access User data. User not set or not logged in.")})
		return
	}

	gc.IndentedJSON(http.StatusOK, user)
}

func (session *CassandraSession) checkAuth(gc *gin.Context) {
	authHeader := gc.GetHeader("Authorization")

	if authHeader == "" {
		gc.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		gc.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		gc.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		gc.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		gc.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		gc.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		gc.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		gc.Abort()
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		gc.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		gc.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := users.Get("Id", claims["id"].(string), session.db)

	gc.Set("currentUser", user)

	gc.Next()
}

func Cassandra() (*gocql.Session, error) {
	session, err := cassandra()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to access Cassandra. Original error: %v", err.Error()))
	}

	return session, nil
}

func (session *CassandraSession) generatePet(gc *gin.Context) {
	var reqJson map[string]string
	err := gc.BindJSON(&reqJson)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind request json"})
		return
	}

	pet, err := pets.New(session.db, gc, reqJson["name"])

	if err != nil || pet.Name == "" {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate pet. Original error: %v", err)})
		return
	}
	gc.IndentedJSON(http.StatusOK, pet)
}

func (session *CassandraSession) savePet(gc *gin.Context) {
	var pet map[string]interface{}
	err := gc.BindJSON(&pet)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to bind JSON. Original error: %v. Failing JSON: %v. Failing request context: %v", err, pet, gc)})
		return
	}

	petJson, err := json.Marshal(pet)

	petId, err := pets.Save(session.db, gc, petJson)

	if err != nil || petId == "" {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save pet. Original error: %v. Failing JSON: %v", err, pet)})
		return
	}
	gc.IndentedJSON(http.StatusOK, petId)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	session, err := Cassandra()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	cSess := CassandraSession{db: session}

	defer cSess.db.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	router.POST("/users", cSess.createUser)
	router.POST("/users/login", cSess.login)
	router.GET("/users/me", cSess.checkAuth, getOwnUser)

	router.POST("/pets/generate", cSess.generatePet)
	router.POST("/pets", cSess.savePet)

	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}
