package service

import (
	"context"
	"errors"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type DepartmentService struct {
	repositoryDepartment repository.Department
	repositoryEquipment  repository.Equipment
	repositoryEmployee   repository.Employee
}

func NewDepartmentService(repositoryDepartment repository.Department, repositoryEquipment repository.Equipment, repositoryEmployee repository.Employee) *DepartmentService {
	return &DepartmentService{
		repositoryDepartment: repositoryDepartment,
		repositoryEquipment:  repositoryEquipment,
		repositoryEmployee:   repositoryEmployee,
	}
}

func (s *DepartmentService) Create(ctx context.Context, title string) error {
	if _, err := s.repositoryDepartment.FindByTitle(ctx, title); err == nil {
		return errors.New("title already exists")
	}

	return s.repositoryDepartment.Create(ctx, title)
}

func (s *DepartmentService) GetById(ctx context.Context, id int64) (*model.Department, error) {
	return s.repositoryDepartment.GetById(ctx, id)
}

func (s *DepartmentService) GetAll(ctx context.Context) ([]*model.Department, error) {
	return s.repositoryDepartment.GetAll(ctx)
}

func (s *DepartmentService) GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error) {
	employee, err := s.repositoryEmployee.GetById(ctx, employeeId)
	if err != nil {
		return nil, err
	}

	if employee.Role == "ADMIN" {
		return s.repositoryDepartment.GetAllButOneForAdmin(ctx, id)
	}

	return s.repositoryDepartment.GetAllButOne(ctx, id, employeeId)
}

func (s *DepartmentService) Update(ctx context.Context, id int64, title string) error {
	if _, err := s.repositoryDepartment.FindByTitle(ctx, title); err == nil {
		return errors.New("title already exists")
	}

	return s.repositoryDepartment.Update(ctx, id, title)
}

func (s *DepartmentService) Delete(ctx context.Context, id int64) error {
	equipments, err := s.repositoryEquipment.GetByLocationDepartment(ctx, id)
	if err != nil {
		return err
	}

	if len(equipments) > 0 {
		return errors.New("have equipment")
	}

	employees, err := s.repositoryEmployee.GetByDepartment(ctx, []int64{}, id)
	if err != nil {
		return err
	}

	if len(employees) > 0 {
		return errors.New("have employees")
	}

	return s.repositoryDepartment.Delete(ctx, id)
}
