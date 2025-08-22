package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type EquipmentService struct {
	EquipmentRepository repository.Equipment
}

func NewEquipmentService(equipmentRepository repository.Equipment) *EquipmentService {
	return &EquipmentService{
		EquipmentRepository: equipmentRepository,
	}
}

// Create is equipment create
func (s *EquipmentService) Create(ctx context.Context, serialNumber string, profileId int64) (int64, error) {
	id, err := s.EquipmentRepository.Create(ctx, strings.ToUpper(serialNumber), profileId)
	if err != nil {
		return 0, logger.Err(err, "")
	}

	return id, nil
}

// Update is equipment update
func (s *EquipmentService) Update(ctx context.Context, id int64, serialNumber string, profileId int64) error {
	if err := s.EquipmentRepository.Update(ctx, id, strings.ToUpper(serialNumber), profileId); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Delete is equipment delete
func (s *EquipmentService) Delete(ctx context.Context, id int64) error {
	if err := s.EquipmentRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Restore is equipment restore
func (s *EquipmentService) Restore(ctx context.Context, id int64) error {
	if err := s.EquipmentRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// GetAll is to get all equipment
func (s *EquipmentService) GetAll(ctx context.Context) ([]*model.Equipment, error) {
	res, err := s.EquipmentRepository.GetAll(ctx)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetByIds is equipment get by id
func (s *EquipmentService) GetByIds(ctx context.Context, ids []int64) ([]*model.Equipment, error) {
	res, err := s.EquipmentRepository.GetByIds(ctx, ids)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// FindBySerialNumber is equipment find by serial number
func (s *EquipmentService) FindBySerialNumber(ctx context.Context, value string) ([]*model.Equipment, error) {
	res, err := s.EquipmentRepository.FindBySerialNumber(ctx, fmt.Sprintf("%%%s%%", value))
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}
