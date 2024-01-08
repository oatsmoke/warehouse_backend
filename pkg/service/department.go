package service

import (
	"errors"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
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

func (s *DepartmentService) Create(title string) error {
	if _, err := s.repositoryDepartment.FindByTitle(title); err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryDepartment.Create(title)
}

func (s *DepartmentService) GetById(id int) (model.Department, error) {
	if id == 0 {
		return model.Department{}, nil
	}
	return s.repositoryDepartment.GetById(id)
}

func (s *DepartmentService) GetAll() ([]model.Department, error) {
	return s.repositoryDepartment.GetAll()
}

func (s *DepartmentService) GetAllButOne(id, employeeId int) ([]model.Department, error) {
	employee, err := s.repositoryEmployee.GetById(employeeId)
	if err != nil {
		return []model.Department{}, err
	}
	if employee.Role == "ADMIN" {
		return s.repositoryDepartment.GetAllButOneForAdmin(id)
	}
	return s.repositoryDepartment.GetAllButOne(id, employeeId)
}

func (s *DepartmentService) Update(id int, title string) error {
	if _, err := s.repositoryDepartment.FindByTitle(title); err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryDepartment.Update(id, title)
}

func (s *DepartmentService) Delete(id int) error {
	equipments, err := s.repositoryEquipment.GetByLocationDepartment(id)
	if err != nil {
		return err
	}
	if len(equipments) > 0 {
		return errors.New("have equipment")
	}
	employees, err := s.repositoryEmployee.GetByDepartment([]int{}, id)
	if err != nil {
		return err
	}
	if len(employees) > 0 {
		return errors.New("have employees")
	}
	return s.repositoryDepartment.Delete(id)
}
