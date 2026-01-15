package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Auth       *AuthRepository
	User       *UserRepository
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

func New(postgresDB *pgxpool.Pool, redisDB *redis.Client, queries queries.Querier) *Repository {
	return &Repository{
		Auth:       NewAuthRepository(redisDB),
		User:       NewUserRepository(queries),
		Employee:   NewEmployeeRepository(queries),
		Department: NewDepartmentRepository(queries),
		Category:   NewCategoryRepository(queries),
		Profile:    NewProfileRepository(queries),
		Equipment:  NewEquipmentRepository(postgresDB),
		Location:   NewLocationRepository(postgresDB),
		Contract:   NewContractRepository(queries),
		Company:    NewCompanyRepository(queries),
		Replace:    NewReplaceRepository(postgresDB),
	}
}

type Auth interface {
	Get(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, claims *jwt_auth.CustomClaims, revoked bool) error
}

type User interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Read(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*model.User, error)
	GetPasswordHash(ctx context.Context, id int64) (string, error)
	SetPasswordHash(ctx context.Context, id int64, passwordHash string) error
	SetEnabled(ctx context.Context, id int64, enabled bool) error
	SetLastLoginAt(ctx context.Context, id int64) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type Employee interface {
	Create(ctx context.Context, employee *model.Employee) (int64, error)
	Read(ctx context.Context, id int64) (*model.Employee, error)
	Update(ctx context.Context, employee *model.Employee) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Employee, int64, error)
	SetDepartment(ctx context.Context, id, departmentID int64) error
	//GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error)
	//GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error)
	//GetFree(ctx context.Context) ([]*model.Employee, error)
	//GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error)
	//AddToDepartment(ctx context.Context, id, department int64) error
	//RemoveFromDepartment(ctx context.Context, id int64) error
	//Activate(ctx context.Context, id int64, password string) error
	//Deactivate(ctx context.Context, id int64) error
	//ResetPassword(ctx context.Context, id int64, password string) error
	//ChangeRole(ctx context.Context, id int64, role string) error
}

type Department interface {
	Create(ctx context.Context, department *model.Department) (int64, error)
	Read(ctx context.Context, id int64) (*model.Department, error)
	Update(ctx context.Context, department *model.Department) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Department, int64, error)
	//GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error)
	//GetAllButOneForAdmin(ctx context.Context, id int64) ([]*model.Department, error)
	//FindByTitle(ctx context.Context, title string) (int64, error)
}

type Category interface {
	Create(ctx context.Context, category *model.Category) (int64, error)
	Read(ctx context.Context, id int64) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Category, int64, error)
}

type Profile interface {
	Create(ctx context.Context, profile *model.Profile) (int64, error)
	Read(ctx context.Context, id int64) (*model.Profile, error)
	Update(ctx context.Context, profile *model.Profile) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Profile, int64, error)
}

type Equipment interface {
	Create(ctx context.Context, equipment *queries.CreateEquipmentParams, location *queries.AddToStorageParams) (int64, error)
	Read(ctx context.Context, id int64) (*model.Equipment, error)
	Update(ctx context.Context, equipment *model.Equipment) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Equipment, int64, error)
}

type Location interface {
	Move(ctx context.Context, location *queries.MoveToLocationParams) error
	List(ctx context.Context, toDepartmentID int64) ([]*model.Equipment, int64, error)
	//AddToStorage(ctx context.Context, date *time.Time, equipmentId, employeeId, companyId int64) error
	//TransferToStorage(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId int64, nowLocation []interface{}) (int64, error)
	//TransferToDepartment(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toDepartment int64, nowLocation []interface{}) (int64, error)
	//TransferToEmployee(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toEmployee int64, nowLocation []interface{}) (int64, error)
	//TransferToEmployeeInDepartment(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toDepartment, toEmployee int64, nowLocation []interface{}) (int64, error)
	//TransferToContract(ctx context.Context, date *time.Time, code string, equipmentId, employeeId, companyId, toContract int64, transferType string, price string, nowLocation []interface{}) (int64, error)
	//Delete(ctx context.Context, id int64) error
	//GetById(ctx context.Context, equipmentId int64) (*model.Location, error)
	//GetHistory(ctx context.Context, equipmentId int64) ([]*model.Location, error)
	//GetLocationNow(ctx context.Context, equipmentId int64) ([]interface{}, error)
	//GetByLocationStorage(ctx context.Context) ([]*model.Location, error)
	//GetByLocationDepartment(ctx context.Context, toDepartment int64) ([]*model.Location, error)
	//GetByLocationEmployee(ctx context.Context, toEmployee int64) ([]*model.Location, error)
	//GetByLocationContract(ctx context.Context, toContract int64) ([]*model.Location, error)
	//GetByLocationDepartmentEmployee(ctx context.Context, toDepartment, toEmployee int64) ([]*model.Location, error)
	//RemainderByCategory(ctx context.Context, categoryId, departmentId int64, date *time.Time) ([]*model.Location, error)
	//TransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate *time.Time, code string) ([]*model.Location, error)
	//ToDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate *time.Time) ([]*model.Location, error)
	//FromDepartmentTransferByCategory(ctx context.Context, categoryId, departmentId int64, fromDate, toDate *time.Time) ([]*model.Location, error)
}

type Contract interface {
	Create(ctx context.Context, contract *model.Contract) (int64, error)
	Read(ctx context.Context, id int64) (*model.Contract, error)
	Update(ctx context.Context, contract *model.Contract) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Contract, int64, error)
}

type Company interface {
	Create(ctx context.Context, company *model.Company) (int64, error)
	Read(ctx context.Context, id int64) (*model.Company, error)
	Update(ctx context.Context, company *model.Company) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Company, int64, error)
}

type Replace interface {
	Create(ctx context.Context, transferIds []int64) error
	FindByLocationId(ctx context.Context, locationId int64) (*model.Replace, error)
}

func validInt64(data pgtype.Int8) int64 {
	if data.Valid {
		return data.Int64
	}
	return 0
}

func validString(data pgtype.Text) string {
	if data.Valid {
		return data.String
	}
	return ""
}

func validTime(data pgtype.Timestamptz) *time.Time {
	if data.Valid {
		return &data.Time
	}
	return nil
}
