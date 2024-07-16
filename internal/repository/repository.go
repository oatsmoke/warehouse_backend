package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/internal/model"
)

type Repository struct {
	Auth       *AuthRepository
	Employee   *EmployeeRepository
	Department *DepartmentRepository
	Category   *CategoryRepository
	Profile    *ProfileRepository
	Equipment  *EquipmentRepository
	Location   *LocationRepository
	Contract   *ContractRepository
	Company    *CompanyRepository
	Replace    *ReplaceRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Auth:       NewAuthRepository(db),
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

type Auth interface {
	FindByPhone(ctx context.Context, user *model.Employee) (*model.Employee, error)
	SetHash(ctx context.Context, id int64, hash string) error
	FindByHash(ctx context.Context, user *model.Employee) (*model.Employee, error)
}

type Employee interface {
	Create(ctx context.Context, name, phone, email string) error
	GetById(ctx context.Context, id int64) (*model.Employee, error)
	GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error)
	GetAll(ctx context.Context) ([]*model.Employee, error)
	GetFree(ctx context.Context) ([]*model.Employee, error)
	GetAllButOne(ctx context.Context, id int64) ([]*model.Employee, error)
	AddToDepartment(ctx context.Context, id, department int64) error
	RemoveFromDepartment(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, name, phone, email string) error
	Delete(ctx context.Context, id int64) error
	Activate(ctx context.Context, id int64, password string) error
	Deactivate(ctx context.Context, id int64) error
	ResetPassword(ctx context.Context, id int64, password string) error
	ChangeRole(ctx context.Context, id int64, role string) error
}

type Department interface {
	Create(ctx context.Context, title string) error
	GetById(ctx context.Context, id int64) (*model.Department, error)
	GetAll(ctx context.Context) ([]*model.Department, error)
	GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error)
	GetAllButOneForAdmin(ctx context.Context, id int64) ([]*model.Department, error)
	FindByTitle(ctx context.Context, title string) (int64, error)
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
}

type Category interface {
	Create(ctx context.Context, title string) error
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, isDeleted bool) ([]*model.Category, error)
	GetById(ctx context.Context, category *model.Category) (*model.Category, error)
	//FindByTitle(ctx context.Context, title string) (int64, error)
}

type Profile interface {
	Create(ctx context.Context, title string, category int64) error
	GetByCategory(ctx context.Context, id int64) ([]*model.Profile, error)
	GetById(ctx context.Context, id int64) (*model.Profile, error)
	GetAll(ctx context.Context) ([]*model.Profile, error)
	FindByTitle(ctx context.Context, title string) (int64, error)
	Update(ctx context.Context, id int64, title string, category int64) error
	Delete(ctx context.Context, id int64) error
}

type Equipment interface {
	Create(ctx context.Context, date int64, company int64, serialNumber string, profile int64, userId int64) (int64, error)
	GetByProfile(ctx context.Context, id int64) ([]*model.Equipment, error)
	GetById(ctx context.Context, id int64) (*model.Location, error)
	GetByLocationStorage(ctx context.Context) ([]*model.Location, error)
	GetByLocationDepartment(ctx context.Context, toDepartment int64) ([]*model.Location, error)
	GetByLocationEmployee(ctx context.Context, toEmployee int64) ([]*model.Location, error)
	GetByLocationContract(ctx context.Context, toContract int64) ([]*model.Location, error)
	GetByLocationDepartmentEmployee(ctx context.Context, toDepartment, toEmployee int64) ([]*model.Location, error)
	GetAll(ctx context.Context) ([]*model.Equipment, error)
	FindBySerialNumber(ctx context.Context, serialNumber string) (int64, error)
	Update(ctx context.Context, id int64, serialNumber string, profile int64) error
	Delete(ctx context.Context, id int64) error
	RemainderByCategory(ctx context.Context, categoryId, departmentId int64, date time.Time) ([]*model.Location, error)
	TransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate time.Time, code string) ([]*model.Location, error)
	ToDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate time.Time) ([]*model.Location, error)
	FromDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate time.Time) ([]*model.Location, error)
}

type Location interface {
	TransferToStorage(ctx context.Context, date int64, code string, equipment, employee, company int64, nowLocation []interface{}) (int64, error)
	TransferToDepartment(ctx context.Context, date int64, code string, equipment, employee, company, toDepartment int64, nowLocation []interface{}) (int64, error)
	TransferToEmployee(ctx context.Context, date int64, code string, equipment, employee, company, toEmployee int64, nowLocation []interface{}) (int64, error)
	TransferToEmployeeInDepartment(ctx context.Context, date int64, code string, equipment, employee, company, toDepartment, toEmployee int64, nowLocation []interface{}) (int64, error)
	TransferToContract(ctx context.Context, date int64, code string, equipment, employee, company, toContract int64, transferType string, price int, nowLocation []interface{}) (int64, error)
	GetHistory(ctx context.Context, id int64) ([]*model.Location, error)
	GetLocationNow(ctx context.Context, id int64) ([]interface{}, error)
	Delete(ctx context.Context, id int64) error
}

type Contract interface {
	Create(ctx context.Context, number, address string) error
	GetById(ctx context.Context, id int64) (*model.Contract, error)
	GetAll(ctx context.Context) ([]*model.Contract, error)
	FindByNumber(ctx context.Context, number string) (int64, error)
	Update(ctx context.Context, id int64, number, address string) error
	Delete(ctx context.Context, id int64) error
}

type Company interface {
	Create(ctx context.Context, title string) error
	GetById(ctx context.Context, id int64) (*model.Company, error)
	GetAll(ctx context.Context) ([]*model.Company, error)
	FindByTitle(ctx context.Context, title string) (int64, error)
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
}

type Replace interface {
	Create(ctx context.Context, transferIds []int64) error
	FindByLocationId(ctx context.Context, id int64) (*model.Replace, error)
}