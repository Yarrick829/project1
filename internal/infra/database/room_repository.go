package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const roomTableName = "rooms"

type room struct {
	Id             uint64     `db:"id,omitempty"`
	OrganizationId uint64     `db:"organization_id"`
	Name           string     `db:"name"`
	Description    *string    `db:"description"`
	CreatedDate    time.Time  `db:"created_date"`
	UpdatedDate    time.Time  `db:"updated_date"`
	DeletedDate    *time.Time `db:"deleted_date"`
}

type RoomRepository interface {
	FindByRoomId(rId uint64) (domain.Room, error)
	FindByOrgId(oId uint64) ([]domain.Room, error)
	Save(o domain.Room) (domain.Room, error)
	Update(rm domain.Room) (domain.Room, error)
	Delete(id uint64) error
	FindByOrganizationId(orgId uint64) ([]domain.Room, error)
}

type roomRepository struct {
	coll db.Collection
	sess db.Session
}

func (r roomRepository) FindByRoomId(rId uint64) (domain.Room, error) {
	var rm room

	err := r.coll.
		Find(db.Cond{"id": rId, "deleted_date": nil}).
		One(&rm)
	if err != nil {
		return domain.Room{}, err
	}
	o := r.mapModelToDomain(rm)
	return o, nil
}

func NewRoomRepository(session db.Session) RoomRepository {
	return roomRepository{
		coll: session.Collection(roomTableName),
		sess: session,
	}
}

func (r roomRepository) Save(rm domain.Room) (domain.Room, error) {
	rms := r.mapDomainToModel(rm)
	now := time.Now()
	rms.CreatedDate = now
	rms.UpdatedDate = now

	err := r.coll.InsertReturning(&rms)
	if err != nil {
		return domain.Room{}, err
	}

	rm = r.mapModelToDomain(rms)
	return rm, nil
}

func (r roomRepository) FindByOrgId(oId uint64) ([]domain.Room, error) {
	var rooms []room

	err := r.coll.Find(db.Cond{"organization_id": oId, "deleted_date": nil}).All(&rooms)
	if err != nil {
		return nil, err
	}
	rms := r.mapModelToDomainCollection(rooms)
	return rms, nil
}

func (r roomRepository) Update(rm domain.Room) (domain.Room, error) {
	rms := r.mapDomainToModel(rm)
	rms.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": rm.Id, "deleted_date": nil}).Update(&rms)
	if err != nil {
		return domain.Room{}, err
	}

	rm = r.mapModelToDomain(rms)
	return rm, nil
}

func (r roomRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r roomRepository) mapDomainToModel(rm domain.Room) room {
	return room{
		Id:             rm.Id,
		OrganizationId: rm.OrganizationId,
		Name:           rm.Name,
		Description:    rm.Description,
		CreatedDate:    rm.CreatedDate,
		UpdatedDate:    rm.UpdatedDate,
		DeletedDate:    rm.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomain(rm room) domain.Room {
	return domain.Room{
		Id:             rm.Id,
		OrganizationId: rm.OrganizationId,
		Name:           rm.Name,
		Description:    rm.Description,
		CreatedDate:    rm.CreatedDate,
		UpdatedDate:    rm.UpdatedDate,
		DeletedDate:    rm.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomainCollection(rooms []room) []domain.Room {
	rms := make([]domain.Room, len(rooms))
	for i := range rooms {
		rms[i] = r.mapModelToDomain(rooms[i])
	}
	return rms
}

func (r roomRepository) FindByOrganizationId(orgId uint64) ([]domain.Room, error) {
	var rms []room

	err := r.coll.
		Find(db.Cond{
			"organization_id": orgId,
			"deleted_date":    nil,
		}).
		All(&rms)

	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(rms), nil
}
