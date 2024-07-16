package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/internal/model"
)

type EquipmentRepository struct {
	db *pgxpool.Pool
}

func NewEquipmentRepository(db *pgxpool.Pool) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) Create(ctx context.Context, date int64, company int64, serialNumber string, profile int64, userId int64) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return
			}
		} else {
			if err := tx.Commit(ctx); err != nil {
				return
			}
		}
	}()

	queryCreateEquipment := `
			INSERT INTO equipments (serial_number, profile) 
			VALUES ($1, $2)
			RETURNING id;`

	equipment := new(model.Equipment)
	if err = tx.QueryRow(ctx, queryCreateEquipment, serialNumber, profile).Scan(&equipment); err != nil {
		return 0, err
	}

	tm := time.Unix(date, 0)
	queryLocationRecord := `	
			INSERT INTO locations (date, code, equipment, employee, company) 
			VALUES ($1, $2, $3, $4, $5);`

	if _, err = tx.Exec(ctx, queryLocationRecord, tm, "ADD_TO_STORAGE", equipment.ID, userId, company); err != nil {
		return 0, err
	}

	return equipment.ID, nil
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

	if err := r.db.QueryRow(ctx, query, id).Scan(
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

func (r *EquipmentRepository) GetByProfile(ctx context.Context, id int64) ([]*model.Equipment, error) {
	var equipments []*model.Equipment
	equipment := new(model.Equipment)

	query := `
			SELECT equipments.id, equipments.serial_number
			FROM equipments
			LEFT JOIN profiles ON profiles.id = equipments.profile
			WHERE profiles.id = $1;`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&equipment.ID, &equipment.SerialNumber); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	return equipments, nil
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

	rows, err := r.db.Query(ctx, query)
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

	rows, err := r.db.Query(ctx, query, toDepartment)
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

	rows, err := r.db.Query(ctx, query, toEmployee)
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

	rows, err := r.db.Query(ctx, query, toContract)
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

	rows, err := r.db.Query(ctx, query, toDepartment, toEmployee)
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

func (r *EquipmentRepository) GetAll(ctx context.Context) ([]*model.Equipment, error) {
	var equipments []*model.Equipment
	equipment := new(model.Equipment)

	query := `
			SELECT equipments.id, equipments.serial_number, 
			       profiles.id, profiles.title, 
			       categories.id, categories.title
			FROM equipments
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			WHERE equipments.is_deleted = FALSE
			ORDER BY profiles.title;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&equipment.ID,
			&equipment.SerialNumber,
			&equipment.Profile.ID,
			&equipment.Profile.Title,
			&equipment.Profile.Category.ID,
			&equipment.Profile.Category.Title); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

func (r *EquipmentRepository) FindBySerialNumber(ctx context.Context, serialNumber string) (int64, error) {
	equipment := new(model.Equipment)

	query := `
			SELECT id 
			FROM equipments 
			WHERE serial_number = $1;`

	if err := r.db.QueryRow(ctx, query, serialNumber).Scan(&equipment.ID); err != nil {
		return 0, err
	}

	return equipment.ID, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, id int64, serialNumber string, profile int64) error {
	query := `
			UPDATE equipments 
			SET serial_number = $2, profile = $3
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, serialNumber, profile); err != nil {
		return err
	}

	return nil
}

func (r *EquipmentRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE equipments 
			SET is_deleted = true
       		WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
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

	rows, err := r.db.Query(ctx, query, categoryId, departmentId, date)
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

	rows, err := r.db.Query(ctx, query, categoryId, departmentId, code, fromDate, toDate)
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

	rows, err := r.db.Query(ctx, query, categoryId, departmentId, fromDate, toDate)
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

	rows, err := r.db.Query(ctx, query, categoryId, departmentId, fromDate, toDate)
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
