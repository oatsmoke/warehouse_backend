package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type EmployeeService struct {
	employeeRepository repository.Employee
}

func NewEmployeeService(employeeRepository repository.Employee) *EmployeeService {
	return &EmployeeService{
		employeeRepository: employeeRepository,
	}
}

func (s *EmployeeService) Create(ctx context.Context, employee *model.Employee) error {
	id, err := s.employeeRepository.Create(ctx, employee)
	if err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("employee with id %d created", id))
	return nil
}

func (s *EmployeeService) Read(ctx context.Context, id int64) (*model.Employee, error) {
	read, err := s.employeeRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("employee with id %d read", id))
	return read, nil
}

func (s *EmployeeService) Update(ctx context.Context, employee *model.Employee) error {
	if err := s.employeeRepository.Update(ctx, employee); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("employee with id %d updated", employee.ID))
	return nil
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	if err := s.employeeRepository.Delete(ctx, id); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("employee with id %d deleted", id))
	return nil
}

func (s *EmployeeService) Restore(ctx context.Context, id int64) error {
	if err := s.employeeRepository.Restore(ctx, id); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("employee with id %d restored", id))
	return nil
}

func (s *EmployeeService) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Employee, error) {
	list, err := s.employeeRepository.List(ctx, qp)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("%d employee listed", len(list)))
	return list, nil
}

// GetAllShort is to get all employees short
//func (s *EmployeeService) GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error) {
//	res, err := s.employeeRepository.GetAllShort(ctx, deleted)
//	if err != nil {
//		return nil, err
//	}
//
//	return res, nil
//}
//
//// GetAllButOne is to get all employees but one
//func (s *EmployeeService) GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error) {
//	res, err := s.employeeRepository.GetAllButOne(ctx, id, deleted)
//	if err != nil {
//		return nil, err
//	}
//
//	return res, nil
//}
//
//// GetById is employee get by id
//
//// GetFree is employee get free
//func (s *EmployeeService) GetFree(ctx context.Context) ([]*model.Employee, error) {
//	res, err := s.employeeRepository.GetFree(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	return res, nil
//}
//
//// GetByDepartment is employee get by department
//func (s *EmployeeService) GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error) {
//	res, err := s.employeeRepository.GetByDepartment(ctx, ids, departmentId)
//	if err != nil {
//		return nil, err
//	}
//
//	return res, nil
//}
//
//// AddToDepartment is employee add to department
//func (s *EmployeeService) AddToDepartment(ctx context.Context, id, department int64) error {
//	if err := s.employeeRepository.AddToDepartment(ctx, id, department); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// RemoveFromDepartment is employee remove from department
//func (s *EmployeeService) RemoveFromDepartment(ctx context.Context, idDepartment, idEmployee int64) error {
//	if err := s.employeeRepository.RemoveFromDepartment(ctx, idEmployee); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Activate is employee activate
//func (s *EmployeeService) Activate(ctx context.Context, id int64) error {
//	employee := &model.Employee{
//		ID: id,
//	}
//
//	res, err := s.employeeRepository.GetById(ctx, employee)
//	if err != nil {
//		return err
//	}
//
//	str := generate.RandString(10)
//	passwordHash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//
//	if err := s.employeeRepository.Activate(ctx, id, string(passwordHash)); err != nil {
//		return err
//	}
//
//	sendTo := &email.SendTo{
//		Name:     res.Name,
//		Email:    res.Email,
//		Phone:    res.Phone,
//		Password: string(passwordHash),
//	}
//
//	if err := email.Send([]*email.SendTo{sendTo}); err != nil {
//		return err
//	}
//
//	//if err := sendMail(res.Email, res.Phone, string(passwordHash)); err != nil {
//	//	return err
//	//}
//
//	return nil
//}
//
//// Deactivate is employee deactivate
//func (s *EmployeeService) Deactivate(ctx context.Context, id int64) error {
//	if err := s.employeeRepository.Deactivate(ctx, id); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// ResetPassword is employee reset password
//func (s *EmployeeService) ResetPassword(ctx context.Context, id int64) error {
//	employee := &model.Employee{
//		ID: id,
//	}
//
//	res, err := s.employeeRepository.GetById(ctx, employee)
//	if err != nil {
//		return err
//	}
//
//	str := generate.RandString(10)
//	passwordHash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//
//	if err := s.employeeRepository.ResetPassword(ctx, id, string(passwordHash)); err != nil {
//		return err
//	}
//
//	sendTo := &email.SendTo{
//		Name:     res.Name,
//		Email:    res.Email,
//		Phone:    res.Phone,
//		Password: string(passwordHash),
//	}
//
//	if err := email.Send([]*email.SendTo{sendTo}); err != nil {
//		return err
//	}
//
//	//if err := sendMail(res.Email, res.Phone, string(passwordHash)); err != nil {
//	//	return err
//	//}
//
//	return nil
//}
//
//// ChangeRole is employee change role
//func (s *EmployeeService) ChangeRole(ctx context.Context, id int64, role string) error {
//	if err := s.employeeRepository.ChangeRole(ctx, id, role); err != nil {
//		return err
//	}
//
//	return nil
//}
