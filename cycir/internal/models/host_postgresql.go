package models

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// InsertHost inserts a host into the database
func (repo *PostgresRepository) InsertHost(h Host) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into hosts (host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	var newID int

	err := repo.DB.QueryRowContext(ctx, query,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.Location,
		h.OS,
		h.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println(err)
		return newID, err
	}

	// get preference map
	var scheduleAmount []byte
	var scheduleUnit []byte
	query = "SELECT preference FROM preferences WHERE name='check_interval_amount'"

	err = repo.DB.QueryRowContext(ctx, query).Scan(&scheduleAmount)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	query = "SELECT preference FROM preferences WHERE name='check_interval_unit'"
	err = repo.DB.QueryRowContext(ctx, query).Scan(&scheduleUnit)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// add host services and set to inactive
	query = `select id from services`
	serviceRows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer serviceRows.Close()

	for serviceRows.Next() {
		var svcID int
		err := serviceRows.Scan(&svcID)
		if err != nil {
			log.Println(err)
			return 0, err
		}

		stmt := `
			insert into host_services 
		    	(host_id, service_id, active, schedule_number, schedule_unit,
				status, created_at, updated_at) values ($1, $2, 0, $3, $4, 'pending', $5, $6)`

		_, err = repo.DB.ExecContext(ctx, stmt, newID, svcID, string(scheduleAmount), string(scheduleUnit), time.Now(), time.Now())
		if err != nil {
			return newID, err
		}
	}

	return newID, nil
}

// BulkInsertHost inserts hosts into the database
func (repo *PostgresRepository) BulkInsertHost(hosts []Host) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// drop foreign key reference to host table: host_services_hosts_id_fk
	// no need to drop the foreign key refernece to services table: host_services_services_id_fk
	stmt := "ALTER TABLE host_services DROP CONSTRAINT host_services_hosts_id_fk;"
	_, err := repo.DB.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return err
	}

	// truncate hosts table
	stmt = "TRUNCATE hosts;"
	_, err = repo.DB.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return err
	}

	// truncate host_services table
	stmt = "TRUNCATE host_services;"
	_, err = repo.DB.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("truncated")

	// get preference map
	var scheduleAmount []byte
	var scheduleUnit []byte
	query := "SELECT preference FROM preferences WHERE name='check_interval_amount'"

	err = repo.DB.QueryRowContext(ctx, query).Scan(&scheduleAmount)
	if err != nil {
		log.Println(err)
		return err
	}

	query = "SELECT preference FROM preferences WHERE name='check_interval_unit'"
	err = repo.DB.QueryRowContext(ctx, query).Scan(&scheduleUnit)
	if err != nil {
		log.Println(err)
		return err
	}
	
	valueStrings := make([]string, 0, len(hosts))
	valueArgs := make([]interface{}, 0, len(hosts) * 10)

	for i, h := range hosts {
		a := []int{10*i+1, 10*i+2, 10*i+3, 10*i+4, 10*i+5, 10*i+6, 10*i+7, 10*i+8, 10*i+9, 10*i+10}
		arg := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), ", $"), "[]")
		if i == (len(hosts) - 1) {
			arg = "($" + arg + ");"
		} else {
			arg = "($" + arg + "), "
		}

		valueStrings = append(valueStrings, arg)
		valueArgs = append(valueArgs, h.HostName)
		valueArgs = append(valueArgs, h.CanonicalName)
		valueArgs = append(valueArgs, h.URL)
		valueArgs = append(valueArgs, h.IP)
		valueArgs = append(valueArgs, h.IPV6)
		valueArgs = append(valueArgs, h.Location)
		valueArgs = append(valueArgs, h.OS)
		valueArgs = append(valueArgs, h.Active)
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, time.Now())
	}

	log.Println("inserting hosts")

	stmt = fmt.Sprintf("insert into hosts (host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at) VALUES %s", 
                        strings.Join(valueStrings, ""))
	_, err = repo.DB.Exec(stmt, valueArgs...)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("finish inserting")

	query = `select id from hosts`
	hostRows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return err
	}
	defer hostRows.Close()
	var hostIDs []int
	for hostRows.Next() {
		var hostID int
		err := hostRows.Scan(&hostID)
		if err != nil {
			log.Println(err)
			return err
		}
		hostIDs = append(hostIDs, hostID)
	}

	// get all service ids
	query = `select id from services`
	serviceRows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return err
	}
	defer serviceRows.Close()

	var svcIDs []int
	for serviceRows.Next() {
		var svcID int
		err := serviceRows.Scan(&svcID)
		if err != nil {
			log.Println(err)
			return err
		}
		svcIDs = append(svcIDs, svcID)
	}
	log.Println(hostIDs)

	// populate into host_services
	for _, hID := range hostIDs {
		for _, svcID := range svcIDs {
			stmt = `
			insert into host_services 
					(host_id, service_id, active, schedule_number, schedule_unit,
				status, created_at, updated_at) values ($1, $2, 0, $3, $4, 'pending', $5, $6)`

			_, err = repo.DB.ExecContext(ctx, stmt, hID, svcID, string(scheduleAmount), string(scheduleUnit), time.Now(), time.Now())
			if err != nil {
				return err
			}
		}
	}

	// add foreign key reference to hosts table
	stmt = "ALTER TABLE host_services ADD CONSTRAINT host_services_hosts_id_fk FOREIGN KEY (host_id) REFERENCES hosts (id);"
	_, err = repo.DB.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// InsertHost inserts a host into the database
func (repo *PostgresRepository) DeleteHost(ID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "ALTER TABLE host_services DROP CONSTRAINT host_services_hosts_id_fk;"
	_, err := repo.DB.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt = `
		delete from hosts where id = $1`

	_, err = repo.DB.ExecContext(ctx, stmt, ID)
	if err != nil {
		return err
	}
		
	if err != nil {
		log.Println(err)
		return err
	}
	stmt = `
		delete from host_services where host_id = $1`
	_, err = repo.DB.ExecContext(ctx, stmt, ID)
	if err != nil {
		return err
	}

	// add foreign key reference to hosts table
	stmt = "ALTER TABLE host_services ADD CONSTRAINT host_services_hosts_id_fk FOREIGN KEY (host_id) REFERENCES hosts (id);"
	_, err = repo.DB.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return err
	}
	
	return nil
}

// GetHostByID gets a host by id and returns Host
func (repo *PostgresRepository) GetHostByID(id int) (Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select 
			id, host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at
		from 
			hosts where id = $1`

	row := repo.DB.QueryRowContext(ctx, query, id)

	var h Host

	err := row.Scan(
		&h.ID,
		&h.HostName,
		&h.CanonicalName,
		&h.URL,
		&h.IP,
		&h.IPV6,
		&h.Location,
		&h.OS,
		&h.Active,
		&h.CreatedAt,
		&h.UpdatedAt,
	)

	if err != nil {
		return h, err
	}

	// get all services for host
	query = `
			select 
				hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, 
				hs.last_check, hs.status, hs.created_at, hs.updated_at,
				s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at, hs.last_message
			from 
				host_services hs
				left join services s on (s.id = hs.service_id)
			where
				host_id = $1
			order by s.service_name`

	rows, err := repo.DB.QueryContext(ctx, query, h.ID)
	if err != nil {
		log.Println(err)
		return h, err
	}
	defer rows.Close()

	var hostServices []HostService
	for rows.Next() {
		var hs HostService
		err := rows.Scan(
			&hs.ID,
			&hs.HostID,
			&hs.ServiceID,
			&hs.Active,
			&hs.ScheduleNumber,
			&hs.ScheduleUnit,
			&hs.LastCheck,
			&hs.Status,
			&hs.CreatedAt,
			&hs.UpdatedAt,
			&hs.Service.ID,
			&hs.Service.ServiceName,
			&hs.Service.Active,
			&hs.Service.Icon,
			&hs.Service.CreatedAt,
			&hs.Service.UpdatedAt,
			&hs.LastMessage,
		)
		if err != nil {
			log.Println(err)
			return h, err
		}
		hostServices = append(hostServices, hs)
	}

	h.HostServices = hostServices

	return h, nil
}

// UpdateHost updates a host in the database
func (repo *PostgresRepository) UpdateHost(h Host) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
			update 
    			hosts 
			set 
			    host_name = $1, canonical_name = $2, url = $3, ip = $4, ipv6 = $5, os = $6,
				active = $7, location = $8, updated_at = $9 
			where 
			    id = $10`

	_, err := repo.DB.ExecContext(ctx, stmt,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.OS,
		h.Active,
		h.Location,
		time.Now(),
		h.ID,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (repo *PostgresRepository) GetAllServiceStatusCounts() (int, int, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		(select count(id) from host_services where active = 1 and status = 'pending') as pending,
		(select count(id) from host_services where active = 1 and status = 'healthy') as healthy,
		(select count(id) from host_services where active = 1 and status = 'warning') as warning,
		(select count(id) from host_services where active = 1 and status = 'problem') as problem`

	var pending, healthy, warning, problem int

	row := repo.DB.QueryRowContext(ctx, query)
	err := row.Scan(
		&pending,
		&healthy,
		&warning,
		&problem,
	)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return pending, healthy, warning, problem, nil
}

// AllHosts returns a slice of hosts
func (repo *PostgresRepository) AllHosts() ([]Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			select 
       			id, host_name, canonical_name, url, ip, ipv6, location, os,
				active, created_at, updated_at 
			from 
			     hosts 
			order by 
				host_name`

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []Host

	for rows.Next() {
		var h Host
		err = rows.Scan(
			&h.ID,
			&h.HostName,
			&h.CanonicalName,
			&h.URL,
			&h.IP,
			&h.IPV6,
			&h.Location,
			&h.OS,
			&h.Active,
			&h.CreatedAt,
			&h.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// get all services for host
		serviceQuery := `
			select 
				hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, 
				hs.last_check, hs.status, hs.created_at, hs.updated_at,
				s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at, hs.last_message
			from 
				host_services hs
				left join services s on (s.id = hs.service_id)
			where
				host_id = $1 
			and hs.active = 1`

		serviceRows, err := repo.DB.QueryContext(ctx, serviceQuery, h.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var hostServices []HostService

		defer serviceRows.Close()
		for serviceRows.Next() {
			var hs HostService
			err = serviceRows.Scan(
				&hs.ID,
				&hs.HostID,
				&hs.ServiceID,
				&hs.Active,
				&hs.ScheduleNumber,
				&hs.ScheduleUnit,
				&hs.LastCheck,
				&hs.Status,
				&hs.CreatedAt,
				&hs.UpdatedAt,
				&hs.Service.ID,
				&hs.Service.ServiceName,
				&hs.Service.Active,
				&hs.Service.Icon,
				&hs.Service.CreatedAt,
				&hs.Service.UpdatedAt,
				&hs.LastMessage,
			)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			hostServices = append(hostServices, hs)
		}
		h.HostServices = hostServices
		hosts = append(hosts, h)
		
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return hosts, nil
}

// UpdateHostServiceStatus updates the active status of a host service
func (repo *PostgresRepository) UpdateHostServiceStatus(hostID, serviceID, active int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update host_services set active = $1 where host_id = $2 and service_id = $3`

	_, err := repo.DB.ExecContext(ctx, stmt, active, hostID, serviceID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateHostService updates a host service in the database
func (repo *PostgresRepository) UpdateHostService(hs HostService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update 
    			host_services set
    				host_id = $1, service_id = $2, active = $3,
				  	schedule_number = $4, schedule_unit = $5, 
				  	last_check = $6, status = $7, updated_at = $8, last_message = $9
				where 
					id = $10`

	_, err := repo.DB.ExecContext(ctx, stmt,
		hs.HostID,
		hs.ServiceID,
		hs.Active,
		hs.ScheduleNumber,
		hs.ScheduleUnit,
		hs.LastCheck,
		hs.Status,
		hs.UpdatedAt,
		hs.LastMessage,
		hs.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetServicesByStatus returns all active services with a given status
func (repo *PostgresRepository) GetServicesByStatus(status string) ([]HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select
			hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit,
			hs.last_check, hs.status, hs.created_at, hs.updated_at,
			h.host_name, s.service_name, hs.last_message
		from
			host_services hs
			left join hosts h on (hs.host_id = h.id)
			left join services s on (hs.service_id = s.id)
		where
			status = $1
			and hs.active = 1
		order by 
			 host_name, service_name`

	var services []HostService

	rows, err := repo.DB.QueryContext(ctx, query, status)
	if err != nil {
		return services, err
	}
	defer rows.Close()

	for rows.Next() {
		var h HostService

		err := rows.Scan(
			&h.ID,
			&h.HostID,
			&h.ServiceID,
			&h.Active,
			&h.ScheduleNumber,
			&h.ScheduleUnit,
			&h.LastCheck,
			&h.Status,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.HostName,
			&h.Service.ServiceName,
			&h.LastMessage,
		)
		if err != nil {
			return nil, err
		}

		services = append(services, h)
	}

	return services, nil
}

// GetHostServiceByID gets a host service by id
func (repo *PostgresRepository) GetHostServiceByID(id int) (HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number,
			hs.schedule_unit, hs.last_check, hs.status, hs.created_at, hs.updated_at,
			s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at, h.host_name,
		    hs.last_message

		from host_services hs
		left join services s on (hs.service_id = s.id)
		left join hosts h on (hs.host_id = h.id)

		where hs.id = $1
`

	var hs HostService

	row := repo.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&hs.ID,
		&hs.HostID,
		&hs.ServiceID,
		&hs.Active,
		&hs.ScheduleNumber,
		&hs.ScheduleUnit,
		&hs.LastCheck,
		&hs.Status,
		&hs.CreatedAt,
		&hs.UpdatedAt,
		&hs.Service.ID,
		&hs.Service.ServiceName,
		&hs.Service.Active,
		&hs.Service.Icon,
		&hs.Service.CreatedAt,
		&hs.Service.UpdatedAt,
		&hs.HostName,
		&hs.LastMessage,
	)

	if err != nil {
		log.Println(err)
		return hs, err
	}

	return hs, nil
}

// GetServicesToMonitor gets all host services we want to monitor
func (repo *PostgresRepository) GetServicesToMonitor() ([]HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number,
			hs.schedule_unit, hs.last_check, hs.status, hs.created_at, hs.updated_at,
			s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at,
			h.host_name, hs.last_message
		from 
		     host_services hs
			left join services s on (hs.service_id = s.id)
			left join hosts h on (h.id = hs.host_id)
		where
			h.active = 1
			and hs.active = 1`

	var services []HostService

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var h HostService
		err := rows.Scan(
			&h.ID,
			&h.HostID,
			&h.ServiceID,
			&h.Active,
			&h.ScheduleNumber,
			&h.ScheduleUnit,
			&h.LastCheck,
			&h.Status,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.Service.ID,
			&h.Service.ServiceName,
			&h.Service.Active,
			&h.Service.Icon,
			&h.Service.CreatedAt,
			&h.Service.UpdatedAt,
			&h.HostName,
			&h.LastMessage,
		)
		if err != nil {
			log.Println(err)
			return services, err
		}
		services = append(services, h)
	}

	return services, nil
}

// GetHostServiceByHostIDServiceID gets a host service by host id and service id
func (repo *PostgresRepository) GetHostServiceByHostIDServiceID(hostID, serviceID int) (HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number,
			hs.schedule_unit, hs.last_check, hs.status, hs.created_at, hs.updated_at,
			s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at, h.host_name,
		    hs.last_message

		from host_services hs
		left join services s on (hs.service_id = s.id)
		left join hosts h on (hs.host_id = h.id)

		where hs.host_id = $1 and hs.service_id = $2
`

	var hs HostService

	row := repo.DB.QueryRowContext(ctx, query, hostID, serviceID)

	err := row.Scan(
		&hs.ID,
		&hs.HostID,
		&hs.ServiceID,
		&hs.Active,
		&hs.ScheduleNumber,
		&hs.ScheduleUnit,
		&hs.LastCheck,
		&hs.Status,
		&hs.CreatedAt,
		&hs.UpdatedAt,
		&hs.Service.ID,
		&hs.Service.ServiceName,
		&hs.Service.Active,
		&hs.Service.Icon,
		&hs.Service.CreatedAt,
		&hs.Service.UpdatedAt,
		&hs.HostName,
		&hs.LastMessage,
	)

	if err != nil {
		log.Println(err)
		return hs, err
	}

	return hs, nil
}

// InsertEvent inserts an event into the database
func (repo *PostgresRepository) InsertEvent(e Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into events (host_service_id, event_type, host_id, service_name, host_name,
			message, created_at, updated_at)
		values
			($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := repo.DB.ExecContext(ctx, stmt,
		e.HostServiceID,
		e.EventType,
		e.HostID,
		e.ServiceName,
		e.HostName,
		e.Message,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// GetAllEvents gets all events
func (repo *PostgresRepository) GetAllEvents() ([]Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, event_type, host_service_id, host_id, service_name, host_name,
			message, created_at, updated_at from events order by created_at`

	var events []Event

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var ev Event
		err := rows.Scan(
			&ev.ID,
			&ev.EventType,
			&ev.HostServiceID,
			&ev.HostID,
			&ev.ServiceName,
			&ev.HostName,
			&ev.Message,
			&ev.CreatedAt,
			&ev.UpdatedAt,
		)
		if err != nil {
			return events, err
		}
		events = append(events, ev)
	}

	return events, nil
}
