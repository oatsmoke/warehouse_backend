package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type DepartmentService struct {
	departmentRepository repository.Department
}

func NewDepartmentService(departmentRepository repository.Department) *DepartmentService {
	return &DepartmentService{
		departmentRepository: departmentRepository,
	}
}

func (s *DepartmentService) Create(ctx context.Context, department *model.Department) error {
	id, err := s.departmentRepository.Create(ctx, department)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("department with id %d created", id))
	return nil
}

func (s *DepartmentService) Read(ctx context.Context, id int64) (*model.Department, error) {
	read, err := s.departmentRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("department with id %d read", id))
	return read, nil
}

func (s *DepartmentService) Update(ctx context.Context, department *model.Department) error {
	if err := s.departmentRepository.Update(ctx, department); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("department with id %d updated", department.ID))
	return nil
}

func (s *DepartmentService) Delete(ctx context.Context, id int64) error {
	if err := s.departmentRepository.Delete(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("department with id %d deleted", id))
	return nil
}

func (s *DepartmentService) Restore(ctx context.Context, id int64) error {
	if err := s.departmentRepository.Restore(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("department with id %d restored", id))
	return nil
}

func (s *DepartmentService) List(ctx context.Context, qp *dto.QueryParams) (*dto.ListResponse[[]*model.Department], error) {
	list, total, err := s.departmentRepository.List(ctx, qp)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%d department listed", len(list)))
	return &dto.ListResponse[[]*model.Department]{
		List:  list,
		Total: total,
	}, nil
}

//func (s *DepartmentService) GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error) {
//	employee := &model.Employee{
//		ID: employeeId,
//	}
//
//	res, err := s.EmployeeRepository.GetById(ctx, employee)
//	if err != nil {
//		return nil, err
//	}
//
//	var departments []*model.Department
//	if res.Role == "ADMIN" {
//		res, err := s.DepartmentRepository.GetAllButOneForAdmin(ctx, id)
//		if err != nil {
//			return nil, err
//		}
//		departments = append(departments, res...)
//	} else {
//		res, err := s.DepartmentRepository.GetAllButOne(ctx, id, employeeId)
//		if err != nil {
//			return nil, err
//		}
//		departments = append(departments, res...)
//	}
//
//	return departments, nil
//}
