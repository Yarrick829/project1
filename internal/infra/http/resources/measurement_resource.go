package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type MeasurementDto struct {
	Id          uint64    `json:"id"`
	DeviceId    uint64    `json:"device_id"`
	RoomId      uint64    `json:"room_id"`
	Value       float64   `json:"value"`
	CreatedDate time.Time `json:"created_date"`
}

func (md MeasurementDto) DomainToDto(m domain.Measurement) MeasurementDto {
	return MeasurementDto{
		Id:          m.Id,
		DeviceId:    m.DeviceId,
		RoomId:      m.RoomId,
		Value:       m.Value,
		CreatedDate: m.CreatedDate,
	}
}

func (md MeasurementDto) DomainToDtoCollection(ms []domain.Measurement) []MeasurementDto {
	msDto := make([]MeasurementDto, len(ms))

	for i := range ms {
		msDto[i] = md.DomainToDto(ms[i])
	}

	return msDto
}
