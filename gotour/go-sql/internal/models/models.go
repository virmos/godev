package models

import (
	"errors"
	"time"
)

var (
	// ErrNoRecord no record found in database error
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials invalid username/password error
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail duplicate email error
	ErrDuplicateEmail = errors.New("models: duplicate email")
	// ErrInactiveAccount inactive account error
	ErrInactiveAccount = errors.New("models: Inactive Account")
)

// User model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	UserActive  int
	AccessLevel int
	Email       string
	Password    []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Preferences map[string]string
}

// Preference model
type Preference struct {
	ID         int
	Name       string
	Preference []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Department model
type Department struct {
	ID         					int
	Name								string
	ParentId   					int
	OrderSequence       string
	TreeLevel						int
	Children						[]*Department
}

// Employee model
type Employee struct {
	ID         		int
	DepartmentId  int
	UserName      string
	FirstName     string
	LastName      string
	Email       	string
	CreatedAt  		time.Time
	UpdatedAt  		time.Time
}

// Alert model
type Alert struct {
	Interval  		string
	Type					string
	Count					int
}
