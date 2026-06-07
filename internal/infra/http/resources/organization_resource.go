package resources

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type OrganizationDto struct {
	Id          uint64    `json:"id"`
	UserId      uint64    `json:"userId"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Rooms       []RoomDto `json:"rooms,omitempty"`
}

func (d OrganizationDto) DomainToDto(o domain.Organization) OrganizationDto {
	roomDto := RoomDto{}

	return OrganizationDto{
		Id:          o.Id,
		UserId:      o.UserId,
		Name:        o.Name,
		Description: o.Description,
		City:        o.City,
		Address:     o.Address,
		Lat:         o.Lat,
		Lon:         o.Lon,
		Rooms:       roomDto.DomainToDtoCollection(o.Rooms),
	}
}

func (d OrganizationDto) DomainToDtoCollection(orgs []domain.Organization) []OrganizationDto {
	orgsDto := make([]OrganizationDto, len(orgs))

	for i := range orgs {
		orgsDto[i] = d.DomainToDto(orgs[i])
	}

	return orgsDto
}
