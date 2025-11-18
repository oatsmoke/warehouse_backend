package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type UserRepository struct {
	queries queries.Querier
}

func NewUserRepository(queries queries.Querier) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	var employeeID pgtype.Int8
	if user.Employee != nil {
		employeeID = pgtype.Int8{
			Int64: user.Employee.ID,
			Valid: user.Employee.ID != 0,
		}
	}

	req, err := r.queries.CreateUser(ctx, &queries.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Email:        user.Email,
		Role:         string(user.Role),
		EmployeeID:   employeeID,
	})
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *UserRepository) Read(ctx context.Context, id int64) (*model.User, error) {
	req, err := r.queries.ReadUser(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	user := &model.User{
		ID:          req.ID,
		Username:    req.Username,
		Email:       req.Email,
		Role:        role.Role(req.Role),
		Enabled:     req.Enabled,
		LastLoginAt: validTime(req.LastLoginAt),
		Employee: &model.Employee{
			ID:         validInt64(req.EmployeeID),
			LastName:   validString(req.EmployeeLastName),
			FirstName:  validString(req.EmployeeFirstName),
			MiddleName: validString(req.EmployeeMiddleName),
			Phone:      validString(req.EmployeePhone),
			Department: &model.Department{
				ID:    validInt64(req.DepartmentID),
				Title: validString(req.DepartmentTitle),
			},
		},
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	ct, err := r.queries.UpdateUser(ctx, &queries.UpdateUserParams{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]*model.User, error) {
	req, err := r.queries.ListUser(ctx)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToSelect, err)
	}

	if len(req) < 1 {
		return nil, nil
	}

	list := make([]*model.User, len(req))
	for i, item := range req {
		user := &model.User{
			ID:          item.ID,
			Username:    item.Username,
			Email:       item.Email,
			Role:        role.Role(item.Role),
			Enabled:     item.Enabled,
			LastLoginAt: validTime(item.LastLoginAt),
			Employee: &model.Employee{
				ID:         validInt64(item.EmployeeID),
				LastName:   validString(item.EmployeeLastName),
				FirstName:  validString(item.EmployeeFirstName),
				MiddleName: validString(item.EmployeeMiddleName),
				Phone:      validString(item.EmployeePhone),
				Department: &model.Department{
					ID:    validInt64(item.DepartmentID),
					Title: validString(item.DepartmentTitle),
				},
			},
		}
		list[i] = user
	}

	return list, nil
}

func (r *UserRepository) GetPasswordHash(ctx context.Context, id int64) (string, error) {
	passwordHash, err := r.queries.GetPasswordHashUser(ctx, id)
	if err != nil {
		return "", logger.Error(logger.MsgFailedToScan, err)
	}

	return passwordHash, nil
}

func (r *UserRepository) SetPasswordHash(ctx context.Context, id int64, passwordHash string) error {
	ct, err := r.queries.SetPasswordHashUser(ctx, &queries.SetPasswordHashUserParams{
		ID:           id,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *UserRepository) SetRole(ctx context.Context, id int64, role role.Role) error {
	ct, err := r.queries.SetRoleUser(ctx, &queries.SetRoleUserParams{
		ID:   id,
		Role: string(role),
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *UserRepository) SetEnabled(ctx context.Context, id int64, enabled bool) error {
	ct, err := r.queries.SetEnabledUser(ctx, &queries.SetEnabledUserParams{
		ID:      id,
		Enabled: enabled,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *UserRepository) SetLastLoginAt(ctx context.Context, id int64) error {
	ct, err := r.queries.SetLastLoginAtUser(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *UserRepository) SetEmployee(ctx context.Context, id, employeeID int64) error {
	var e pgtype.Int8
	if employeeID != 0 {
		e = pgtype.Int8{
			Int64: employeeID,
			Valid: true,
		}
	}

	ct, err := r.queries.SetEmployeeUser(ctx, &queries.SetEmployeeUserParams{
		ID:         id,
		EmployeeID: e,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	req, err := r.queries.GetByUsernameUser(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	user := &model.User{
		ID:           req.ID,
		Username:     req.Username,
		PasswordHash: req.PasswordHash,
		Email:        req.Email,
		Role:         role.Role(req.Role),
		Enabled:      req.Enabled,
		LastLoginAt:  validTime(req.LastLoginAt),
	}

	return user, nil
}
