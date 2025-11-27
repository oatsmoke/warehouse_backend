package repository

import (
	"context"

	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type ContractRepository struct {
	queries queries.Querier
}

func NewContractRepository(queries queries.Querier) *ContractRepository {
	return &ContractRepository{
		queries: queries,
	}
}

func (r *ContractRepository) Create(ctx context.Context, contract *model.Contract) (int64, error) {
	req, err := r.queries.CreateContract(ctx, &queries.CreateContractParams{
		Number:  contract.Number,
		Address: contract.Address,
	})
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *ContractRepository) Read(ctx context.Context, id int64) (*model.Contract, error) {
	req, err := r.queries.ReadContract(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	contract := &model.Contract{
		ID:        req.ID,
		Number:    req.Number,
		Address:   req.Address,
		DeletedAt: validTime(req.DeletedAt),
	}

	return contract, nil
}

func (r *ContractRepository) Update(ctx context.Context, contract *model.Contract) error {
	ct, err := r.queries.UpdateContract(ctx, &queries.UpdateContractParams{
		ID:      contract.ID,
		Number:  contract.Number,
		Address: contract.Address,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *ContractRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteContract(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *ContractRepository) Restore(ctx context.Context, id int64) error {
	ct, err := r.queries.RestoreContract(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *ContractRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Contract, int64, error) {
	req, err := r.queries.ListContract(ctx, &queries.ListContractParams{
		WithDeleted:      qp.WithDeleted,
		Search:           qp.Search,
		Ids:              qp.Ids,
		SortColumn:       qp.SortColumn,
		SortOrder:        qp.SortOrder,
		PaginationLimit:  qp.PaginationLimit,
		PaginationOffset: qp.PaginationOffset,
	})
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}

	if len(req) < 1 {
		return []*model.Contract{}, 0, nil
	}

	list := make([]*model.Contract, len(req))
	for i, item := range req {
		contract := &model.Contract{
			ID:        item.ID,
			Number:    item.Number,
			Address:   item.Address,
			DeletedAt: validTime(item.DeletedAt),
		}
		list[i] = contract
	}

	return list, req[0].Total, nil
}
