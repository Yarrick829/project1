package requests

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type DeviceRequest struct {
	RoomId          *uint64                `json:"room_id"`
	InventoryNumber string                 `json:"inventory_number" validate:"required"`
	SerialNumber    string                 `json:"serial_number" validate:"required"`
	Characteristic  string                 `json:"characteristic" validate:"required"`
	Category        domain.DeviceCathegory `json:"category" validate:"required,oneof=SENSOR ACTUATOR"`
	Units           *string                `json:"units" validate:"required"`
	PowerConsumtion *float64               `json:"power_consumtion" validate:"required"`
}

func (dr DeviceRequest) ToDomainModel() (interface{}, error) {
	return domain.Device{
		RoomId:          dr.RoomId,
		InventoryNumber: dr.InventoryNumber,
		SerialNumber:    dr.SerialNumber,
		Characteristic:  dr.Characteristic,
		Category:        dr.Category,
		Units:           dr.Units,
		PowerConsumtion: dr.PowerConsumtion,
	}, nil
}
