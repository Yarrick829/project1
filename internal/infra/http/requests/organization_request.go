package requests

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type OrganizationRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	City        string  `json:"city" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
}

func (r OrganizationRequest) ToDomainModel() (interface{}, error) {
	return domain.Organization{
		Name:        r.Name,
		Description: r.Description,
		City:        r.City,
		Address:     r.Address,
		Lat:         r.Lat,
		Lon:         r.Lon,
	}, nil
}
