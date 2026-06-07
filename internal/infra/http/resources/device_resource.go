package resources

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type DeviceDto struct {
	Id              uint64                 `json:"id"`
	OrganizationId  uint64                 `json:"organization_id"`
	RoomId          *uint64                `json:"room_id,omitempty"`
	GUID            string                 `json:"guid"`
	InventoryNumber string                 `json:"inventory_number"`
	SerialNumber    string                 `json:"serial_number"`
	Characteristic  string                 `json:"characteristic"`
	Category        domain.DeviceCathegory `json:"category"`
	Units           *string                `json:"units,omitempty"`
	PowerConsumtion *float64               `json:"power_consumtion,omitempty"`
	Measurements    []MeasurementDto       `json:"measurements,omitempty"`
}

func (dd DeviceDto) DomainToDto(dm domain.Device) DeviceDto {

	measurementDto := MeasurementDto{}

	return DeviceDto{
		Id:              dm.Id,
		OrganizationId:  dm.OrganizationId,
		RoomId:          dm.RoomId,
		GUID:            dm.GUID,
		InventoryNumber: dm.InventoryNumber,
		SerialNumber:    dm.SerialNumber,
		Characteristic:  dm.Characteristic,
		Category:        dm.Category,
		Units:           dm.Units,
		PowerConsumtion: dm.PowerConsumtion,
		Measurements: measurementDto.
			DomainToDtoCollection(dm.Measurements),
	}
}

func (dd DeviceDto) DomainToDtoCollection(dms []domain.Device) []DeviceDto {
	dmsDto := make([]DeviceDto, len(dms))

	for i := range dms {
		dmsDto[i] = dd.DomainToDto(dms[i])
	}

	return dmsDto
}
