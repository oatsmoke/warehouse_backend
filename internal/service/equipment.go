package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type EquipmentService struct {
	equipmentRepository repository.Equipment
	locationRepository  repository.Location
}

func NewEquipmentService(equipmentRepository repository.Equipment, locationRepository repository.Location) *EquipmentService {
	return &EquipmentService{
		equipmentRepository: equipmentRepository,
		locationRepository:  locationRepository,
	}
}

func (s *EquipmentService) Create(ctx context.Context, userId int64, req *dto.CreateEquipmentRequest) error {
	d, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return err
	}
	date := pgtype.Timestamptz{
		Time:  d,
		Valid: true,
	}

	l := &queries.AddToStorageParams{
		UserID:   userId,
		MoveAt:   date,
		MoveCode: "AddToStorage",
	}

	for _, sn := range req.SerialNumbers {

		e := &queries.CreateEquipmentParams{
			SerialNumber: sn,
			ProfileID:    req.ProfileID,
			CompanyID:    req.CompanyID,
		}

		id, err := s.equipmentRepository.Create(ctx, e, l)
		if err != nil {
			logger.Warn(fmt.Sprintf("equipment [%s] create error: %v", sn, err))
			continue
		}

		if req.ParamID != 0 {
			move := &queries.MoveToLocationParams{
				EquipmentID:    id,
				UserID:         userId,
				MoveAt:         date,
				MoveCode:       fmt.Sprintf("StorageTo%s", req.Param),
				ToDepartmentID: toPGTypeInt8(req.ParamID),
			}

			if err := s.locationRepository.Move(ctx, move); err != nil {
				logger.Warn(fmt.Sprintf("equipment [%s] move error: %v", sn, err))
			}
		}

		logger.Info(fmt.Sprintf("equipment [%s] created", sn))
	}

	return nil
}

func (s *EquipmentService) Read(ctx context.Context, id int64) (*model.Equipment, error) {
	read, err := s.equipmentRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("equipment with id %d read", id))
	return read, nil
}

func (s *EquipmentService) Update(ctx context.Context, equipment *model.Equipment) error {
	if err := s.equipmentRepository.Update(ctx, equipment); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("equipment with id %d updated", equipment.ID))
	return nil
}

func (s *EquipmentService) Delete(ctx context.Context, id int64) error {
	if err := s.equipmentRepository.Delete(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("equipment with id %d deleted", id))
	return nil
}

func (s *EquipmentService) Restore(ctx context.Context, id int64) error {
	if err := s.equipmentRepository.Restore(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("equipment with id %d restored", id))
	return nil
}

func (s *EquipmentService) List(ctx context.Context, qp *dto.QueryParams) (*dto.ListResponse[[]*model.Equipment], error) {
	list, total, err := s.equipmentRepository.List(ctx, qp)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%d equipment listed", len(list)))
	return &dto.ListResponse[[]*model.Equipment]{
		List:  list,
		Total: total,
	}, nil
}
