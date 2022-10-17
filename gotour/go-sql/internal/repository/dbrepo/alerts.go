package dbrepo

import (
	"context"
	"github.com/virmos/cycir/internal/models"
	"log"
	"time"
)

func (m *postgresDBRepo) AllAlerts() ([]*models.Alert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := 
		`with alerts_groupby_type_intervals as (select id, FORMAT('%s/%s', date_part('year', _interval::timestamp), date_part('month', _interval::timestamp)) as date_interval, "type" from alerts
		order by "_interval" )
		select date_interval, type, COUNT(type) from alerts_groupby_type_intervals
		group by date_interval, type
		order by date_interval`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert

	for rows.Next() {
		al := &models.Alert{}
		err = rows.Scan(&al.Interval, &al.Type, &al.Count)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Append it to the slice
		alerts = append(alerts, al)
	}
	
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return alerts, nil
}
