package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type EventDto struct {
	Id          uint64             `json:"id"`
	DeviceId    uint64             `json:"device_id"`
	RoomId      uint64             `json:"room_id"`
	Action      domain.EventAction `json:"action"`
	CreatedDate time.Time          `json:"created_date"`
}

func (ed EventDto) DomainToDto(e domain.Event) EventDto {
	return EventDto{
		Id:          e.Id,
		DeviceId:    e.DeviceId,
		RoomId:      e.RoomId,
		Action:      e.Action,
		CreatedDate: e.CreatedDate,
	}
}

func (ed EventDto) DomainToDtoCollection(events []domain.Event) []EventDto {
	result := make([]EventDto, len(events))

	for i := range events {
		result[i] = ed.DomainToDto(events[i])
	}

	return result
}
