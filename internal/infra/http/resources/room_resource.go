package resources

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type RoomDto struct {
	Id             uint64      `json:"id"`
	OrganizationId uint64      `json:"orgId"`
	Name           string      `json:"name"`
	Description    *string     `json:"description,omitempty"`
	Devices        []DeviceDto `json:"devices,omitempty"`
}

func (rd RoomDto) DomainToDto(rm domain.Room) RoomDto {
	deviceDto := DeviceDto{}

	return RoomDto{
		Id:             rm.Id,
		OrganizationId: rm.OrganizationId,
		Name:           rm.Name,
		Description:    rm.Description,
		Devices:        deviceDto.DomainToDtoCollection(rm.Devices),
	}
}

func (rd RoomDto) DomainToDtoCollection(rms []domain.Room) []RoomDto {
	rmsDto := make([]RoomDto, len(rms))
	for i := range rms {
		rmsDto[i] = rd.DomainToDto(rms[i])
	}
	return rmsDto
}
