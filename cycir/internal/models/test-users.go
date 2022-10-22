package models

// AllUsers returns all users
func (repo *TestRepository) AllUsers() ([]*User, error) {
	var u []*User
	return u, nil
}

// GetUserById returns a user by id
func (repo *TestRepository) GetUserById(id int) (User, error) {
	var u User
	return u, nil
}

// GetUserByEmail gets a user by email address
func (repo *TestRepository) GetUserByEmail(email string) (User, error) {
	var u User
	return u, nil
}

// Authenticate authenticates
func (repo *TestRepository) Authenticate(email, testPassword string) (int, string, error) {
	if email == "admin@example.com" {
		return 1, "", nil
	}
	return 0, "", ErrInvalidCredentials
}

// InsertRememberMeToken inserts a remember me token into remember_tokens for a user
func (repo *TestRepository) InsertRememberMeToken(id int, token string) error {
	return nil
}

// DeleteToken deletes a remember me token
func (repo *TestRepository) DeleteToken(token string) error {
	return nil
}

// CheckForToken checks for a valid remember me token
func (repo *TestRepository) CheckForToken(id int, token string) bool {
	if token == "xyz" {
		return true
	}
	return false
}

// Insert method to add a new record to the users table.
func (repo *TestRepository) InsertUser(u User) (int, error) {
	return 1, nil
}

// UpdateUser updates a user by id
func (repo *TestRepository) UpdateUser(u User) error {
	return nil
}

// DeleteUser sets a user to deleted by populating deleted_at value
func (repo *TestRepository) DeleteUser(id int) error {
	return nil
}

// UpdatePassword resets a password
func (repo *TestRepository) UpdatePassword(id int, newPassword string) error {
	return nil
}
