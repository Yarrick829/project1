package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type roomService struct {
	roomRepo   database.RoomRepository
	deviceRepo database.DeviceRepository
}

type RoomService interface {
	Save(o domain.Room) (domain.Room, error)
	Update(o domain.Room) (domain.Room, error)
	Delete(id uint64) error
	Find(rId uint64) (interface{}, error)
	FindList(orgId uint) ([]domain.Room, error)
}

func NewRoomService(
	rr database.RoomRepository,
	dr database.DeviceRepository,
) RoomService {
	return roomService{
		roomRepo:   rr,
		deviceRepo: dr,
	}
}

func (rs roomService) Save(rm domain.Room) (domain.Room, error) {
	rm, err := rs.roomRepo.Save(rm)
	if err != nil {
		log.Printf("roomService.Save(rs.roomRepo.Save): %s", err)
		return domain.Room{}, err
	}

	return rm, nil
}

func (rs roomService) Update(rm domain.Room) (domain.Room, error) {
	rm, err := rs.roomRepo.Update(rm)
	if err != nil {
		log.Printf("roomService.Update(rs.roomRepo.Update): %s", err)
		return domain.Room{}, err
	}

	return rm, nil
}

func (rs roomService) Delete(id uint64) error {
	err := rs.roomRepo.Delete(id)
	if err != nil {
		log.Printf("roomService.Delete(rs.roomRepo.Delete): %s", err)
		return err
	}

	return nil
}

func (rs roomService) Find(rId uint64) (interface{}, error) {
	rm, err := rs.roomRepo.FindByRoomId(rId)
	if err != nil {
		return nil, err
	}

	rm.Devices, err = rs.deviceRepo.FindByRoomId(rm.Id)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

func (rs roomService) FindList(orgId uint) ([]domain.Room, error) {
	rooms, err := rs.roomRepo.FindByOrgId(uint64(orgId))
	if err != nil {
		log.Printf("roomService.FindList(rs.roomRepo.FindByOrgId): %s", err)
		return nil, err
	}

	return rooms, nil
}
