package models

func (repo *TestRepository) InsertToken(t *Token, u User) error {
	return nil
}

func (repo *TestRepository) GetUserForToken(token string) (*User, error) {
	var user User
	return &user, nil
}
