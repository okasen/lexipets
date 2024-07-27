package users

import "time"

type User struct {
	Id        string    `cql:"id" json:"id"`
	Username  string    `cql:"username" json:"username"`
	Email     string    `cql:"email" json:"email"`
	Password  string    `cql:"password" json:"password"`
	CreatedAt time.Time `cql:"created_at" json:"created_at"`
	UpdatedAt time.Time `cql:"updated_at" json:"updated_at"`
	DeletedAt time.Time `cql:"deleted_at" json:"deleted_at,omitempty"`
}

type New struct {
	Username string `binding:"required" json:"username"`
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

type Login struct {
	Username string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
}
