package users

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func Exists(username string, email string, session *gocql.Session) (bool, string, error) {
	exists, err := Check("username", username, session)
	if err != nil {
		return false, "error occured", err
	}
	if exists {
		return true, "username is taken", nil
	}
	if email != "" {
		exists, err = Check("email", email, session)
		if err != nil {
			return false, "error occured", err
		}
		if exists {
			return true, "Email is in use", nil
		}
	}
	return false, "", nil
}

func Create(session *gocql.Session, user User) (User, error) {
	userId, err := Persist(user, session)
	if err != nil {
		return User{}, err
	}
	user.Id = userId
	return user, nil
}

func Authenticate(login Login, session *gocql.Session) (User, string, error) {
	user, err := Get("username", login.Username, session)
	if err != nil {
		return User{}, "", errors.New(fmt.Sprintf("Cannot authenticate. Original Error: %v", err.Error()))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return User{}, "", errors.New(fmt.Sprintf("Invalid Username or Password", err.Error()))
	}

	generate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generate.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return User{}, "", errors.New(fmt.Sprintf("Failed to generate token", err.Error()))
	}

	return user, token, nil
}
