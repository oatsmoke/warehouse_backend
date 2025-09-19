package repository

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/redis/go-redis/v9"
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

func New(postgresDB *pgxpool.Pool, redisDB *redis.Client) *Repository {
	return &Repository{
		Auth:       NewAuthRepository(postgresDB, redisDB),
		Employee:   NewEmployeeRepository(postgresDB),
		Department: NewDepartmentRepository(postgresDB),
		Category:   NewCategoryRepository(postgresDB),
		Profile:    NewProfileRepository(postgresDB),
		Equipment:  NewEquipmentRepository(postgresDB),
		Location:   NewLocationRepository(postgresDB),
		Contract:   NewContractRepository(postgresDB),
		Company:    NewCompanyRepository(postgresDB),
		Replace:    NewReplaceRepository(postgresDB),
	}
}

type Auth interface {
	FindByPhone(ctx context.Context, user *model.Employee) (*model.Employee, error)
	Set(ctx context.Context, claims *jwt.RegisteredClaims, revoked bool) error
	Get(ctx context.Context, key string) (bool, error)
	//SetHash(ctx context.Context, id int64, hash string) error
	//FindByHash(ctx context.Context, user *model.Employee) (*model.Employee, error)
}

type Employee interface {
	Create(ctx context.Context, name, phone, email string) error
	Update(ctx context.Context, id int64, name, phone, email string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Employee, error)
	GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error)
	GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error)
	GetById(ctx context.Context, employee *model.Employee) (*model.Employee, error)
	GetFree(ctx context.Context) ([]*model.Employee, error)
	GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error)
	AddToDepartment(ctx context.Context, id, department int64) error
	RemoveFromDepartment(ctx context.Context, id int64) error
	Activate(ctx context.Context, id int64, password string) error
	Deactivate(ctx context.Context, id int64) error
	ResetPassword(ctx context.Context, id int64, password string) error
	ChangeRole(ctx context.Context, id int64, role string) error
}

type Department interface {
	Create(ctx context.Context, title string) error
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Department, error)
	GetById(ctx context.Context, department *model.Department) (*model.Department, error)
	GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error)
	GetAllButOneForAdmin(ctx context.Context, id int64) ([]*model.Department, error)
	//FindByTitle(ctx context.Context, title string) (int64, error)
}

type Category interface {
	Create(ctx context.Context, category *model.Category) error
	Read(ctx context.Context, id int64) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, withDeleted bool) ([]*model.Category, error)
}

type Profile interface {
	Create(ctx context.Context, title string, categoryId int64) error
	Update(ctx context.Context, id int64, title string, categoryId int64) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Profile, error)
	GetById(ctx context.Context, profile *model.Profile) (*model.Profile, error)
	//GetByCategory(ctx context.Context, categoryId int64) ([]*model.Profile, error)
	//FindByTitle(ctx context.Context, title string) (int64, error)
}

type Equipment interface {
	Create(ctx context.Context, serialNumber string, profileId int64) (int64, error)
	Update(ctx context.Context, id int64, serialNumber string, profileId int64) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*model.Equipment, error)
	GetByIds(ctx context.Context, ids []int64) ([]*model.Equipment, error)
	FindBySerialNumber(ctx context.Context, value string) ([]*model.Equipment, error)
	//GetByProfile(ctx context.Context, id int64) ([]*model.Equipment, error)
	//GetBySerialNumber(ctx context.Context, equipment *model.Equipment) (*model.Equipment, error)
}

type Location interface {
	AddToStorage(ctx context.Context, date *time.Time, equipmentId, employeeId, companyId int64) error
	TransferToStorage(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId int64, nowLocation []interface{}) (int64, error)
	TransferToDepartment(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toDepartment int64, nowLocation []interface{}) (int64, error)
	TransferToEmployee(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toEmployee int64, nowLocation []interface{}) (int64, error)
	TransferToEmployeeInDepartment(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toDepartment, toEmployee int64, nowLocation []interface{}) (int64, error)
	TransferToContract(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toContract int64, transferType string, price string, nowLocation []interface{}) (int64, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, equipmentId int64) (*model.Location, error)
	GetHistory(ctx context.Context, equipmentId int64) ([]*model.Location, error)
	GetLocationNow(ctx context.Context, equipmentId int64) ([]interface{}, error)
	GetByLocationStorage(ctx context.Context) ([]*model.Location, error)
	GetByLocationDepartment(ctx context.Context, toDepartment int64) ([]*model.Location, error)
	GetByLocationEmployee(ctx context.Context, toEmployee int64) ([]*model.Location, error)
	GetByLocationContract(ctx context.Context, toContract int64) ([]*model.Location, error)
	GetByLocationDepartmentEmployee(ctx context.Context, toDepartment, toEmployee int64) ([]*model.Location, error)
	RemainderByCategory(ctx context.Context, categoryId, departmentId int64, date *time.Time) ([]*model.Location, error)
	TransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate *time.Time, code string) ([]*model.Location, error)
	ToDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate *time.Time) ([]*model.Location, error)
	FromDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate *time.Time) ([]*model.Location, error)
}

type Contract interface {
	Create(ctx context.Context, number, address string) error
	Update(ctx context.Context, id int64, number, address string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Contract, error)
	GetById(ctx context.Context, contract *model.Contract) (*model.Contract, error)
	//FindByNumber(ctx context.Context, number string) (int64, error)
}

type Company interface {
	Create(ctx context.Context, title string) error
	Update(ctx context.Context, id int64, title string) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, deleted bool) ([]*model.Company, error)
	GetById(ctx context.Context, company *model.Company) (*model.Company, error)
	//FindByTitle(ctx context.Context, title string) (int64, error)
}

type Replace interface {
	Create(ctx context.Context, transferIds []int64) error
	FindByLocationId(ctx context.Context, locationId int64) (*model.Replace, error)
}
