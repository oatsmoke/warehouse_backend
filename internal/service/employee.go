package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"warehouse_backend/internal/lib/generate"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

//const (
//	salt       = "12345678"
//	dictionary = "abcdefghijklmnopqrstuvwxyz0123456789"
//	length     = 10
//)

type EmployeeService struct {
	repositoryEmployee  repository.Employee
	repositoryEquipment repository.Equipment
	repositoryAuth      repository.Auth
}

func NewEmployeeService(
	repositoryEmployee repository.Employee,
	repositoryEquipment repository.Equipment,
	repositoryAuth repository.Auth) *EmployeeService {
	return &EmployeeService{
		repositoryEmployee:  repositoryEmployee,
		repositoryEquipment: repositoryEquipment,
		repositoryAuth:      repositoryAuth,
	}
}

func (s *EmployeeService) Create(ctx context.Context, name, phone, email string) error {
	//if _, err := s.repositoryAuth.FindByPhone(ctx, phone); err == nil {
	//	return errors.New("phone already exists")
	//}

	return s.repositoryEmployee.Create(ctx, name, phone, email)
}

func (s *EmployeeService) GetById(ctx context.Context, id int64) (*model.Employee, error) {
	return s.repositoryEmployee.GetById(ctx, id)
}

func (s *EmployeeService) GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error) {
	return s.repositoryEmployee.GetByDepartment(ctx, ids, departmentId)
}

func (s *EmployeeService) GetAll(ctx context.Context) ([]*model.Employee, error) {
	return s.repositoryEmployee.GetAll(ctx)
}

func (s *EmployeeService) GetFree(ctx context.Context) ([]*model.Employee, error) {
	return s.repositoryEmployee.GetFree(ctx)
}

func (s *EmployeeService) GetAllButOne(ctx context.Context, id int64) ([]*model.Employee, error) {
	return s.repositoryEmployee.GetAllButOne(ctx, id)
}

func (s *EmployeeService) AddToDepartment(ctx context.Context, id, department int64) error {
	return s.repositoryEmployee.AddToDepartment(ctx, id, department)
}

func (s *EmployeeService) RemoveFromDepartment(ctx context.Context, idDepartment, idEmployee int64) error {
	equipments, err := s.repositoryEquipment.GetByLocationDepartmentEmployee(ctx, idDepartment, idEmployee)
	if err != nil {
		return err
	}

	if len(equipments) > 0 {
		return errors.New("have equipment")
	}

	return s.repositoryEmployee.RemoveFromDepartment(ctx, idEmployee)
}

func (s *EmployeeService) Update(ctx context.Context, id int64, name, phone, email string) error {
	//findId, err := s.repositoryAuth.FindByPhone(ctx, phone)
	//if findId != id && err == nil {
	//	return errors.New("phone already exists")
	//}

	return s.repositoryEmployee.Update(ctx, id, name, phone, email)
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	equipments, err := s.repositoryEquipment.GetByLocationEmployee(ctx, id)
	if err != nil {
		return err
	}

	if len(equipments) > 0 {
		return errors.New("have equipment")
	}

	return s.repositoryEmployee.Delete(ctx, id)
}

func (s *EmployeeService) Activate(ctx context.Context, id int64) error {
	employee, err := s.repositoryEmployee.GetById(ctx, id)
	if err != nil {
		return err
	}

	str := generate.RandString(10)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repositoryEmployee.Activate(ctx, id, string(passwordHash)); err != nil {
		return err
	}

	if err := sendMail(employee.Email, employee.Phone, string(passwordHash)); err != nil {
		return err
	}

	return nil
}

func (s *EmployeeService) Deactivate(ctx context.Context, id int64) error {
	return s.repositoryEmployee.Deactivate(ctx, id)
}

func (s *EmployeeService) ResetPassword(ctx context.Context, id int64) error {
	employee, err := s.repositoryEmployee.GetById(ctx, id)
	if err != nil {
		return err
	}

	str := generate.RandString(10)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repositoryEmployee.ResetPassword(ctx, id, string(passwordHash)); err != nil {
		return err
	}

	if err := sendMail(employee.Email, employee.Phone, string(passwordHash)); err != nil {
		return err
	}

	return nil
}

func (s *EmployeeService) ChangeRole(ctx context.Context, id int64, role string) error {
	return s.repositoryEmployee.ChangeRole(ctx, id, role)
}
