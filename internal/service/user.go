package service

import (
	"context"
	"errors"
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
		return logger.Error(logger.MsgFailedToGenerateHash, err)
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

	logger.Info(fmt.Sprintf("user with id %d created", id))
	return nil
}

func (s *UserService) Read(ctx context.Context, id int64) (*model.User, error) {
	read, err := s.userRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("user with id %d read", id))
	return read, nil
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d updated", user.ID))
	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d deleted", id))
	return nil
}

func (s *UserService) List(ctx context.Context) ([]*model.User, error) {
	list, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%d user listed", len(list)))
	return list, nil
}

func (s *UserService) SetPassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	oldPasswordHash, err := s.userRepository.GetPasswordHash(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oldPasswordHash), []byte(oldPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return logger.Error(logger.MsgFailedToValidate, logger.ErrWrongPassword)
		}
		return err
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return logger.Error(logger.MsgFailedToGenerateHash, err)
	}

	err = s.userRepository.SetPasswordHash(ctx, id, string(newPasswordHash))
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d changed password", id))
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, id int64) error {
	newPassword := generate.RandString(10)
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return logger.Error(logger.MsgFailedToGenerateHash, err)
	}

	err = s.userRepository.SetPasswordHash(ctx, id, string(newPasswordHash))
	if err != nil {
		return err
	}

	user, err := s.userRepository.Read(ctx, id)
	if err != nil {
		return err
	}

	sendTo := &email.SendTo{
		Name:     user.Employee.FirstName,
		Email:    user.Email,
		Username: user.Username,
		Password: newPassword,
	}

	if err := email.Send([]*email.SendTo{sendTo}); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d reset password", id))
	return nil
}

func (s *UserService) SetRole(ctx context.Context, id int64, role role.Role) error {
	if err := s.userRepository.SetRole(ctx, id, role); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d changed role to %s", id, role))
	return nil
}

func (s *UserService) SetEnabled(ctx context.Context, id int64, enabled bool) error {
	if err := s.userRepository.SetEnabled(ctx, id, enabled); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d set enabled to %t", id, enabled))
	return nil
}

func (s *UserService) SetEmployee(ctx context.Context, id, employeeID int64) error {
	if err := s.userRepository.SetEmployee(ctx, id, employeeID); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d set employee id %d", id, employeeID))
	return nil
}
