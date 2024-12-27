package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type NewUser struct {
	Email    string
	Password string
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

// CreateUser creates a user by hashing the password, storing it and the email in the DB
// and returning a user
func (us *UserService) CreateUser(nu NewUser) (*User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, fmt.Errorf("create user: password hashing: %w", err)
	}
	passwordHash := string(hashedBytes)
	user := User{
		Email:        nu.Email,
		PasswordHash: passwordHash,
	}
	row := us.DB.QueryRow(`
	INSERT INTO users (email, password_hash)
	VALUES ($1, $2) RETURNING id`, nu.Email, passwordHash)
	err = row.Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return &User{}, fmt.Errorf("create user: email already taken: %w", err)
			}
		}
		return &User{}, fmt.Errorf("create user: query: %w", err)
	}
	return &user, nil
}

func (us *UserService) FindUser(userID int) (*User, error) {
	user := User{
		ID: userID,
	}

	row := us.DB.QueryRow(`SELECT email, password_hash FROM users WHERE id=$1`, userID)
	err := row.Scan(user.Email, user.PasswordHash)
	if err != nil {
		// this means there is no such user with id=userID
		return &User{}, fmt.Errorf("find user: no such user: %w", err)
	}
	return &user, nil
}
