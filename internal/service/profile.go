package service

import (
	"context"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type ProfileService struct {
	ProfileRepository repository.Profile
}

func NewProfileService(profileRepository repository.Profile) *ProfileService {
	return &ProfileService{
		ProfileRepository: profileRepository,
	}
}

// Create is profile create
func (s *ProfileService) Create(ctx context.Context, title string, category int64) error {
	const fn = "service.Profile.Create"

	if err := s.ProfileRepository.Create(ctx, title, category); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Update is profile update
func (s *ProfileService) Update(ctx context.Context, id int64, title string, category int64) error {
	const fn = "service.Profile.Update"

	if err := s.ProfileRepository.Update(ctx, id, title, category); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Delete is profile delete
func (s *ProfileService) Delete(ctx context.Context, id int64) error {
	const fn = "service.Profile.Delete"

	if err := s.ProfileRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Restore is profile restore
func (s *ProfileService) Restore(ctx context.Context, id int64) error {
	const fn = "service.Profile.Restore"

	if err := s.ProfileRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// GetAll is to get all profiles
func (s *ProfileService) GetAll(ctx context.Context, deleted bool) ([]*model.Profile, error) {
	const fn = "service.Profile.GetAll"

	profiles, err := s.ProfileRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return profiles, nil
}

// GetById is to get profile by id
func (s *ProfileService) GetById(ctx context.Context, id int64) (*model.Profile, error) {
	const fn = "service.Profile.GetById"

	category, err := s.ProfileRepository.GetById(ctx, &model.Profile{ID: id})
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return category, nil
}
