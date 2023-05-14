package repository

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
	"time"
	"warehouse_backend/pkg/model"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) TransferToStorage(date int64, code string, equipment, employee, company int) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, transfer_type, price) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, "", "").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToDepartment(date int64, code string, equipment, employee, company, toDepartment int) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_department, transfer_type, price) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toDepartment, "", "").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToEmployee(date int64, code string, equipment, employee, company, toEmployee int) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_employee, transfer_type, price) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toEmployee, "", "").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToEmployeeInDepartment(date int64, code string, equipment, employee, company, toDepartment, toEmployee int) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_department, to_employee, transfer_type, price) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toDepartment, toEmployee, "", "").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToContract(date int64, code string, equipment, employee, company, toContract int, transferType, price string) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_contract, transfer_type, price) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toContract, transferType, price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) GetHistory(id int) ([]model.Location, error) {
	var histories []model.Location
	var history model.Location
	var toD, toE, toCNumber, toCAddress interface{}
	var date time.Time
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
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&history.Id,
			&date,
			&history.Code,
			&history.TransferType,
			&history.Price,
			&history.Employee.Name,
			&history.Company.Title,
			&toD,
			&toE,
			&toCNumber,
			&toCAddress)
		history.Date = date.Unix()
		history.ToDepartment.Title = InterfaceToString(toD)
		history.ToEmployee.Name = InterfaceToString(toE)
		history.ToContract.Number = InterfaceToString(toCNumber)
		history.ToContract.Address = InterfaceToString(toCAddress)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, err
}

func (r *LocationRepository) Delete(id int) error {
	query := `
			DELETE FROM locations 
       		WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
