package users

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

func Check(fieldName string, fieldValue string, session *gocql.Session) (bool, error) {
	var count int
	statement := fmt.Sprintf(`SELECT COUNT(1) FROM lexipets.users WHERE %v = ? ALLOW FILTERING`, fieldName)
	err := session.Query(statement, fieldValue).Scan(&count)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error while checking if user exists. Original error: %v", err))
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

func Persist(user User, session *gocql.Session) (string, error) {
	newId := gocql.TimeUUID()
	err := session.Query(`INSERT INTO lexipets.users (id, username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`, newId, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Exec()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while persisting user. Original error: %v", err))
	}
	return newId.String(), nil
}

func Get(fieldName string, fieldValue string, session *gocql.Session) (User, error) {
	if fieldName == "password" {
		return User{}, errors.New(fmt.Sprintf("Bad request!"))
	}
	var (
		userId     string
		username   string
		email      string
		password   string
		created_at time.Time
		updated_at time.Time
		deleted_at time.Time
	)
	statement := fmt.Sprintf(`SELECT id, username, email, password, created_at, updated_at, deleted_at FROM lexipets.users WHERE %v = ? ALLOW FILTERING`, fieldName)
	err := session.Query(statement, fieldValue).Scan(&userId, &username, &email, &password, &created_at, &updated_at, &deleted_at)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Error while checking if user exists. Original error: %v", err))
	}

	if !deleted_at.IsZero() {
		return User{}, errors.New("User has been deleted.")
	}

	return User{Id: userId, Username: username, Email: email, Password: password, CreatedAt: created_at, UpdatedAt: updated_at, DeletedAt: deleted_at}, nil
}
