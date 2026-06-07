package database

import (
	"fmt"
	"math"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const measurementTableName = "measurements"

type measurement struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      uint64     `db:"room_id"`
	Value       float64    `db:"value"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type MeasurementRepository interface {
	FindByMeasurementId(mId uint64) (domain.Measurement, error)
	FindByDeviceId(dId uint64) ([]domain.Measurement, error)
	Save(m domain.Measurement) (domain.Measurement, error)
	Update(m domain.Measurement) (domain.Measurement, error)
	Delete(id uint64) error
	FindList(p domain.Pagination, mf MeasurementFilters) (domain.Measurements, error)
}

type measurementRepository struct {
	coll db.Collection
	sess db.Session
}

func NewMeasurementRepository(session db.Session) MeasurementRepository {
	return measurementRepository{
		coll: session.Collection(measurementTableName),
		sess: session,
	}
}

func (mr measurementRepository) FindByMeasurementId(mId uint64) (domain.Measurement, error) {
	var meas measurement

	err := mr.coll.
		Find(db.Cond{"id": mId, "deleted_date": nil}).
		One(&meas)

	if err != nil {
		return domain.Measurement{}, err
	}

	return mr.mapModelToDomain(meas), nil
}

func (mr measurementRepository) FindByDeviceId(dId uint64) ([]domain.Measurement, error) {
	var measurements []measurement

	err := mr.coll.
		Find(db.Cond{"device_id": dId, "deleted_date": nil}).
		All(&measurements)

	if err != nil {
		return nil, err
	}

	return mr.mapModelToDomainCollection(measurements), nil
}

func (mr measurementRepository) Save(m domain.Measurement) (domain.Measurement, error) {
	meas := mr.mapDomainToModel(m)

	now := time.Now()
	meas.CreatedDate = now
	meas.UpdatedDate = now

	err := mr.coll.InsertReturning(&meas)
	if err != nil {
		return domain.Measurement{}, err
	}

	return mr.mapModelToDomain(meas), nil
}

func (mr measurementRepository) Update(m domain.Measurement) (domain.Measurement, error) {
	meas := mr.mapDomainToModel(m)

	meas.UpdatedDate = time.Now()

	err := mr.coll.
		Find(db.Cond{"id": m.Id, "deleted_date": nil}).
		Update(&meas)

	if err != nil {
		return domain.Measurement{}, err
	}

	return mr.mapModelToDomain(meas), nil
}

func (mr measurementRepository) Delete(id uint64) error {
	return mr.coll.
		Find(db.Cond{"id": id, "deleted_date": nil}).
		Update(map[string]interface{}{
			"deleted_date": time.Now(),
		})
}

func (mr measurementRepository) mapDomainToModel(m domain.Measurement) measurement {
	return measurement{
		Id:          m.Id,
		DeviceId:    m.DeviceId,
		RoomId:      m.RoomId,
		Value:       m.Value,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (mr measurementRepository) mapModelToDomain(m measurement) domain.Measurement {
	return domain.Measurement{
		Id:          m.Id,
		DeviceId:    m.DeviceId,
		RoomId:      m.RoomId,
		Value:       m.Value,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (mr measurementRepository) mapModelToDomainCollection(measurements []measurement) []domain.Measurement {
	ms := make([]domain.Measurement, len(measurements))

	for i := range measurements {
		ms[i] = mr.mapModelToDomain(measurements[i])
	}

	return ms
}

func (mr measurementRepository) FindByDeviceIdAndPeriod(
	dId uint64,
	from time.Time,
	to time.Time,
) ([]domain.Measurement, error) {

	var measurements []measurement

	err := mr.coll.
		Find(db.Cond{
			"device_id":       dId,
			"deleted_date":    nil,
			"created_date >=": from,
			"created_date <=": to,
		}).
		All(&measurements)

	if err != nil {
		return nil, err
	}

	return mr.mapModelToDomainCollection(measurements), nil
}

type MeasurementFilters struct {
	DeviceId uint64
	RoomId   uint64

	CreatedDateFrom *time.Time
	CreatedDateTo   *time.Time

	Sort string
}

func (r measurementRepository) FindList(
	p domain.Pagination,
	mf MeasurementFilters,
) (domain.Measurements, error) {

	var ms []measurement

	if mf.DeviceId == 0 {
		return domain.Measurements{}, fmt.Errorf("device_id is required")
	}

	if mf.RoomId == 0 {
		return domain.Measurements{}, fmt.Errorf("room_id is required")
	}

	if p.Page == 0 {
		p.Page = 1
	}

	if p.CountPerPage == 0 {
		p.CountPerPage = 20
	}

	query := r.coll.
		Find(db.Cond{
			"device_id":    mf.DeviceId,
			"room_id":      mf.RoomId,
			"deleted_date": nil,
		})

	if mf.CreatedDateFrom != nil {
		query = query.And("created_date >= ?", *mf.CreatedDateFrom)
	}

	if mf.CreatedDateTo != nil {
		query = query.And("created_date <= ?", *mf.CreatedDateTo)
	}

	switch mf.Sort {
	case "created_date":
		query = query.OrderBy("created_date")
	case "-created_date":
		query = query.OrderBy("-created_date")
	case "value":
		query = query.OrderBy("value")
	case "-value":
		query = query.OrderBy("-value")
	default:
		query = query.OrderBy("-created_date")
	}

	res := query.Paginate(uint(p.CountPerPage))

	err := res.Page(uint(p.Page)).All(&ms)
	if err != nil {
		return domain.Measurements{}, err
	}

	measurements := r.mapModelToDomainPagination(ms)

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Measurements{}, err
	}

	measurements.Total = totalCount
	measurements.Pages = uint(
		math.Ceil(
			float64(measurements.Total) /
				float64(p.CountPerPage),
		),
	)

	return measurements, nil
}

func (r measurementRepository) mapModelToDomainPagination(
	ms []measurement,
) domain.Measurements {

	measurements := make([]domain.Measurement, 0, len(ms))

	for _, m := range ms {
		measurements = append(measurements, domain.Measurement{
			Id:          m.Id,
			DeviceId:    m.DeviceId,
			RoomId:      m.RoomId,
			Value:       m.Value,
			CreatedDate: m.CreatedDate,
			UpdatedDate: m.UpdatedDate,
			DeletedDate: m.DeletedDate,
		})
	}

	return domain.Measurements{
		Items: measurements,
	}
}
