package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type LocationRepository struct {
	DB *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{DB: db}
}

// AddToStorage is equipment add to storage
func (r *LocationRepository) AddToStorage(ctx context.Context, date string, equipmentId, employeeId, companyId int64) error {
	const query = `
		INSERT INTO locations (date, code, equipment, employee, company)
		VALUES ($1, $2, $3, $4, $5);`

	if _, err := r.DB.Exec(ctx, query, date, "ADD_TO_STORAGE", equipmentId, employeeId, companyId); err != nil {
		return err
	}

	return nil
}

// TransferToStorage is equipment transfer to storage
func (r *LocationRepository) TransferToStorage(ctx context.Context, date string, code string, equipmentId, employeeId, companyId int64, nowLocation []interface{}) (int64, error) {
	var id int64
	const query = `
		INSERT INTO locations (date, code, equipment, employee, company, from_department, from_employee, from_contract) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;`

	if err := r.DB.QueryRow(ctx, query, date, code, equipmentId, employeeId, companyId, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// TransferToDepartment is equipment transfer to department
func (r *LocationRepository) TransferToDepartment(ctx context.Context, date string, code string, equipmentId, employeeId, companyId, toDepartment int64, nowLocation []interface{}) (int64, error) {
	var id int64
	const query = `
		INSERT INTO locations (date, code, equipment, employee, company, to_department, from_department, from_employee, from_contract) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;`

	if err := r.DB.QueryRow(ctx, query, date, code, equipmentId, employeeId, companyId, toDepartment, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// TransferToEmployee is equipment transfer to employee
func (r *LocationRepository) TransferToEmployee(ctx context.Context, date string, code string, equipmentId, employeeId, companyId, toEmployee int64, nowLocation []interface{}) (int64, error) {
	var id int64
	const query = `
		INSERT INTO locations (date, code, equipment, employee, company, to_employee, from_department, from_employee, from_contract) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;`

	if err := r.DB.QueryRow(ctx, query, date, code, equipmentId, employeeId, companyId, toEmployee, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// TransferToEmployeeInDepartment is equipment transfer to employee in department
func (r *LocationRepository) TransferToEmployeeInDepartment(ctx context.Context, date string, code string, equipmentId, employeeId, companyId, toDepartment, toEmployee int64, nowLocation []interface{}) (int64, error) {
	var id int64
	const query = `
		INSERT INTO locations (date, code, equipment, employee, company, to_department, to_employee, from_department, from_employee, from_contract) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;`

	if err := r.DB.QueryRow(ctx, query, date, code, equipmentId, employeeId, companyId, toDepartment, toEmployee, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// TransferToContract is equipment transfer to contract
func (r *LocationRepository) TransferToContract(ctx context.Context, date string, code string, equipmentId, employeeId, companyId, toContract int64, transferType string, price int, nowLocation []interface{}) (int64, error) {
	var id int64
	const query = `
		INSERT INTO locations (date, code, equipment, employee, company, to_contract, transfer_type, price, from_department, from_employee, from_contract) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id;`

	if err := r.DB.QueryRow(ctx, query, date, code, equipmentId, employeeId, companyId, toContract, transferType, price, nowLocation[0], nowLocation[1], nowLocation[2]).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// Delete is equipment transfer delete
func (r *LocationRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM locations 
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetById is equipment get by id
func (r *LocationRepository) GetById(ctx context.Context, equipmentId int64) (*model.Location, error) {
	equipmentByLoc := new(model.Location)

	const query = `
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

	if err := r.DB.QueryRow(ctx, query, equipmentId).Scan(
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

// GetHistory is equipment get history
func (r *LocationRepository) GetHistory(ctx context.Context, equipmentId int64) ([]*model.Location, error) {
	var histories []*model.Location

	const query = `
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

	rows, err := r.DB.Query(ctx, query, equipmentId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		history := new(model.Location)
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
			&history.ToContract.Address,
		); err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, err
}

// GetLocationNow is equipment get location now
func (r *LocationRepository) GetLocationNow(ctx context.Context, equipmentId int64) ([]interface{}, error) {
	var department, employee, contract interface{}
	var arr []interface{}

	const query = `
		SELECT to_department, to_employee, to_contract 
		FROM locations 
       	WHERE locations.id IN 
       	      (SELECT MAX(locations.id)
       	       FROM locations
       	       GROUP BY locations.equipment)
		AND locations.equipment = $1;`

	if err := r.DB.QueryRow(ctx, query, equipmentId).Scan(
		&department,
		&employee,
		&contract,
	); err != nil {
		return []interface{}{}, err
	}
	arr = append(arr, department, employee, contract)

	return arr, nil
}

// GetByLocationStorage is equipment get by location storage
func (r *LocationRepository) GetByLocationStorage(ctx context.Context) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location

	const query = `
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
		AND equipments.deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		equipmentByLoc := newLocation()
		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
		); err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

// GetByLocationDepartment is equipment get by location department
func (r *LocationRepository) GetByLocationDepartment(ctx context.Context, toDepartment int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location

	const query = `
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
		AND equipments.deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toDepartment)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		equipmentByLoc := newLocation()
		var (
			toDepartmentId, toEmployeeId      sql.NullInt64
			toDepartmentTitle, toEmployeeName sql.NullString
		)

		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
			&toDepartmentId,
			&toDepartmentTitle,
			&toEmployeeId,
			&toEmployeeName,
		); err != nil {
			return nil, err
		}

		equipmentByLoc.ToDepartment.ID = validInt64(toDepartmentId)
		equipmentByLoc.ToDepartment.Title = validString(toDepartmentTitle)
		equipmentByLoc.ToEmployee.ID = validInt64(toEmployeeId)
		equipmentByLoc.ToEmployee.Name = validString(toEmployeeName)

		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

// GetByLocationEmployee is equipment get by location employee
func (r *LocationRepository) GetByLocationEmployee(ctx context.Context, toEmployee int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location

	const query = `
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
		AND equipments.deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toEmployee)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		equipmentByLoc := newLocation()
		var toDepartmentTitle, toEmployeeName sql.NullString

		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
			&toDepartmentTitle,
			&toEmployeeName,
		); err != nil {
			return nil, err
		}

		equipmentByLoc.ToDepartment.Title = validString(toDepartmentTitle)
		equipmentByLoc.ToEmployee.Name = validString(toEmployeeName)

		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

// GetByLocationContract is equipment get by location contract
func (r *LocationRepository) GetByLocationContract(ctx context.Context, toContract int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location

	const query = `
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
		AND equipments.deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toContract)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		equipmentByLoc := newLocation()

		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
		); err != nil {
			return nil, err
		}

		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

// GetByLocationDepartmentEmployee is equipment get by location department employee
func (r *LocationRepository) GetByLocationDepartmentEmployee(ctx context.Context, toDepartment, toEmployee int64) ([]*model.Location, error) {
	var equipmentsByLoc []*model.Location

	const query = `
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
		AND equipments.deleted = FALSE;`

	rows, err := r.DB.Query(ctx, query, toDepartment, toEmployee)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		equipmentByLoc := newLocation()
		var toDepartmentTitle, toEmployeeName sql.NullString

		if err := rows.Scan(
			&equipmentByLoc.Equipment.ID,
			&equipmentByLoc.Equipment.SerialNumber,
			&equipmentByLoc.Equipment.Profile.Title,
			&equipmentByLoc.Equipment.Profile.Category.Title,
			&equipmentByLoc.Company.ID,
			&equipmentByLoc.Company.Title,
			&toDepartmentTitle,
			&toEmployeeName,
		); err != nil {
			return nil, err
		}

		equipmentByLoc.ToDepartment.Title = validString(toDepartmentTitle)
		equipmentByLoc.ToEmployee.Name = validString(toEmployeeName)

		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}

	return equipmentsByLoc, nil
}

// RemainderByCategory is remainder equipment get by category
func (r *LocationRepository) RemainderByCategory(ctx context.Context, categoryId, departmentId int64, date string) ([]*model.Location, error) {
	var locations []*model.Location

	const query = `
		SELECT locations.date, 
		       profiles.title, 
		       equipments.serial_number
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
		location := newLocation()

		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
		); err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}

// TransferByCategory is transfer equipment get by category
func (r *LocationRepository) TransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate string, code string) ([]*model.Location, error) {
	var locations []*model.Location
	var query string

	switch {
	case code == "STORAGE_TO_DEPARTMENT" || code == "CONTRACT_TO_DEPARTMENT":
		query = `
			SELECT locations.date, 
			       profiles.title, 
			       equipments.serial_number
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
			SELECT locations.date, 
			       profiles.title, 
			       equipments.serial_number
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
		location := newLocation()

		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
		); err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}

// ToDepartmentTransferByCategory is transfer equipment to department by category
func (r *LocationRepository) ToDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate string) ([]*model.Location, error) {
	var locations []*model.Location

	const query = `
		SELECT locations.date, 
		       profiles.title, 
		       equipments.serial_number, 
		       from_department.id, from_department.title
		FROM locations
		LEFT JOIN equipments ON locations.equipment = equipments.id
		LEFT JOIN profiles ON equipments.profile = profiles.id
		LEFT JOIN departments from_department ON from_department.id = locations.from_department
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
		location := newLocation()
		var (
			fromDepartmentId    sql.NullInt64
			fromDepartmentTitle sql.NullString
		)

		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
			&fromDepartmentId,
			&fromDepartmentTitle,
		); err != nil {
			return nil, err
		}

		location.ToDepartment.ID = validInt64(fromDepartmentId)
		location.ToDepartment.Title = validString(fromDepartmentTitle)

		locations = append(locations, location)
	}

	return locations, nil
}

// FromDepartmentTransferByCategory is transfer equipment from department by category
func (r *LocationRepository) FromDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate string) ([]*model.Location, error) {
	var locations []*model.Location

	const query = `
		SELECT locations.date, 
		       profiles.title, 
		       equipments.serial_number, 
		       to_department.id, to_department.title
		FROM locations
		LEFT JOIN equipments ON locations.equipment = equipments.id
		LEFT JOIN profiles ON equipments.profile = profiles.id
		LEFT JOIN departments to_department ON to_department.id = locations.to_department
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
		location := newLocation()
		var (
			toDepartmentId    sql.NullInt64
			toDepartmentTitle sql.NullString
		)

		if err := rows.Scan(
			&location.Date,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
			&toDepartmentId,
			&toDepartmentTitle,
		); err != nil {
			return nil, err
		}

		location.ToDepartment.ID = validInt64(toDepartmentId)
		location.ToDepartment.Title = validString(toDepartmentTitle)

		locations = append(locations, location)
	}

	return locations, nil
}

func newLocation() *model.Location {
	return &model.Location{
		Equipment: &model.Equipment{
			Profile: &model.Profile{
				Category: &model.Category{},
			},
		},
		Employee:       &model.Employee{},
		Company:        &model.Company{},
		FromDepartment: &model.Department{},
		FromEmployee:   &model.Employee{},
		FromContract:   &model.Contract{},
		ToDepartment:   &model.Department{},
		ToEmployee:     &model.Employee{},
		ToContract:     &model.Contract{},
	}
}

func validInt64(num sql.NullInt64) int64 {
	if num.Valid {
		return num.Int64
	}
	return 0
}

func validString(num sql.NullString) string {
	if num.Valid {
		return num.String
	}
	return ""
}
