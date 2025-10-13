package service

import (
	"context"
	"time"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type Service struct {
	Auth       *AuthService
	User       *UserService
	Employee   *EmployeeService
	Department *DepartmentService
	Category   *CategoryService
	Profile    *ProfileService
	Equipment  *EquipmentService
	Location   *LocationService
	Contract   *ContractService
	Company    *CompanyService
}

func New(repository *repository.Repository) *Service {
	return &Service{
		Auth:       NewAuthService(repository.Auth, repository.User),
		User:       NewUserService(repository.User, repository.Employee),
		Employee:   NewEmployeeService(repository.Employee),
		Department: NewDepartmentService(repository.Department),
		Category:   NewCategoryService(repository.Category),
		Profile:    NewProfileService(repository.Profile),
		Equipment:  NewEquipmentService(repository.Equipment),
		Location:   NewLocationService(repository.Location, repository.Replace, repository.Category),
		Contract:   NewContractService(repository.Contract),
		Company:    NewCompanyService(repository.Company),
	}
}

type Auth interface {
	AuthUser(ctx context.Context, login *dto.UserLogin) (*jwt_auth.Token, error)
	CheckToken(ctx context.Context, token *jwt_auth.Token) (*jwt_auth.Token, error)
}

type User interface {
	Create(ctx context.Context, user *model.User) error
	Read(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*model.User, error)
	SetPassword(ctx context.Context, id int64, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, id int64) error
	SetRole(ctx context.Context, id int64, role role.Role) error
	SetEnabled(ctx context.Context, id int64, enabled bool) error
	SetEmployee(ctx context.Context, id, employeeID int64) error
}

type Employee interface {
	Create(ctx context.Context, employee *model.Employee) error
	Read(ctx context.Context, id int64) (*model.Employee, error)
	Update(ctx context.Context, employee *model.Employee) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Employee, error)
	SetDepartment(ctx context.Context, id, departmentID int64) error
	//GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error)
	//GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error)
	//GetFree(ctx context.Context) ([]*model.Employee, error)
	//GetByDepartment(ctx context.Context, ids []int64, id int64) ([]*model.Employee, error)
	//AddToDepartment(ctx context.Context, id, department int64) error
	//RemoveFromDepartment(ctx context.Context, idDepartment, idEmployee int64) error
	//Activate(ctx context.Context, id int64) error
	//Deactivate(ctx context.Context, id int64) error
	//ResetPassword(ctx context.Context, id int64) error
	//ChangeRole(ctx context.Context, id int64, role string) error
}

type Department interface {
	Create(ctx context.Context, department *model.Department) error
	Read(ctx context.Context, id int64) (*model.Department, error)
	Update(ctx context.Context, department *model.Department) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Department, error)
	//GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error)
}

type Category interface {
	Create(ctx context.Context, category *model.Category) error
	Read(ctx context.Context, id int64) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Category, error)
}

type Profile interface {
	Create(ctx context.Context, profile *model.Profile) error
	Read(ctx context.Context, id int64) (*model.Profile, error)
	Update(ctx context.Context, profile *model.Profile) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Profile, error)
}

type Equipment interface {
	Create(ctx context.Context, equipment *model.Equipment) error
	Read(ctx context.Context, id int64) (*model.Equipment, error)
	Update(ctx context.Context, equipment *model.Equipment) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Equipment, error)
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
	Create(ctx context.Context, contract *model.Contract) error
	Read(ctx context.Context, id int64) (*model.Contract, error)
	Update(ctx context.Context, contract *model.Contract) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Contract, error)
}

type Company interface {
	Create(ctx context.Context, company *model.Company) error
	Read(ctx context.Context, id int64) (*model.Company, error)
	Update(ctx context.Context, company *model.Company) error
	Delete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, qp *dto.QueryParams) ([]*model.Company, error)
}
