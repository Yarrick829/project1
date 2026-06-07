package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const eventTableName = "events"

type event struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      uint64     `db:"room_id"`
	Action      string     `db:"action"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type EventRepository interface {
	FindByEventId(eId uint64) (domain.Event, error)
	FindByDeviceId(dId uint64) ([]domain.Event, error)
	FindByDeviceIdAndPeriod(
		dId uint64,
		from time.Time,
		to time.Time,
	) ([]domain.Event, error)
	Save(e domain.Event) (domain.Event, error)
	Update(e domain.Event) (domain.Event, error)
	Delete(id uint64) error
}

type eventRepository struct {
	coll db.Collection
	sess db.Session
}

func NewEventRepository(session db.Session) EventRepository {
	return eventRepository{
		coll: session.Collection(eventTableName),
		sess: session,
	}
}

func (er eventRepository) FindByEventId(eId uint64) (domain.Event, error) {
	var ev event

	err := er.coll.
		Find(db.Cond{
			"id":           eId,
			"deleted_date": nil,
		}).
		One(&ev)

	if err != nil {
		return domain.Event{}, err
	}

	return er.mapModelToDomain(ev), nil
}

func (er eventRepository) FindByDeviceId(dId uint64) ([]domain.Event, error) {
	var events []event

	err := er.coll.
		Find(db.Cond{
			"device_id":    dId,
			"deleted_date": nil,
		}).
		All(&events)

	if err != nil {
		return nil, err
	}

	return er.mapModelToDomainCollection(events), nil
}

func (er eventRepository) FindByDeviceIdAndPeriod(
	dId uint64,
	from time.Time,
	to time.Time,
) ([]domain.Event, error) {

	var events []event

	err := er.coll.
		Find(db.Cond{
			"device_id":       dId,
			"deleted_date":    nil,
			"created_date >=": from,
			"created_date <=": to,
		}).
		All(&events)

	if err != nil {
		return nil, err
	}

	return er.mapModelToDomainCollection(events), nil
}

func (er eventRepository) Save(e domain.Event) (domain.Event, error) {

	ev := er.mapDomainToModel(e)

	now := time.Now()

	ev.CreatedDate = now
	ev.UpdatedDate = now

	err := er.coll.InsertReturning(&ev)
	if err != nil {
		return domain.Event{}, err
	}

	return er.mapModelToDomain(ev), nil
}

func (er eventRepository) Update(e domain.Event) (domain.Event, error) {

	ev := er.mapDomainToModel(e)

	ev.UpdatedDate = time.Now()

	err := er.coll.
		Find(db.Cond{
			"id":           e.Id,
			"deleted_date": nil,
		}).
		Update(&ev)

	if err != nil {
		return domain.Event{}, err
	}

	return er.mapModelToDomain(ev), nil
}

func (er eventRepository) Delete(id uint64) error {

	return er.coll.
		Find(db.Cond{
			"id":           id,
			"deleted_date": nil,
		}).
		Update(map[string]interface{}{
			"deleted_date": time.Now(),
		})
}

func (er eventRepository) mapDomainToModel(e domain.Event) event {
	return event{
		Id:          e.Id,
		DeviceId:    e.DeviceId,
		RoomId:      e.RoomId,
		Action:      string(e.Action),
		CreatedDate: e.CreatedDate,
		UpdatedDate: e.UpdatedDate,
		DeletedDate: e.DeletedDate,
	}
}

func (er eventRepository) mapModelToDomain(e event) domain.Event {
	return domain.Event{
		Id:          e.Id,
		DeviceId:    e.DeviceId,
		RoomId:      e.RoomId,
		Action:      domain.EventAction(e.Action),
		CreatedDate: e.CreatedDate,
		UpdatedDate: e.UpdatedDate,
		DeletedDate: e.DeletedDate,
	}
}

func (er eventRepository) mapModelToDomainCollection(events []event) []domain.Event {

	es := make([]domain.Event, len(events))

	for i := range events {
		es[i] = er.mapModelToDomain(events[i])
	}

	return es
}
