package dbrepo

import (
	"context"
	"github.com/virmos/cycir/internal/models"
	"log"
	"time"
)

// GetUserById returns a user by id
func (m *postgresDBRepo) GetEmployeesByDepartment(id int) ([]*models.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, first_name, last_name,  username, email, 
			created_at, updated_at
			FROM employees where department_id = $1`
	rows, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*models.Employee

	for rows.Next() {
		e := &models.Employee{}
		err := rows.Scan(
			&e.ID,
			&e.FirstName,
			&e.LastName,
			&e.UserName,
			&e.Email,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Append it to the slice
		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return employees, nil
}
