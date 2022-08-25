package dbrepo

import (
	"context"
	"github.com/virmos/cycir/internal/models"
	"log"
	"time"
)

func (m *postgresDBRepo) AllDepartments() ([]*models.Department, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := 
		`WITH RECURSIVE department_tree_view AS (
				SELECT id,
						COALESCE(parent_id,0) AS parent_id,
						name,
						1 AS generation_number,
						CAST(id AS varchar(50)) AS order_sequence
				FROM departments
				WHERE parent_id IS NULL
				
		UNION ALL
		
				SELECT parent.id,
						parent.parent_id,
						parent.name,
						generation_number + 1 AS generation_number,
						CAST(order_sequence || '_' || CAST(parent.id AS VARCHAR (50)) AS VARCHAR(50)) AS order_sequence
				FROM departments parent
				JOIN department_tree_view tv
					ON parent.parent_id = tv.id
		) SELECT
				id, name, parent_id, order_sequence, generation_number
				AS departments_tree
		FROM department_tree_view
		ORDER BY order_sequence;`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*models.Department

	for rows.Next() {
		s := &models.Department{}
		err = rows.Scan(&s.ID, &s.Name, &s.ParentId, &s.OrderSequence, &s.TreeLevel)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Append it to the slice
		departments = append(departments, s)
	}
	
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return departments, nil
}
