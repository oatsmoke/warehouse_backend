package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/pkg/model"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) TransferToStorage(date int64, code string, equipment, employee, company int, nowLocation []interface{}) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToDepartment(date int64, code string, equipment, employee, company, toDepartment int, nowLocation []interface{}) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_department, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toDepartment, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToEmployee(date int64, code string, equipment, employee, company, toEmployee int, nowLocation []interface{}) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_employee, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toEmployee, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToEmployeeInDepartment(date int64, code string, equipment, employee, company, toDepartment, toEmployee int, nowLocation []interface{}) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_department, to_employee, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toDepartment, toEmployee, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) TransferToContract(date int64, code string, equipment, employee, company, toContract int, transferType string, price int, nowLocation []interface{}) (int, error) {
	query := `
			INSERT INTO locations (date, code, equipment, employee, company, to_contract, transfer_type, price, from_department, from_employee, from_contract) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id;`
	tm := time.Unix(date, 0)
	id := 0
	err := r.db.QueryRow(context.Background(), query, tm, code, equipment, employee, company, toContract, transferType, price, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LocationRepository) GetHistory(id int) ([]model.Location, error) {
	var histories []model.Location
	var history model.Location
	var toD, toE, toCNumber, toCAddress, transferType, price interface{}
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
		err = rows.Scan(&history.Id, &date, &history.Code, &transferType, &price, &history.Employee.Name, &history.Company.Title, &toD, &toE, &toCNumber, &toCAddress)
		history.Date = date.Unix()
		history.ToDepartment.Title = InterfaceToString(toD)
		history.ToEmployee.Name = InterfaceToString(toE)
		history.ToContract.Number = InterfaceToString(toCNumber)
		history.ToContract.Address = InterfaceToString(toCAddress)
		history.TransferType = InterfaceToString(transferType)
		history.Price = InterfaceToInt(price)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, err
}

func (r *LocationRepository) GetLocationNow(id int) ([]interface{}, error) {
	query := `
			SELECT to_department, to_employee, to_contract 
			    FROM locations 
       		WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND locations.equipment = $1;`
	var department, employee, contract interface{}
	var arr []interface{}
	if err := r.db.QueryRow(context.Background(), query, id).Scan(&department, &employee, &contract); err != nil {
		return []interface{}{}, err
	}
	arr = append(arr, department, employee, contract)
	return arr, nil
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
