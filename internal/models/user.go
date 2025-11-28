package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	FirstName    string
	LastName   string
	CanVote      bool
	CreatedAt    time.Time
}

type UserService struct {
	DB *sql.DB
}

func (u *UserService) Create(email, password, firstName, secondName string) (*User, error) {
	email = strings.ToLower(email)
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating the hash: %v", err)
	}

	NewUser := User{
		Email: email,
		PasswordHash: string(hashedPassword),
		FirstName: firstName,
		LastName: secondName,
		CanVote: true,
		CreatedAt: time.Now(),
	}

	row := u.DB.QueryRow(`INSERT INTO users (email, password_hash, first_name, last_name, can_vote, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, NewUser.Email, 
	NewUser.PasswordHash, 
	NewUser.FirstName, 
	NewUser.LastName, 
	NewUser.CanVote,
	NewUser.CreatedAt)

	err = row.Scan(&NewUser.ID)
	if err != nil {
		return nil, fmt.Errorf("could insert: %v", err)
	}

	return &NewUser, nil
}
 
func (u *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)

	user := User{
		Email: email,
	}

	row := u.DB.QueryRow(`SELECT password_hash, first_name, last_name FROM users WHERE email = $1`, user.Email)
	err := row.Scan(&user.PasswordHash, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, fmt.Errorf("could not find: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %v", err)
	}

	return &user, nil
}