package service

import (
	"context"
	"errors"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type ProfileService struct {
	repositoryProfile   repository.Profile
	repositoryEquipment repository.Equipment
}

func NewProfileService(repositoryProfile repository.Profile,
	repositoryEquipment repository.Equipment) *ProfileService {
	return &ProfileService{repositoryProfile: repositoryProfile,
		repositoryEquipment: repositoryEquipment,
	}
}

func (s *ProfileService) Create(ctx context.Context, title string, category int64) error {
	if _, err := s.repositoryProfile.FindByTitle(ctx, title); err == nil {
		return errors.New("title already exists")
	}

	return s.repositoryProfile.Create(ctx, title, category)
}

func (s *ProfileService) GetById(ctx context.Context, id int64) (*model.Profile, error) {
	return s.repositoryProfile.GetById(ctx, id)
}
func (s *ProfileService) GetAll(ctx context.Context) ([]*model.Profile, error) {
	return s.repositoryProfile.GetAll(ctx)
}

func (s *ProfileService) Update(ctx context.Context, id int64, title string, category int64) error {
	findId, err := s.repositoryProfile.FindByTitle(ctx, title)
	if findId != id && err == nil {
		return errors.New("title already exists")
	}

	return s.repositoryProfile.Update(ctx, id, title, category)
}

func (s *ProfileService) Delete(ctx context.Context, id int64) error {
	equipments, err := s.repositoryEquipment.GetByProfile(ctx, id)
	if err != nil {
		return err
	}

	if len(equipments) > 0 {
		return errors.New("used in equipment")
	}

	return s.repositoryProfile.Delete(ctx, id)
}
