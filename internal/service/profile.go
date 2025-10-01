package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type ProfileService struct {
	profileRepository repository.Profile
}

func NewProfileService(profileRepository repository.Profile) *ProfileService {
	return &ProfileService{
		profileRepository: profileRepository,
	}
}

func (s *ProfileService) Create(ctx context.Context, profile *model.Profile) error {
	id, err := s.profileRepository.Create(ctx, profile)
	if err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("profile with id %d created", id))
	return nil
}

func (s *ProfileService) Read(ctx context.Context, id int64) (*model.Profile, error) {
	read, err := s.profileRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("profile with id %d read", id))
	return read, nil
}

func (s *ProfileService) Update(ctx context.Context, profile *model.Profile) error {
	if err := s.profileRepository.Update(ctx, profile); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("profile with id %d updated", profile.ID))
	return nil
}

func (s *ProfileService) Delete(ctx context.Context, id int64) error {
	if err := s.profileRepository.Delete(ctx, id); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("profile with id %d deleted", id))
	return nil
}

func (s *ProfileService) Restore(ctx context.Context, id int64) error {
	if err := s.profileRepository.Restore(ctx, id); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("profile with id %d restored", id))
	return nil
}

func (s *ProfileService) List(ctx context.Context, withDeleted bool) ([]*model.Profile, error) {
	list, err := s.profileRepository.List(ctx, withDeleted)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("%d profile listed", len(list)))
	return list, nil
}
