package requests

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type MeasurementRequest struct {
	DeviceId uint64  `json:"device_id" validate:"required"`
	RoomId   uint64  `json:"room_id" validate:"required"`
	Value    float64 `json:"value" validate:"required"`
}

func (mr MeasurementRequest) ToDomainModel() (interface{}, error) {
	return domain.Measurement{
		DeviceId: mr.DeviceId,
		RoomId:   mr.RoomId,
		Value:    mr.Value,
	}, nil
}
