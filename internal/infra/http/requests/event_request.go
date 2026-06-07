package requests

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type EventRequest struct {
	DeviceId uint64             `json:"device_id" validate:"required"`
	RoomId   uint64             `json:"room_id" validate:"required"`
	Action   domain.EventAction `json:"action" validate:"required,oneof=ON OFF"`
}

func (er EventRequest) ToDomainModel() (interface{}, error) {
	return domain.Event{
		DeviceId: er.DeviceId,
		RoomId:   er.RoomId,
		Action:   er.Action,
	}, nil
}
