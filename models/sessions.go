package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
)

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

// func (ss SessionService) Create(userID int) (*Session, error) {
// 	bytesPerToken := ss.BytesPerToken
// 	if bytesPerToken < 32 {
// 		bytesPerToken = 32
// 	}
// 	token, err := String(bytesPerToken)
// 	if err != nil {
// 		return nil, fmt.Errorf("create session: %w", err)
// 	}
// 	tokenHash := ss.hash(token)
// 	session := Session{
// 		UserID: userID,
// 		Token: token,
// 		TokenHash: tokenHash,
// 	}

// }

// hash used to hash the token
func (ss SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
