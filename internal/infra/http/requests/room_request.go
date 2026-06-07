package requests

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type RoomRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"required"`
}

func (rr RoomRequest) ToDomainModel() (interface{}, error) {
	return domain.Room{
		Name:        rr.Name,
		Description: rr.Description,
	}, nil
}
