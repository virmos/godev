package models

import (
	"errors"
)

// InsertHost inserts a host into the database
func (repo *TestRepository) InsertHost(h Host) (int, error) {
	return 1, nil
}

// BulkInsertHost inserts hosts into the database
func (repo *TestRepository) BulkInsertHost(hosts []Host) error {
	return nil
}

// InsertHost inserts a host into the database
func (repo *TestRepository) DeleteHost(ID int) error {
	return nil
}

// GetHostByID gets a host by id and returns Host
func (repo *TestRepository) GetHostByID(id int) (Host, error) {
	var host Host
	if id == 0 {
		return host, errors.New("id must be greater than 0")
	}
	return host, nil
}

// UpdateHost updates a host in the database
func (repo *TestRepository) UpdateHost(h Host) error {
	return nil
}

func (repo *TestRepository) GetAllServiceStatusCounts() (int, int, int, int, error) {
	return 1, 1, 1, 1, nil
}

// AllHosts returns a slice of hosts
func (repo *TestRepository) AllHosts() ([]Host, error) {
	var hosts []Host
	return hosts, nil
}

// UpdateHostServiceStatus updates the active status of a host service
func (repo *TestRepository) UpdateHostServiceStatus(hostID, serviceID, active int) error {
	return nil
}

// UpdateHostService updates a host service in the database
func (repo *TestRepository) UpdateHostService(hs HostService) error {
	return nil
}

// GetServicesByStatus returns all active services with a given status
func (repo *TestRepository) GetServicesByStatus(status string) ([]HostService, error) {
	var hs []HostService
	return hs, nil
}

// GetHostServiceByID gets a host service by id
func (repo *TestRepository) GetHostServiceByID(id int) (HostService, error) {
	var hs HostService
	return hs, nil
}

// GetServicesToMonitor gets all host services we want to monitor
func (repo *TestRepository) GetServicesToMonitor() ([]HostService, error) {
	var hs []HostService
	return hs, nil
}

// GetHostServiceByHostIDServiceID gets a host service by host id and service id
func (repo *TestRepository) GetHostServiceByHostIDServiceID(hostID, serviceID int) (HostService, error) {
	var hs HostService
	if serviceID > 3 {
		return hs, errors.New("services id in range 1-3(HTTP, HTTPS, SSL) respectively")
	}
	return hs, nil
}

// InsertEvent inserts an event into the database
func (repo *TestRepository) InsertEvent(e Event) error {
	return nil
}

// GetAllEvents gets all events
func (repo *TestRepository) GetAllEvents() ([]Event, error) {
	var events []Event
	return events, nil
}
