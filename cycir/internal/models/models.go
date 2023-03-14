package models

import (
	"database/sql"
	"errors"
	"github.com/robfig/cron/v3"
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

var DB *sql.DB

type PostgresRepository struct {
	DB *sql.DB
}

type TestRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(pool *sql.DB) *PostgresRepository {
	DB = pool
	return &PostgresRepository{
		DB: pool,
	}
}

func NewTestRepository() *TestRepository {
	return &TestRepository{ }
}

// User model
type User struct {
	ID          int               `json:"id"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	UserActive  int               `json:"user_active"`
	AccessLevel int               `json:"access_level"`
	Email       string            `json:"email"`
	Password    []byte            `json:"password"`
	CreatedAt   time.Time         `json:"-"`
	UpdatedAt   time.Time         `json:"-"`
	DeletedAt   time.Time         `json:"deleted_at"`
	Preferences map[string]string `json:"preferences"`
}

// Preference model
type Preference struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Preference []byte    `json:"preference"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

// Setting for preference model
type Setting struct {
	SiteURL            string `json:"site_url"`
	NotifyName         string `json:"notify_name"`
	NotifyEmail        string `json:"notify_email"`
	SMTP_Server        string `json:"smtp_server"`
	SMTP_Port          string `json:"smtp_port"`
	SMTP_User          string `json:"smtp_user"`
	SMTP_Password      string `json:"smtp_password"`
	SMS_Enabled        string `json:"sms_enabled"`
	SMS_Provider       string `json:"sms_provider"`
	TWILIO_PhoneNumber string `json:"twilio_phone_number"`
	TWILIO_SID         string `json:"twilio_sid"`
	TWILIO_AuthToken   string `json:"twilio_auth_token"`
	SMTP_FromEmail     string `json:"smtp_from_email"`
	SMTP_FromName      string `json:"smtp_from_name"`
	NotifyViaSMS       string `json:"notify_via_sms"`
	NotifyViaEmail     string `json:"notify_via_email"`
	SMS_NotifyNumber   string `json:"sms_notify_number"`

	Action string `json:"action"`
}

// Host is the model for hosts
type Host struct {
	ID            int           `json:"id"`
	HostName      string        `json:"host_name"`
	CanonicalName string        `json:"canonical_name"`
	URL           string        `json:"url"`
	IP            string        `json:"ip"`
	IPV6          string        `json:"ipv6"`
	Location      string        `json:"location"`
	OS            string        `json:"os"`
	Active        int           `json:"active"`
	CreatedAt     time.Time     `json:"-"`
	UpdatedAt     time.Time     `json:"-"`
	HostServices  []HostService `json:"host_services"`

	Action string `json:"action"`
}

// Services is the model for services
type Services struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Active      int       `json:"active"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// HostService is the model for host services
type HostService struct {
	ID             int       `json:"id"`
	HostID         int       `json:"host_id"`
	ServiceID      int       `json:"service_id"`
	Active         int       `json:"active"`
	ScheduleNumber int       `json:"schedule_number"`
	ScheduleUnit   string    `json:"schedule_unit"`
	Status         string    `json:"status"`
	LastCheck      time.Time `json:"last_check"`
	LastMessage    string    `json:"last_message"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	Service        Services  `json:"service"`
	HostName       string    `json:"host_name"`
}

// Schedule model
type Schedule struct {
	ID            int          `json:"id"`
	EntryID       cron.EntryID `json:"entry_id"`
	Entry         cron.Entry   `json:"entry"`
	Host          string       `json:"host"`
	Service       string       `json:"service"`
	LastRunFromHS time.Time    `json:"last_run_from_hs"`
	HostServiceID int          `json:"host_service_id"`
	ScheduleText  string       `json:"schedule_text"`
}

// Event model
type Event struct {
	ID            int       `json:"id"`
	EventType     string    `json:"event_type"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	ServiceName   string    `json:"service_name"`
	HostName      string    `json:"host_name"`
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}
