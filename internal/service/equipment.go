package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type EquipmentService struct {
	equipmentRepository repository.Equipment
}

func NewEquipmentService(equipmentRepository repository.Equipment) *EquipmentService {
	return &EquipmentService{
		equipmentRepository: equipmentRepository,
	}
}

func (s *EquipmentService) Create(ctx context.Context, equipment *model.Equipment) error {
	id, err := s.equipmentRepository.Create(ctx, equipment)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("equipment with id %d created", id))
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

func (s *EquipmentService) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Equipment, error) {
	list, err := s.equipmentRepository.List(ctx, qp)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%d equipment listed", len(list)))
	return list, nil
}
