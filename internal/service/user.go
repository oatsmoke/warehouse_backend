package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/lib/email"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository     repository.User
	employeeRepository repository.Employee
}

func NewUserService(userRepository repository.User, employeeRepository repository.Employee) *UserService {
	return &UserService{
		userRepository:     userRepository,
		employeeRepository: employeeRepository,
	}
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	employee := new(model.Employee)
	var employeeID int64
	if user != nil && user.Employee.ID != 0 {
		read, err := s.employeeRepository.Read(ctx, user.Employee.ID)
		if err != nil {
			return err
		}

		employeeID = user.Employee.ID
		employee = read
	}

	var username string
	if user != nil && user.Username != "" {
		username = user.Username
	} else if employee != nil && employee.Phone != "" {
		username = employee.Phone
	} else {
		username = fmt.Sprintf("user-%s", generate.RandString(10))
	}

	password := generate.RandString(10)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var userEmail string
	if user != nil {
		userEmail = user.Email
	}

	var userRole role.Role
	if user != nil && user.Role != "" {
		userRole = user.Role
	} else {
		userRole = role.EmployeeRole
	}

	u := &model.User{
		Username:     username,
		PasswordHash: string(passwordHash),
		Email:        userEmail,
		Role:         userRole,
		Employee: &model.Employee{
			ID: employeeID,
		},
	}

	id, err := s.userRepository.Create(ctx, u)
	if err != nil {
		return err
	}

	var employeeFirstName string
	if employee != nil {
		employeeFirstName = employee.FirstName
	}

	sendTo := &email.SendTo{
		Name:     employeeFirstName,
		Email:    u.Email,
		Username: username,
		Password: password,
	}

	if err := email.Send([]*email.SendTo{sendTo}); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("user with id %d created", id))
	return nil
}

func (s *UserService) Read(ctx context.Context, id int64) (*model.User, error) {
	read, err := s.userRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("user with id %d read", id))
	return read, nil
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("user with id %d updated", user.ID))
	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("user with id %d deleted", id))
	return nil
}

func (s *UserService) List(ctx context.Context) ([]*model.User, error) {
	list, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("%d user listed", len(list)))
	return list, nil
}
