package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
	"github.com/virmos/django"
	"golang.org/x/crypto/bcrypt"
)

// User is the type for a user
type User struct {
	ID        string    `db:"_id,omitempty" json:"_id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Table returns the table name associated with this model in the database
func (u *User) Table() string {
	return "users"
}

func (u *User) Validate(validator *django.Validation) {
	validator.Check(u.Name != "", "name", "Name must be provided")
	validator.Check(u.Email != "", "email", "Email must be provided")
	validator.IsEmail("email", u.Email)
}

// GetAll returns a slice of all users
func (u *User) GetAll() ([]*User, error) {
	collection := upper.Collection(u.Table())

	var all []*User

	res := collection.Find().OrderBy("name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetByEmail gets one user, by email
func (u *User) GetByEmail(email string) (*User, error) {
	var theUser User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"email =": email})
	err := res.One(&theUser)
	if err != nil {
		return nil, err
	}

	return &theUser, nil
}

// Get gets one user by id
func (u *User) Get(id string) (*User, error) {
	var theUser User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"_id =": id})

	err := res.One(&theUser)
	if err != nil {
		return nil, err
	}

	return &theUser, nil
}

// Update updates a user record in the database
func (u *User) Update(theUser User) error {
	theUser.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"_id =": theUser.ID})
	err := res.Update(&theUser)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a user by id
func (u *User) Delete(id string) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"_id =": id})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

// Insert inserts a new user, and returns the newly inserted id
func (u *User) Insert(theUser User) (string, error) {
	newHash, err := bcrypt.GenerateFromPassword([]byte(theUser.Password), 12)
	if err != nil {
		return "0", err
	}

	theUser.CreatedAt = time.Now()
	theUser.UpdatedAt = time.Now()
	theUser.Password = string(newHash)

	collection := upper.Collection(u.Table())
	_, err = collection.Insert(theUser)
	if err != nil {
		return "0", err
	}

	return theUser.ID, nil
}

// PasswordMatches verifies a supplied password against the hash stored in the database.
// It returns true if valid, and false if the password does not match, or if there is an
// error. Note that an error is only returned if something goes wrong (since an invalid password
// is not an error -- it's just the wrong password))
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			// some kind of error occurred
			return false, err
		}
	}

	return true, nil
}
