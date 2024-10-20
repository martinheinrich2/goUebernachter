package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// The Token struct holds the data for a token
type Token struct {
	ID          int
	UserID      int
	Token       string
	TokenExpiry time.Time
}

// The TokenModelInterface describes the methods that the actual TokenModel struct has.
type TokenModelInterface interface {
	InsertToken(userID int, token string) error
	DeleteToken(userID int) error
	GetToken(userID int) (Token, error)
	UserIdExists(userID int) (bool, error)
}

// The TokenModel struct which wraps a database connection pool.
type TokenModel struct {
	DB *sql.DB
}

// The InsertToken method inserts the user id, token and timestamp into table tokens.
func (m *TokenModel) InsertToken(userID int, token string) error {

	stmt := `INSERT INTO password_reset_tokens (user_id, token, token_expiry) VALUES(?, ?, DATETIME(CURRENT_TIMESTAMP, '+2 hours'))`
	_, err := m.DB.Exec(stmt, userID, token)
	if err != nil {
		fmt.Println("err != nil: ", err)
	}
	return nil
}

// The DeleteToken method deletes a token for a userid.
func (m *TokenModel) DeleteToken(userID int) error {
	stmt := `DELETE FROM password_reset_tokens WHERE user_id = ?`
	_, err := m.DB.Exec(stmt, userID)
	if err != nil {
		fmt.Println("err != nil: ", err)
	}
	return nil
}

// The GetToken method gets the token for a userId.
func (m *TokenModel) GetToken(userID int) (Token, error) {
	var token Token

	stmt := `SELECT user_id, token, token_expiry FROM password_reset_tokens WHERE user_id = ?`
	err := m.DB.QueryRow(stmt, userID).Scan(&token.UserID, &token.Token, &token.TokenExpiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Token{}, ErrNoRecord
		} else {
			return Token{}, err
		}
	}
	return token, nil
}

// The UserIdExists method checks if a user id exists in the database.
func (m *TokenModel) UserIdExists(userID int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS (SELECT true FROM password_reset_tokens WHERE user_id = ?)`
	err := m.DB.QueryRow(stmt, userID).Scan(&exists)
	return exists, err
}
