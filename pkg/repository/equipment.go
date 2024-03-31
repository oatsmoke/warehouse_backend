package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/pkg/model"
)

type EquipmentRepository struct {
	db *pgxpool.Pool
}

func NewEquipmentRepository(db *pgxpool.Pool) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) Create(date int64, company int, serialNumber string, profile int, userId int) (int, error) {
	ctx := context.Background()
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
	id := 0
	err = tx.QueryRow(ctx, queryCreateEquipment, serialNumber, profile).Scan(&id)
	if err != nil {
		return 0, err
	}
	tm := time.Unix(date, 0)
	queryLocationRecord := `	
			INSERT INTO locations (date, code, equipment, employee, company) 
			VALUES ($1, $2, $3, $4, $5);`
	_, err = tx.Exec(ctx, queryLocationRecord, tm, "ADD_TO_STORAGE", id, userId, company)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *EquipmentRepository) GetById(id int) (model.Location, error) {
	var equipmentByLoc model.Location
	var transferType, price, toD, toE, toC interface{}
	query := `
			SELECT locations.transfer_type, locations.price,
			       equipments.id, equipments.serial_number, 
			       profiles.id, profiles.title, 
			       categories.id, categories.title,
			       companies.id, companies.title,
			       to_department.id,
			       to_employee.id,
			       to_contract.id
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
	err := r.db.QueryRow(context.Background(), query, id).Scan(&transferType, &price, &equipmentByLoc.Equipment.Id, &equipmentByLoc.Equipment.SerialNumber, &equipmentByLoc.Equipment.Profile.Id, &equipmentByLoc.Equipment.Profile.Title, &equipmentByLoc.Equipment.Profile.Category.Id, &equipmentByLoc.Equipment.Profile.Category.Title, &equipmentByLoc.Company.Id, &equipmentByLoc.Company.Title, &toD, &toE, &toC)
	equipmentByLoc.TransferType = InterfaceToString(transferType)
	equipmentByLoc.Price = InterfaceToInt(price)
	equipmentByLoc.ToDepartment.Id = InterfaceToInt(toD)
	equipmentByLoc.ToEmployee.Id = InterfaceToInt(toE)
	equipmentByLoc.ToContract.Id = InterfaceToInt(toC)
	if err != nil {
		return model.Location{}, err
	}
	return equipmentByLoc, nil
}

func (r *EquipmentRepository) GetByProfile(id int) ([]model.Equipment, error) {
	var equipments []model.Equipment
	var equipment model.Equipment
	query := `
			SELECT equipments.id, equipments.serial_number
			FROM equipments
			LEFT JOIN profiles ON profiles.id = equipments.profile
			WHERE profiles.id = $1;`
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&equipment.Id, &equipment.SerialNumber)
		if err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}
	return equipments, nil
}

func (r *EquipmentRepository) GetByLocationStorage() ([]model.Location, error) {
	var equipmentsByLoc []model.Location
	var equipmentByLoc model.Location
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
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&equipmentByLoc.Equipment.Id, &equipmentByLoc.Equipment.SerialNumber, &equipmentByLoc.Equipment.Profile.Title, &equipmentByLoc.Equipment.Profile.Category.Title, &equipmentByLoc.Company.Id, &equipmentByLoc.Company.Title)
		if err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}
	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationDepartment(toDepartment int) ([]model.Location, error) {
	var equipmentsByLoc []model.Location
	var equipmentByLoc model.Location
	var toEId, toEName interface{}
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
	rows, err := r.db.Query(context.Background(), query, toDepartment)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&equipmentByLoc.Equipment.Id, &equipmentByLoc.Equipment.SerialNumber, &equipmentByLoc.Equipment.Profile.Title, &equipmentByLoc.Equipment.Profile.Category.Title, &equipmentByLoc.Company.Id, &equipmentByLoc.Company.Title, &equipmentByLoc.ToDepartment.Id, &equipmentByLoc.ToDepartment.Title, &toEId, &toEName)
		equipmentByLoc.ToEmployee.Id = InterfaceToInt(toEId)
		equipmentByLoc.ToEmployee.Name = InterfaceToString(toEName)
		if err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}
	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationEmployee(toEmployee int) ([]model.Location, error) {
	var equipmentsByLoc []model.Location
	var equipmentByLoc model.Location
	var toD interface{}
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
	rows, err := r.db.Query(context.Background(), query, toEmployee)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&equipmentByLoc.Equipment.Id, &equipmentByLoc.Equipment.SerialNumber, &equipmentByLoc.Equipment.Profile.Title, &equipmentByLoc.Equipment.Profile.Category.Title, &equipmentByLoc.Company.Id, &equipmentByLoc.Company.Title, &toD, &equipmentByLoc.ToEmployee.Name)
		equipmentByLoc.ToDepartment.Title = InterfaceToString(toD)
		if err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}
	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationContract(toContract int) ([]model.Location, error) {
	var equipmentsByLoc []model.Location
	var equipmentByLoc model.Location
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
	rows, err := r.db.Query(context.Background(), query, toContract)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&equipmentByLoc.Equipment.Id, &equipmentByLoc.Equipment.SerialNumber, &equipmentByLoc.Equipment.Profile.Title, &equipmentByLoc.Equipment.Profile.Category.Title, &equipmentByLoc.Company.Id, &equipmentByLoc.Company.Title)
		if err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}
	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetByLocationDepartmentEmployee(toDepartment, toEmployee int) ([]model.Location, error) {
	var equipmentsByLoc []model.Location
	var equipmentByLoc model.Location
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
	rows, err := r.db.Query(context.Background(), query, toDepartment, toEmployee)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&equipmentByLoc.Equipment.Id, &equipmentByLoc.Equipment.SerialNumber, &equipmentByLoc.Equipment.Profile.Title, &equipmentByLoc.Equipment.Profile.Category.Title, &equipmentByLoc.Company.Id, &equipmentByLoc.Company.Title, &equipmentByLoc.ToDepartment.Title, &equipmentByLoc.ToEmployee.Name)
		if err != nil {
			return nil, err
		}
		equipmentsByLoc = append(equipmentsByLoc, equipmentByLoc)
	}
	return equipmentsByLoc, nil
}

func (r *EquipmentRepository) GetAll() ([]model.Equipment, error) {
	var equipments []model.Equipment
	var equipment model.Equipment
	query := `
			SELECT equipments.id, equipments.serial_number, profiles.id, profiles.title, categories.id, categories.title
			FROM equipments
			LEFT JOIN profiles ON profiles.id = equipments.profile
			LEFT JOIN categories ON categories.id = profiles.category
			WHERE equipments.is_deleted = FALSE
			ORDER BY profiles.title;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&equipment.Id, &equipment.SerialNumber, &equipment.Profile.Id, &equipment.Profile.Title, &equipment.Profile.Category.Id, &equipment.Profile.Category.Title)
		if err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}
	return equipments, nil
}

func (r *EquipmentRepository) FindBySerialNumber(serialNumber string) (int, error) {
	var equipment model.Equipment
	query := `
			SELECT id 
			FROM equipments 
			WHERE serial_number = $1;`
	err := r.db.QueryRow(context.Background(), query, serialNumber).Scan(&equipment.Id)
	if err != nil {
		return 0, err
	}
	return equipment.Id, nil
}

func (r *EquipmentRepository) Update(id int, serialNumber string, profile int) error {
	query := `
			UPDATE equipments 
			SET serial_number = $2, profile = $3
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, serialNumber, profile)
	if err != nil {
		return err
	}
	return nil
}

func (r *EquipmentRepository) Delete(id int) error {
	query := `
			UPDATE equipments 
			SET is_deleted = true
       		WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *EquipmentRepository) RemainderByCategory(categoryId, departmentId int, date time.Time) ([]model.Location, error) {
	var locations []model.Location
	var location model.Location
	var dateLoc time.Time
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
	rows, err := r.db.Query(context.Background(), query, categoryId, departmentId, date)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&dateLoc, &location.Equipment.Profile.Title, &location.Equipment.SerialNumber)
		location.Date = dateLoc.Unix()
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func (r *EquipmentRepository) TransferByCategory(categoryId, departmentId int, fromDate, toDate time.Time, code string) ([]model.Location, error) {
	var locations []model.Location
	var location model.Location
	var dateLoc time.Time
	query := ""
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
	rows, err := r.db.Query(context.Background(), query, categoryId, departmentId, code, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&dateLoc, &location.Equipment.Profile.Title, &location.Equipment.SerialNumber)
		location.Date = dateLoc.Unix()
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func (r *EquipmentRepository) ToDepartmentTransferByCategory(categoryId, departmentId int, fromDate, toDate time.Time) ([]model.Location, error) {
	var locations []model.Location
	var location model.Location
	var dateLoc time.Time
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
	rows, err := r.db.Query(context.Background(), query, categoryId, departmentId, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&dateLoc,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
			&location.ToDepartment.Id,
			&location.ToDepartment.Title)
		location.Date = dateLoc.Unix()
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func (r *EquipmentRepository) FromDepartmentTransferByCategory(categoryId, departmentId int, fromDate, toDate time.Time) ([]model.Location, error) {
	var locations []model.Location
	var location model.Location
	var dateLoc time.Time
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
	rows, err := r.db.Query(context.Background(), query, categoryId, departmentId, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&dateLoc,
			&location.Equipment.Profile.Title,
			&location.Equipment.SerialNumber,
			&location.FromDepartment.Id,
			&location.FromDepartment.Title)
		location.Date = dateLoc.Unix()
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}
