package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/internal/model"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) TransferToStorage(ctx context.Context, date int64, code string, equipment, employee, company int64, nowLocation []interface{}) (int64, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id;`

	tm := time.Unix(date, 0)
	location := new(model.Location)

	if err := r.db.QueryRow(ctx, query, tm, code, equipment, employee, company, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&location.ID); err != nil {
		return 0, err
	}

	return location.ID, nil
}

func (r *LocationRepository) TransferToDepartment(ctx context.Context, date int64, code string, equipment, employee, company, toDepartment int64, nowLocation []interface{}) (int64, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_department, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id;`

	tm := time.Unix(date, 0)
	location := new(model.Location)

	if err := r.db.QueryRow(ctx, query, tm, code, equipment, employee, company, toDepartment, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&location.ID); err != nil {
		return 0, err
	}

	return location.ID, nil
}

func (r *LocationRepository) TransferToEmployee(ctx context.Context, date int64, code string, equipment, employee, company, toEmployee int64, nowLocation []interface{}) (int64, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_employee, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id;`

	tm := time.Unix(date, 0)
	location := new(model.Location)

	if err := r.db.QueryRow(ctx, query, tm, code, equipment, employee, company, toEmployee, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&location.ID); err != nil {
		return 0, err
	}

	return location.ID, nil
}

func (r *LocationRepository) TransferToEmployeeInDepartment(ctx context.Context, date int64, code string, equipment, employee, company, toDepartment, toEmployee int64, nowLocation []interface{}) (int64, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_department, to_employee, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id;`

	tm := time.Unix(date, 0)
	location := new(model.Location)
	if err := r.db.QueryRow(ctx, query, tm, code, equipment, employee, company, toDepartment, toEmployee, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&location.ID); err != nil {
		return 0, err
	}

	return location.ID, nil
}

func (r *LocationRepository) TransferToContract(ctx context.Context, date int64, code string, equipment, employee, company, toContract int64, transferType string, price int, nowLocation []interface{}) (int64, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_contract, transfer_type, price, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id;`

	tm := time.Unix(date, 0)
	location := new(model.Location)

	if err := r.db.QueryRow(ctx, query, tm, code, equipment, employee, company, toContract, transferType, price, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&location.ID); err != nil {
		return 0, err
	}
	return location.ID, nil
}

func (r *LocationRepository) GetHistory(ctx context.Context, id int64) ([]*model.Location, error) {
	var histories []*model.Location
	history := new(model.Location)

	query := `
			SELECT locations.id, locations.date, locations.code, locations.transfer_type, locations.price,
				employees.name,
				companies.title,
				to_department.title,
				to_employee.name,
				to_contract.number, to_contract.address
			FROM locations
			LEFT JOIN employees ON employees.id = locations.employee
			LEFT JOIN companies ON companies.id = locations.company
			LEFT JOIN departments to_department ON to_department.id = locations.to_department
			LEFT JOIN employees to_employee ON to_employee.id = locations.to_employee
			LEFT JOIN contracts to_contract ON to_contract.id = locations.to_contract
			WHERE locations.equipment = $1
			ORDER BY locations.id DESC;`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&history.ID,
			&history.Date,
			&history.Code,
			&history.TransferType,
			&history.Price,
			&history.Employee.Name,
			&history.Company.Title,
			&history.ToDepartment.Title,
			&history.ToEmployee.Name,
			&history.ToContract.Number,
			&history.ToContract.Address); err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, err
}

func (r *LocationRepository) GetLocationNow(ctx context.Context, id int64) ([]interface{}, error) {
	var department, employee, contract interface{}
	var arr []interface{}

	query := `
			SELECT to_department, to_employee, to_contract 
			    FROM locations 
       		WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND locations.equipment = $1;`

	if err := r.db.QueryRow(ctx, query, id).Scan(&department, &employee, &contract); err != nil {
		return []interface{}{}, err
	}
	arr = append(arr, department, employee, contract)

	return arr, nil
}

func (r *LocationRepository) Delete(ctx context.Context, id int64) error {
	query := `
			DELETE FROM locations 
       		WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}
