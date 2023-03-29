package service

import (
	"errors"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
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

func (s *ProfileService) Create(title string, category int) error {
	if _, err := s.repositoryProfile.FindByTitle(title); err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryProfile.Create(title, category)
}

func (s *ProfileService) GetById(id int) (model.Profile, error) {
	return s.repositoryProfile.GetById(id)
}
func (s *ProfileService) GetAll() ([]model.Profile, error) {
	return s.repositoryProfile.GetAll()
}

func (s *ProfileService) Update(id int, title string, category int) error {
	findId, err := s.repositoryProfile.FindByTitle(title)
	if findId != id && err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryProfile.Update(id, title, category)
}

func (s *ProfileService) Delete(id int) error {
	equipments, err := s.repositoryEquipment.GetByProfile(id)
	if err != nil {
		return err
	}
	if len(equipments) > 0 {
		return errors.New("used in equipment")
	}
	return s.repositoryProfile.Delete(id)
}
