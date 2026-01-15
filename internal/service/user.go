package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/email"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
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
	name := user.Username
	if user.Employee.ID != 0 {
		read, err := s.employeeRepository.Read(ctx, user.Employee.ID)
		if err != nil {
			return err
		}
		name = read.FirstName
	}

	password := generate.RandString(10)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return logger.Error(logger.MsgFailedToGenerateHash, err)
	}
	user.PasswordHash = string(passwordHash)

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	sendTo := &email.SendTo{
		Name:     name,
		Email:    user.Email,
		Username: user.Username,
		Password: password,
	}

	go email.Send([]*email.SendTo{sendTo})

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

func (s *UserService) List(ctx context.Context) ([]*dto.UserResponse, error) {
	list, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	userRes := make([]*dto.UserResponse, 0, len(list))
	for _, user := range list {
		u := &dto.UserResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Role:         user.Role.String(),
			Enabled:      user.Enabled,
			EmployeeName: shortEmployeeName(user.Employee.LastName, user.Employee.FirstName, user.Employee.MiddleName),
		}

		if user.LastLoginAt != nil {
			u.LastLoginAt = user.LastLoginAt.Format("02.01.2006 15:04:05")
		}

		userRes = append(userRes, u)
	}

	logger.Info(fmt.Sprintf("%d user listed", len(list)))
	return userRes, nil
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

	go email.Send([]*email.SendTo{sendTo})

	logger.Info(fmt.Sprintf("user with id %d reset password", id))
	return nil
}

func (s *UserService) SetEnabled(ctx context.Context, id int64, enabled bool) error {
	if err := s.userRepository.SetEnabled(ctx, id, enabled); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("user with id %d set enabled to %t", id, enabled))
	return nil
}
