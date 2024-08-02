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
func (r *EquipmentRepository) GetById(ctx context.Context, id int64) (*model.Location, error) {
	equipmentByLoc := new(model.Location)

	query := `
			SELECT equipments.id, equipments.serial_number, 
			       profiles.id, profiles.title, 
			       categories.id, categories.title,
			       companies.id, companies.title,
			       to_department.id,
			       to_employee.id,
			       to_contract.id,
				   locations.transfer_type, locations.price
			FROM locations
			LEFT JOIN equipments ON equipments.id = locations.equipment    
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			LEFT JOIN companies ON companies.id = locations.company
			LEFT JOIN departments to_department ON to_department.id = locations.to_department
			LEFT JOIN employees to_employee ON to_employee.id = locations.to_employee
			LEFT JOIN contracts to_contract ON to_contract.id = locations.to_contract
			WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND equipments.id = $1;`

	if err := r.DB.QueryRow(ctx, query, id).Scan(
		&equipmentByLoc.Equipment.ID,
		&equipmentByLoc.Equipment.SerialNumber,
		&equipmentByLoc.Equipment.Profile.ID,
		&equipmentByLoc.Equipment.Profile.Title,
		&equipmentByLoc.Equipment.Profile.Category.ID,
		&equipmentByLoc.Equipment.Profile.Category.Title,
		&equipmentByLoc.Company.ID,
		&equipmentByLoc.Company.Title,
		&equipmentByLoc.ToDepartment.ID,
		&equipmentByLoc.ToEmployee.ID,
		&equipmentByLoc.ToContract.ID,
		&equipmentByLoc.TransferType,
		&equipmentByLoc.Price); err != nil {
		return nil, err
	}

	return equipmentByLoc, nil
}

func (r *EquipmentRepository) GetByLocationStorage(ctx context.Context) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location
	equipmentByLoc := new(model.Location)

	query := `
			SELECT equipments.id, equipments.serial_number, 
			       profiles.title, 
			       categories.title,
			       companies.id, companies.title
			FROM locations
			LEFT JOIN equipments ON equipments.id = locations.equipment
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			LEFT JOIN companies ON companies.id = locations.company
			WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND locations.to_department IS NULL
			AND locations.to_employee IS NULL
			AND locations.to_contract IS NULL
			AND equipments.is_deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title); err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationDepartment(ctx context.Context, toDepartment int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location
	equipmentByLoc := new(model.Location)

	query := `
			SELECT equipments.id, equipments.serial_number, 
			       profiles.title, 
			       categories.title,
			       companies.id, companies.title,
			       to_department.id, to_department.title,
			       to_employee.id, to_employee.name
			FROM locations
			LEFT JOIN equipments ON equipments.id = locations.equipment
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			LEFT JOIN companies ON companies.id = locations.company
			LEFT JOIN departments to_department ON to_department.id = locations.to_department
			LEFT JOIN employees to_employee ON to_employee.id = locations.to_employee
			WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND locations.to_department = $1
			AND equipments.is_deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toDepartment)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
			&equipmentByLoc.ToDepartment.ID,
			&equipmentByLoc.ToDepartment.Title,
			&equipmentByLoc.ToEmployee.ID,
			&equipmentByLoc.ToEmployee.Name); err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationEmployee(ctx context.Context, toEmployee int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location
	equipmentByLoc := new(model.Location)

	query := `
			SELECT equipments.id, equipments.serial_number, 
			       profiles.title, 
			       categories.title,
			       companies.id, companies.title,
			       to_department.title,
			       to_employee.name
			FROM locations
			LEFT JOIN equipments ON equipments.id = locations.equipment
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			LEFT JOIN companies ON companies.id = locations.company
			LEFT JOIN departments to_department ON to_department.id = locations.to_department
			LEFT JOIN employees to_employee ON to_employee.id = locations.to_employee
			WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND locations.to_employee = $1
			AND equipments.is_deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toEmployee)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
			&equipmentByLoc.ToDepartment.Title,
			&equipmentByLoc.ToEmployee.Name); err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationContract(ctx context.Context, toContract int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location
	equipmentByLoc := new(model.Location)

	query := `
			SELECT equipments.id, equipments.serial_number, 
			       profiles.title, 
			       categories.title,
			       companies.id, companies.title
			FROM locations
			LEFT JOIN equipments ON equipments.id = locations.equipment
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			LEFT JOIN companies ON companies.id = locations.company
			LEFT JOIN departments ON departments.id = locations.to_department
			LEFT JOIN employees ON employees.id = locations.to_employee
			LEFT JOIN contracts ON contracts.id = locations.to_contract
			WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND locations.to_department IS NULL
			AND locations.to_employee IS NULL
			AND locations.to_contract = $1
			AND equipments.is_deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toContract)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title); err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationDepartmentEmployee(ctx context.Context, toDepartment, toEmployee int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location
	equipmentByLoc := new(model.Location)

	query := `
			SELECT equipments.id, equipments.serial_number, 
			       profiles.title, 
			       categories.title,
			       companies.id, companies.title,
			       to_department.title,
			       to_employee.name
			FROM locations
			LEFT JOIN equipments ON equipments.id = locations.equipment
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			LEFT JOIN companies ON companies.id = locations.company
			LEFT JOIN departments to_department ON to_department.id = locations.to_department
			LEFT JOIN employees to_employee ON to_employee.id = locations.to_employee
			WHERE locations.id IN 
			(SELECT MAX(locations.id)
			 FROM locations
			 GROUP BY locations.equipment)
			AND locations.to_department = $1
			AND locations.to_employee = $2
			AND equipments.is_deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toDepartment, toEmployee)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
			&equipmentByLoc.ToDepartment.Title,
			&equipmentByLoc.ToEmployee.Name); err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) RemainderByCategory(ctx context.Context, categoryId, departmentId int64, date time.Time) ([]*model.Location, error) {
	var locations []*model.Location
	location := new(model.Location)

	query := `
			SELECT locations.date, profiles.title, equipments.serial_number
			FROM locations
			         LEFT JOIN equipments ON locations.equipment = equipments.id
			         LEFT JOIN profiles ON equipments.profile = profiles.id
			         LEFT JOIN departments ON departments.id = locations.to_department
			WHERE locations.id IN
			      (SELECT MAX(locations.id)
				   FROM locations
				   WHERE locations.date < $3
				   GROUP BY locations.equipment)
			  AND profiles.category = $1
			  AND locations.to_department = $2
			  AND locations.date < $3;`

	rows, err := r.DB.Query(ctx, query, categoryId, departmentId, date)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (r *EquipmentRepository) TransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate time.Time, code string) ([]*model.Location, error) {
	var locations []*model.Location
	location := new(model.Location)
	var query string

	switch {
	case code == "STORAGE_TO_DEPARTMENT" || code == "CONTRACT_TO_DEPARTMENT":
		query = `
			SELECT locations.date, profiles.title, equipments.serial_number
			FROM locations
			         LEFT JOIN equipments ON locations.equipment = equipments.id
			         LEFT JOIN profiles ON equipments.profile = profiles.id
			WHERE profiles.category = $1
			  AND locations.to_department = $2
			  AND locations.code = $3
			  AND locations.date >= $4
			  AND locations.date < $5;`
	case code == "DEPARTMENT_TO_STORAGE" || code == "DEPARTMENT_TO_CONTRACT":
		query = `
			SELECT locations.date, profiles.title, equipments.serial_number
			FROM locations
			         LEFT JOIN equipments ON locations.equipment = equipments.id
			         LEFT JOIN profiles ON equipments.profile = profiles.id
			WHERE profiles.category = $1
			  AND locations.from_department = $2
			  AND locations.code = $3
			  AND locations.date >= $4
			  AND locations.date < $5;`
	}

	rows, err := r.DB.Query(ctx, query, categoryId, departmentId, code, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (r *EquipmentRepository) ToDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate time.Time) ([]*model.Location, error) {
	var locations []*model.Location
	location := new(model.Location)

	query := `
			SELECT locations.date, profiles.title, equipments.serial_number, departments.id, departments.title
			FROM locations
			         LEFT JOIN equipments ON locations.equipment = equipments.id
			         LEFT JOIN profiles ON equipments.profile = profiles.id
			         LEFT JOIN departments on departments.id = locations.from_department
			WHERE profiles.category = $1
			  AND locations.to_department = $2
			  AND locations.code = 'DEPARTMENT_TO_DEPARTMENT'
			  AND locations.date >= $3
			  AND locations.date < $4;`

	rows, err := r.DB.Query(ctx, query, categoryId, departmentId, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
			&location.ToDepartment.ID,
			&location.ToDepartment.Title); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (r *EquipmentRepository) FromDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate time.Time) ([]*model.Location, error) {
	var locations []*model.Location
	location := new(model.Location)

	query := `
			SELECT locations.date, profiles.title, equipments.serial_number, departments.id, departments.title
			FROM locations
			         LEFT JOIN equipments ON locations.equipment = equipments.id
			         LEFT JOIN profiles ON equipments.profile = profiles.id
			         LEFT JOIN departments on departments.id = locations.to_department
			WHERE profiles.category = $1
			  AND locations.from_department = $2
			  AND locations.code = 'DEPARTMENT_TO_DEPARTMENT'
			  AND locations.date >= $3
			  AND locations.date < $4;`

	rows, err := r.DB.Query(ctx, query, categoryId, departmentId, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
			&location.FromDepartment.ID,
			&location.FromDepartment.Title); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}
