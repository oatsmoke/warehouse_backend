package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/pkg/model"
)

type Repository struct {
	Employee
	Department
	Category
	Profile
	Equipment
	Location
	Contract
	Company
	Replace
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Employee:   NewEmployeeRepository(db),
		Department: NewDepartmentRepository(db),
		Category:   NewCategoryRepository(db),
		Profile:    NewProfileRepository(db),
		Equipment:  NewEquipmentRepository(db),
		Location:   NewLocationRepository(db),
		Contract:   NewContractRepository(db),
		Company:    NewCompanyRepository(db),
		Replace:    NewReplaceRepository(db),
	}
}

type Employee interface {
	Create(name, phone, email string) error
	GetById(id int) (model.Employee, error)
	GetByDepartment(ids []int, id int) ([]model.Employee, error)
	GetAll() ([]model.Employee, error)
	GetFree() ([]model.Employee, error)
	GetAllButOne(id int) ([]model.Employee, error)
	FindUser(userName, password string) (int, error)
	FindByPhone(phone string) (int, error)
	FindByHash(hash string) (int, error)
	SetHash(id int, hash string) error
	AddToDepartment(id, department int) error
	RemoveFromDepartment(id int) error
	Update(id int, name, phone, email string) error
	Delete(id int) error
	Activate(id int, password string) error
	Deactivate(id int) error
	ResetPassword(id int, password string) error
	ChangeRole(id int, role string) error
}

type Department interface {
	Create(title string) error
	GetById(id int) (model.Department, error)
	GetAll() ([]model.Department, error)
	GetAllButOne(id, employeeId int) ([]model.Department, error)
	GetAllButOneForAdmin(id int) ([]model.Department, error)
	FindByTitle(title string) (int, error)
	Update(id int, title string) error
	Delete(id int) error
}

type Category interface {
	Create(title string) error
	GetById(id int) (model.Category, error)
	GetAll() ([]model.Category, error)
	FindByTitle(title string) (int, error)
	Update(id int, title string) error
	Delete(id int) error
}

type Profile interface {
	Create(title string, category int) error
	GetByCategory(id int) ([]model.Profile, error)
	GetById(id int) (model.Profile, error)
	GetAll() ([]model.Profile, error)
	FindByTitle(title string) (int, error)
	Update(id int, title string, category int) error
	Delete(id int) error
}

type Equipment interface {
	Create(date int64, company int, serialNumber string, profile int, userId int) (int, error)
	GetByProfile(id int) ([]model.Equipment, error)
	GetById(id int) (model.Location, error)
	GetByLocationStorage() ([]model.Location, error)
	GetByLocationDepartment(toDepartment int) ([]model.Location, error)
	GetByLocationEmployee(toEmployee int) ([]model.Location, error)
	GetByLocationContract(toContract int) ([]model.Location, error)
	GetByLocationDepartmentEmployee(toDepartment, toEmployee int) ([]model.Location, error)
	GetAll() ([]model.Equipment, error)
	FindBySerialNumber(serialNumber string) (int, error)
	Update(id int, serialNumber string, profile int) error
	Delete(id int) error
	RemainderByCategory(categoryId, departmentId int, date time.Time) ([]model.Location, error)
	TransferByCategory(categoryId, departmentId int, fromDate, toDate time.Time, code string) ([]model.Location, error)
	ToDepartmentTransferByCategory(categoryId, departmentId int, fromDate, toDate time.Time) ([]model.Location, error)
	FromDepartmentTransferByCategory(categoryId, departmentId int, fromDate, toDate time.Time) ([]model.Location, error)
}

type Location interface {
	TransferToStorage(date int64, code string, equipment, employee, company int, nowLocation []interface{}) (int, error)
	TransferToDepartment(date int64, code string, equipment, employee, company, toDepartment int, nowLocation []interface{}) (int, error)
	TransferToEmployee(date int64, code string, equipment, employee, company, toEmployee int, nowLocation []interface{}) (int, error)
	TransferToEmployeeInDepartment(date int64, code string, equipment, employee, company, toDepartment, toEmployee int, nowLocation []interface{}) (int, error)
	TransferToContract(date int64, code string, equipment, employee, company, toContract int, transferType string, price int, nowLocation []interface{}) (int, error)
	GetHistory(id int) ([]model.Location, error)
	GetLocationNow(id int) ([]interface{}, error)
	Delete(id int) error
}

type Contract interface {
	Create(number, address string) error
	GetById(id int) (model.Contract, error)
	GetAll() ([]model.Contract, error)
	FindByNumber(number string) (int, error)
	Update(id int, number, address string) error
	Delete(id int) error
}

type Company interface {
	Create(title string) error
	GetById(id int) (model.Company, error)
	GetAll() ([]model.Company, error)
	FindByTitle(title string) (int, error)
	Update(id int, title string) error
	Delete(id int) error
}

type Replace interface {
	Create(ids []int) error
	FindByLocationId(id int) (model.Replace, error)
}

func InterfaceToInt(value interface{}) int {
	if value != nil {
		return int(value.(int32))
	}
	return 0
}

func InterfaceToString(value interface{}) string {
	if value != nil {
		return value.(string)
	}
	return ""
}
