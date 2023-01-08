package models

import (
	"time"
)

type Repository interface {
	// users
	AllUsers() ([]*User, error)
	GetUserById(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	Authenticate(email, testPassword string) (int, string, error)
	InsertRememberMeToken(id int, token string) error
	DeleteToken(token string) error
	CheckForToken(id int, token string) bool
	InsertUser(u User) (int, error)
	UpdateUser(u User) error
	DeleteUser(id int) error
	UpdatePassword(id int, newPassword string) error

	// tokens
	InsertToken(t *Token, u User) error
	GetUserForToken(token string) (*User, error)
	GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error)

	// preferences
	AllPreferences() ([]Preference, error)
	SetSystemPref(name, value string) error
	UpdateSystemPref(name, value string) error
	InsertOrUpdateSitePreferences(pm map[string]string) error

	// hosts
	InsertHost(h Host) (int, error)
	DeleteHost(ID int) error
	BulkInsertHost(h []Host) error
	GetHostByID(id int) (Host, error)
	UpdateHost(h Host) error
	GetAllServiceStatusCounts() (int, int, int, int, error)
	AllHosts() ([]Host, error)
	UpdateHostServiceStatus(hostID, serviceID, active int) error
	UpdateHostService(hs HostService) error
	GetServicesByStatus(status string) ([]HostService, error)
	GetHostServiceByID(id int) (HostService, error)
	GetServicesToMonitor() ([]HostService, error)
	GetHostServiceByHostIDServiceID(hostID, serviceID int) (HostService, error)
	InsertEvent(e Event) error
	GetAllEvents() ([]Event, error)
}