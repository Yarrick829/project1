package requests

import (
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type PaginationRequest struct {
	Page         string `query:"page"`
	CountPerPage string `query:"count_per_page"`
}

func (pr PaginationRequest) ToDomainModel() (interface{}, error) {
	p := domain.Pagination{
		Page:         1,
		CountPerPage: 20,
	}

	if pr.Page != "" {
		page, err := strconv.ParseUint(pr.Page, 10, 64)
		if err != nil {
			return nil, err
		}
		p.Page = page
	}

	if pr.CountPerPage != "" {
		count, err := strconv.ParseUint(pr.CountPerPage, 10, 64)
		if err != nil {
			return nil, err
		}
		p.CountPerPage = count
	}

	return p, nil
}
