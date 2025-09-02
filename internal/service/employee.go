package service

import (
	"context"

	"github.com/oatsmoke/warehouse_backend/internal/lib/email"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService struct {
	EmployeeRepository repository.Employee
	//repositoryEquipment repository.Equipment
	//repositoryAuth      repository.Auth
}

func NewEmployeeService(
	employeeRepository repository.Employee,
	repositoryEquipment repository.Equipment,
	repositoryAuth repository.Auth) *EmployeeService {
	return &EmployeeService{
		EmployeeRepository: employeeRepository,
		//repositoryEquipment: repositoryEquipment,
		//repositoryAuth:      repositoryAuth,
	}
}

// Create is employee create
func (s *EmployeeService) Create(ctx context.Context, name, phone, email string) error {
	if err := s.EmployeeRepository.Create(ctx, name, phone, email); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Update is employee update
func (s *EmployeeService) Update(ctx context.Context, id int64, name, phone, email string) error {
	if err := s.EmployeeRepository.Update(ctx, id, name, phone, email); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Delete is employee delete
func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	if err := s.EmployeeRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Restore is employee restore
func (s *EmployeeService) Restore(ctx context.Context, id int64) error {
	if err := s.EmployeeRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// GetAll is to get all employees
func (s *EmployeeService) GetAll(ctx context.Context, deleted bool) ([]*model.Employee, error) {
	res, err := s.EmployeeRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetAllShort is to get all employees short
func (s *EmployeeService) GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error) {
	res, err := s.EmployeeRepository.GetAllShort(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetAllButOne is to get all employees but one
func (s *EmployeeService) GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error) {
	res, err := s.EmployeeRepository.GetAllButOne(ctx, id, deleted)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetById is employee get by id
func (s *EmployeeService) GetById(ctx context.Context, id int64) (*model.Employee, error) {
	employee := &model.Employee{
		ID: id,
	}

	res, err := s.EmployeeRepository.GetById(ctx, employee)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetFree is employee get free
func (s *EmployeeService) GetFree(ctx context.Context) ([]*model.Employee, error) {
	res, err := s.EmployeeRepository.GetFree(ctx)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetByDepartment is employee get by department
func (s *EmployeeService) GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error) {
	res, err := s.EmployeeRepository.GetByDepartment(ctx, ids, departmentId)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// AddToDepartment is employee add to department
func (s *EmployeeService) AddToDepartment(ctx context.Context, id, department int64) error {
	if err := s.EmployeeRepository.AddToDepartment(ctx, id, department); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// RemoveFromDepartment is employee remove from department
func (s *EmployeeService) RemoveFromDepartment(ctx context.Context, idDepartment, idEmployee int64) error {
	if err := s.EmployeeRepository.RemoveFromDepartment(ctx, idEmployee); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Activate is employee activate
func (s *EmployeeService) Activate(ctx context.Context, id int64) error {
	employee := &model.Employee{
		ID: id,
	}

	res, err := s.EmployeeRepository.GetById(ctx, employee)
	if err != nil {
		return logger.Err(err, "")
	}

	str := generate.RandString(10)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return logger.Err(err, "")
	}

	if err := s.EmployeeRepository.Activate(ctx, id, string(passwordHash)); err != nil {
		return logger.Err(err, "")
	}

	sendTo := &email.SendTo{
		Name:     res.Name,
		Email:    res.Email,
		Phone:    res.Phone,
		Password: string(passwordHash),
	}

	if err := email.Send([]*email.SendTo{sendTo}); err != nil {
		return logger.Err(err, "")
	}

	//if err := sendMail(res.Email, res.Phone, string(passwordHash)); err != nil {
	//	return logger.Err(err, "")
	//}

	return nil
}

// Deactivate is employee deactivate
func (s *EmployeeService) Deactivate(ctx context.Context, id int64) error {
	if err := s.EmployeeRepository.Deactivate(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// ResetPassword is employee reset password
func (s *EmployeeService) ResetPassword(ctx context.Context, id int64) error {
	employee := &model.Employee{
		ID: id,
	}

	res, err := s.EmployeeRepository.GetById(ctx, employee)
	if err != nil {
		return logger.Err(err, "")
	}

	str := generate.RandString(10)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return logger.Err(err, "")
	}

	if err := s.EmployeeRepository.ResetPassword(ctx, id, string(passwordHash)); err != nil {
		return logger.Err(err, "")
	}

	sendTo := &email.SendTo{
		Name:     res.Name,
		Email:    res.Email,
		Phone:    res.Phone,
		Password: string(passwordHash),
	}

	if err := email.Send([]*email.SendTo{sendTo}); err != nil {
		return logger.Err(err, "")
	}

	//if err := sendMail(res.Email, res.Phone, string(passwordHash)); err != nil {
	//	return logger.Err(err, "")
	//}

	return nil
}

// ChangeRole is employee change role
func (s *EmployeeService) ChangeRole(ctx context.Context, id int64, role string) error {
	if err := s.EmployeeRepository.ChangeRole(ctx, id, role); err != nil {
		return logger.Err(err, "")
	}

	return nil
}
