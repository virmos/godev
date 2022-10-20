package models

import (
	"time"
)

func (repo *TestRepository) InsertToken(t *Token, u User) error {
	return nil
}

func (repo *TestRepository) GetUserForToken(token string) (*User, error) {
	var user User
	return &user, nil
}

func (repo *TestRepository) GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	var token Token
	return &token, nil
}