package service

import (
	"context"
	"time"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type Service struct {
	Auth       *AuthService
	Employee   *EmployeeService
	Department *DepartmentService
	Category   *CategoryService
	Profile    *ProfileService
	Equipment  *EquipmentService
	Location   *LocationService
	Contract   *ContractService
	Company    *CompanyService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Auth:       NewAuthService(repository.Auth),
		Employee:   NewEmployeeService(repository.Employee, repository.Equipment, repository.Auth),
		Department: NewDepartmentService(repository.Department, repository.Employee),
		Category:   NewCategoryService(repository.Category),
		Profile:    NewProfileService(repository.Profile),
		Equipment:  NewEquipmentService(repository.Equipment),
		Location:   NewLocationService(repository.Location, repository.Replace, repository.Category),
		Contract:   NewContractService(repository.Contract),
		Company:    NewCompanyService(repository.Company),
	}
}

type Auth interface {
	AuthUser(ctx context.Context, login, password string) (int64, error)
	GenerateHash(ctx context.Context, id int64) (string, error)
	FindByHash(ctx context.Context, hash string) (int64, error)
}

type Employee interface {
	Create(ctx context.Context, name, phone, email string) error
	Update(ctx context.Context, id int64, name, phone, email string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Employee, error)
	GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error)
	GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error)
	GetById(ctx context.Context, id int64) (*model.Employee, error)
	GetFree(ctx context.Context) ([]*model.Employee, error)
	GetByDepartment(ctx context.Context, ids []int64, id int64) ([]*model.Employee, error)
	AddToDepartment(ctx context.Context, id, department int64) error
	RemoveFromDepartment(ctx context.Context, idDepartment, idEmployee int64) error
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
	ResetPassword(ctx context.Context, id int64) error
	ChangeRole(ctx context.Context, id int64, role string) error
}

type Department interface {
	Create(ctx context.Context, title string) error
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Department, error)
	GetById(ctx context.Context, id int64) (*model.Department, error)
	GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error)
}

type Category interface {
	Create(ctx context.Context, title string) error
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Category, error)
	GetById(ctx context.Context, id int64) (*model.Category, error)
}

type Profile interface {
	Create(ctx context.Context, title string, categoryId int64) error
	Update(ctx context.Context, id int64, title string, categoryId int64) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Profile, error)
	GetById(ctx context.Context, id int64) (*model.Profile, error)
}

type Equipment interface {
	Create(ctx context.Context, serialNumber string, profileId int64) (int64, error)
	Update(ctx context.Context, id int64, serialNumber string, profileId int64) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*model.Equipment, error)
	GetByIds(ctx context.Context, ids []int64) ([]*model.Equipment, error)
	FindBySerialNumber(ctx context.Context, value string) ([]*model.Equipment, error)
}

type Location interface {
	AddToStorage(ctx context.Context, date *time.Time, equipmentId, employeeId, companyId int64) error
	TransferTo(ctx context.Context, EmployeeId int64, requests []*model.RequestLocation) error
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, equipmentId int64) (*model.Location, error)
	GetByIds(ctx context.Context, equipmentIds []int64) ([]*model.Location, error)
	GetHistory(ctx context.Context, equipmentId int64) ([]*model.Location, error)
	GetByLocation(ctx context.Context, toDepartment, toEmployee, toContract int64) ([]*model.Location, error)
	ReportByCategory(ctx context.Context, departmentId int64, date *time.Time) (*model.Report, error)
}

type Contract interface {
	Create(ctx context.Context, number, address string) error
	Update(ctx context.Context, id int64, number, address string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Contract, error)
	GetById(ctx context.Context, id int64) (*model.Contract, error)
}

type Company interface {
	Create(ctx context.Context, title string) error
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Company, error)
	GetById(ctx context.Context, id int64) (*model.Company, error)
}
