package service

import (
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type Service struct {
	Employee
	Department
	Category
	Profile
	Equipment
	Location
	Contract
	Company
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Employee:   NewEmployeeService(repository.Employee, repository.Equipment),
		Department: NewDepartmentService(repository.Department, repository.Equipment, repository.Employee),
		Category:   NewCategoryService(repository.Category, repository.Profile),
		Profile:    NewProfileService(repository.Profile, repository.Equipment),
		Equipment:  NewEquipmentService(repository.Equipment),
		Location:   NewLocationService(repository.Location),
		Contract:   NewContractService(repository.Contract, repository.Equipment),
		Company:    NewCompanyService(repository.Company),
	}
}

type Employee interface {
	Create(name, phone, email string) error
	GetById(id int) (model.Employee, error)
	GetByDepartment(ids []int, id int) ([]model.Employee, error)
	GetAll() ([]model.Employee, error)
	GetFree() ([]model.Employee, error)
	GetAllButOne(id int) ([]model.Employee, error)
	FindUser(login, password string) (int, error)
	FindByHash(hash string) (int, error)
	AddToDepartment(id, department int) error
	RemoveFromDepartment(idDepartment, idEmployee int) error
	Update(id int, name, phone, email string) error
	Delete(id int) error
	GenerateToken(id int) (string, error)
	ParseToken(token string) (interface{}, error)
	GenerateHash(id int) (string, error)
	Activate(id int) error
	Deactivate(id int) error
	ResetPassword(id int) error
}

type Department interface {
	Create(title string) error
	GetById(id int) (model.Department, error)
	GetAll() ([]model.Department, error)
	GetAllButOne(id int) ([]model.Department, error)
	Update(id int, title string) error
	Delete(id int) error
}

type Category interface {
	Create(title string) error
	GetById(id int) (model.Category, error)
	GetAll() ([]model.Category, error)
	Update(id int, title string) error
	Delete(id int) error
}

type Profile interface {
	Create(title string, category int) error
	GetById(id int) (model.Profile, error)
	GetAll() ([]model.Profile, error)
	Update(id int, title string, category int) error
	Delete(id int) error
}

type Equipment interface {
	Create(date int64, company int, serialNumber string, profile int, userId int) (int, error)
	GetById(id int) (model.Location, error)
	GetByIds(ids []int) ([]model.Location, error)
	GetByLocation(toDepartment, toEmployee, toContract int) ([]model.Location, error)
	GetAll() ([]model.Equipment, error)
	Update(id int, serialNumber string, profile int) error
	Delete(id int) error
}
type Location interface {
	TransferTo(id int, request []model.RequestLocation) error
	GetHistory(id int) ([]model.Location, error)
	Delete(id int) error
}

type Contract interface {
	Create(number, address string) error
	GetById(id int) (model.Contract, error)
	GetAll() ([]model.Contract, error)
	Update(id int, number, address string) error
	Delete(id int) error
}

type Company interface {
	Create(title string) error
	GetById(id int) (model.Company, error)
	GetAll() ([]model.Company, error)
	Update(id int, title string) error
	Delete(id int) error
}
