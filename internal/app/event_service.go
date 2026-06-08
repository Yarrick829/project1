package app

import (
	"errors"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type eventService struct {
	eventRepo database.EventRepository
}

type EventService interface {
	Save(e domain.Event) (domain.Event, error)
	Update(e domain.Event) (domain.Event, error)
	Delete(id uint64) error
	Find(eId uint64) (interface{}, error)
	FindList(p domain.Pagination, ef database.EventFilters) (domain.Events, error)
}

func NewEventService(er database.EventRepository) EventService {
	return eventService{
		eventRepo: er,
	}
}

func (es eventService) Save(e domain.Event) (domain.Event, error) {

	err := es.validateEvent(e)
	if err != nil {
		log.Printf("eventService.Save(validateEvent): %s", err)
		return domain.Event{}, err
	}

	e, err = es.eventRepo.Save(e)
	if err != nil {
		log.Printf("eventService.Save(es.eventRepo.Save): %s", err)
		return domain.Event{}, err
	}

	return e, nil
}

func (es eventService) Update(e domain.Event) (domain.Event, error) {

	err := es.validateEvent(e)
	if err != nil {
		log.Printf("eventService.Update(validateEvent): %s", err)
		return domain.Event{}, err
	}

	e, err = es.eventRepo.Update(e)
	if err != nil {
		log.Printf("eventService.Update(es.eventRepo.Update): %s", err)
		return domain.Event{}, err
	}

	return e, nil
}

func (es eventService) Delete(id uint64) error {

	err := es.eventRepo.Delete(id)
	if err != nil {
		log.Printf("eventService.Delete(es.eventRepo.Delete): %s", err)
		return err
	}

	return nil
}

func (es eventService) Find(eId uint64) (interface{}, error) {

	e, err := es.eventRepo.FindByEventId(eId)
	if err != nil {
		log.Printf("eventService.Find(es.eventRepo.FindByEventId): %s", err)
		return nil, err
	}

	return e, nil
}

func (es eventService) validateEvent(e domain.Event) error {

	if e.DeviceId == 0 {
		return errors.New("device id is required")
	}

	if e.RoomId == 0 {
		return errors.New("room id is required")
	}

	switch e.Action {
	case domain.EventOn:
		return nil

	case domain.EventOff:
		return nil

	default:
		return errors.New("invalid event action")
	}
}

func (es eventService) FindList(p domain.Pagination, ef database.EventFilters) (domain.Events, error) {

	events, err := es.eventRepo.FindList(p, ef)

	if err != nil {
		log.Printf(
			"eventService.FindList(es.eventRepo.FindList): %s",
			err,
		)

		return domain.Events{}, err
	}

	return events, nil
}
