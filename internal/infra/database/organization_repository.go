package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const OrganisationTableName = "organizations"

type organization struct {
	Id          uint64     `db:"id,omitempty"`
	UserId      uint64     `db:"user_id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	City        string     `db:"city"`
	Address     string     `db:"address"`
	Lat         float64    `db:"lat"`
	Lon         float64    `db:"lon"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type organizationRepository struct {
	coll db.Collection
	sess db.Session
}

type OrganizationRepository interface {
	Save(o domain.Organization) (domain.Organization, error)
	FindList(uId uint64) ([]domain.Organization, error)
	Find(id uint64) (domain.Organization, error)
	Update(o domain.Organization) (domain.Organization, error)
	Delete(id uint64) error
}

func NewOrganizationRepository(session db.Session) organizationRepository {
	return organizationRepository{
		sess: session,
		coll: session.Collection(OrganisationTableName),
	}
}

func (r organizationRepository) mapDomainToModel(o domain.Organization) organization {
	return organization{
		Id:          o.Id,
		UserId:      o.UserId,
		Name:        o.Name,
		Description: o.Description,
		City:        o.City,
		Address:     o.Address,
		Lat:         o.Lat,
		Lon:         o.Lon,
		CreatedDate: o.CreatedDate,
		UpdatedDate: o.UpdatedDate,
		DeletedDate: o.DeletedDate,
	}
}

func (r organizationRepository) Save(o domain.Organization) (domain.Organization, error) {
	org := r.mapDomainToModel(o)
	now := time.Now()
	org.CreatedDate = now
	org.UpdatedDate = now

	err := r.coll.InsertReturning(&org)
	if err != nil {
		return domain.Organization{}, err
	}

	o = r.mapModelToDomain(org)
	return o, nil
}

func (r organizationRepository) FindList(uId uint64) ([]domain.Organization, error) {
	var orgs []organization

	err := r.coll.
		Find(db.Cond{
			"user_id":      uId,
			"deleted_date": nil,
		}).All(&orgs)
	if err != nil {
		return nil, err
	}

	organizations := r.mapModelToDomainCollection(orgs)
	return organizations, nil

}
func (r organizationRepository) Find(id uint64) (domain.Organization, error) {
	var org organization
	err := r.coll.Find(db.Cond{
		"id": id, "deleted_date": nil}).
		One(&org)
	if err != nil {
		return domain.Organization{}, err
	}
	o := r.mapModelToDomain(org)
	return o, nil
}

func (r organizationRepository) Update(o domain.Organization) (domain.Organization, error) {
	org := r.mapDomainToModel(o)
	org.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": o.Id, "deleted_date": nil}).Update(&org)
	if err != nil {
		return domain.Organization{}, err
	}

	o = r.mapModelToDomain(org)
	return o, nil
}

func (r organizationRepository) mapModelToDomain(o organization) domain.Organization {
	return domain.Organization{
		Id:          o.Id,
		UserId:      o.UserId,
		Name:        o.Name,
		Description: o.Description,
		City:        o.City,
		Address:     o.Address,
		Lat:         o.Lat,
		Lon:         o.Lon,
		CreatedDate: o.CreatedDate,
		UpdatedDate: o.UpdatedDate,
		DeletedDate: o.DeletedDate,
	}
}

func (r organizationRepository) mapModelToDomainCollection(orgs []organization) []domain.Organization {
	organizations := make([]domain.Organization, len(orgs))
	for i := range orgs {
		organizations[i] = r.mapModelToDomain(orgs[i])
	}
	return organizations
}

func (r organizationRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}
