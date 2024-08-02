package service

import (
	"context"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type DepartmentService struct {
	DepartmentRepository repository.Department
	EmployeeRepository   repository.Employee
}

func NewDepartmentService(departmentRepository repository.Department, employeeRepository repository.Employee) *DepartmentService {
	return &DepartmentService{
		DepartmentRepository: departmentRepository,
		EmployeeRepository:   employeeRepository,
	}
}

// Create is department create
func (s *DepartmentService) Create(ctx context.Context, title string) error {
	const fn = "service.Department.Create"

	if err := s.DepartmentRepository.Create(ctx, title); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Update is department update
func (s *DepartmentService) Update(ctx context.Context, id int64, title string) error {
	const fn = "service.Department.Update"

	if err := s.DepartmentRepository.Update(ctx, id, title); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Delete is department delete
func (s *DepartmentService) Delete(ctx context.Context, id int64) error {
	const fn = "service.Department.Delete"

	if err := s.DepartmentRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Restore is department restore
func (s *DepartmentService) Restore(ctx context.Context, id int64) error {
	const fn = "service.Department.Restore"

	if err := s.DepartmentRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// GetAll is to get all departments
func (s *DepartmentService) GetAll(ctx context.Context, deleted bool) ([]*model.Department, error) {
	const fn = "service.Department.GetAll"

	res, err := s.DepartmentRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return res, nil
}

// GetById is to get department by id
func (s *DepartmentService) GetById(ctx context.Context, id int64) (*model.Department, error) {
	const fn = "service.Department.GetById"

	res, err := s.DepartmentRepository.GetById(ctx, &model.Department{ID: id})
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return res, nil
}

// GetAllButOne is to get all departments but one
func (s *DepartmentService) GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error) {
	const fn = "service.Department.GetAllButOne"

	employee := &model.Employee{
		ID: employeeId,
	}

	res, err := s.EmployeeRepository.GetById(ctx, employee)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	var departments []*model.Department
	if res.Role == "ADMIN" {
		res, err := s.DepartmentRepository.GetAllButOneForAdmin(ctx, id)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		departments = append(departments, res...)
	} else {
		res, err := s.DepartmentRepository.GetAllButOne(ctx, id, employeeId)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		departments = append(departments, res...)
	}

	return departments, nil
}
