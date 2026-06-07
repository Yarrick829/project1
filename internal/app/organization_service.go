package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type organizationService struct {
	orgRepo  database.OrganizationRepository
	roomRepo database.RoomRepository
}

type OrganizationService interface {
	Save(o domain.Organization) (domain.Organization, error)
	FindList(uId uint64) ([]domain.Organization, error)
	Find(id uint64) (interface{}, error)
	Update(o domain.Organization) (domain.Organization, error)
	Delete(id uint64) error
}

func NewOrganizationService(or database.OrganizationRepository, rr database.RoomRepository) OrganizationService {
	return organizationService{
		orgRepo:  or,
		roomRepo: rr,
	}
}

func (s organizationService) Save(o domain.Organization) (domain.Organization, error) {
	org, err := s.orgRepo.Save(o)
	if err != nil {
		log.Printf("organizationService.Save(s.orgRepo.Save): %s", err)
		return domain.Organization{}, err
	}

	return org, nil
}

func (s organizationService) FindList(uId uint64) ([]domain.Organization, error) {
	orgs, err := s.orgRepo.FindList(uId)
	if err != nil {
		log.Printf("organizationService.FindList(s.orgRepo.FindList): %s", err)
		return nil, err
	}

	return orgs, nil

}

func (s organizationService) Find(id uint64) (interface{}, error) {
	org, err := s.orgRepo.Find(id)
	if err != nil {
		return nil, err
	}

	org.Rooms, err = s.roomRepo.FindByOrgId(org.Id)
	if err != nil {
		return nil, err
	}

	log.Printf("organization %d has %d rooms", org.Id, len(org.Rooms))

	return org, nil
}

func (s organizationService) Update(o domain.Organization) (domain.Organization, error) {
	org, err := s.orgRepo.Update(o)
	if err != nil {
		log.Printf("organizationService.Update(s.orgRepo.Update): %s", err)
		return domain.Organization{}, err
	}

	return org, nil
}

func (s organizationService) Delete(id uint64) error {
	err := s.orgRepo.Delete(id)
	if err != nil {
		log.Printf("organizationServise.Delete(s.orgRepo.Delete): %s", err)
		return err
	}

	return nil
}
