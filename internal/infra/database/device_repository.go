package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/google/uuid"
	"github.com/upper/db/v4"
)

const deviceTableName = "devices"

type device struct {
	Id              uint64     `db:"id,omitempty"`
	OrganizationId  uint64     `db:"organization_id"`
	RoomId          *uint64    `db:"room_id"`
	GUID            string     `db:"guid"`
	InventoryNumber string     `db:"inventory_number"`
	SerialNumber    string     `db:"serial_number"`
	Characteristic  string     `db:"characteristics"`
	Category        string     `db:"category"`
	Units           *string    `db:"units"`
	PowerConsumtion *float64   `db:"power_consumption"`
	CreatedDate     time.Time  `db:"created_date"`
	UpdatedDate     time.Time  `db:"updated_date"`
	DeletedDate     *time.Time `db:"deleted_date"`
}

type DeviceRepository interface {
	FindByDeviceId(dId uint64) (domain.Device, error)
	FindList(orgId uint64) ([]domain.Device, error)
	FindByRoomId(roomId uint64) ([]domain.Device, error)
	Save(d domain.Device) (domain.Device, error)
	Update(d domain.Device) (domain.Device, error)
	Delete(id uint64) error
}

type deviceRepository struct {
	coll db.Collection
	sess db.Session
}

func NewDeviceRepository(session db.Session) DeviceRepository {
	return deviceRepository{
		coll: session.Collection(deviceTableName),
		sess: session,
	}
}

func (dr deviceRepository) FindByDeviceId(dId uint64) (domain.Device, error) {
	var dev device

	err := dr.coll.
		Find(db.Cond{"id": dId, "deleted_date": nil}).
		One(&dev)

	if err != nil {
		return domain.Device{}, err
	}

	d := dr.mapModelToDomain(dev)

	return d, nil
}

func (dr deviceRepository) FindList(orgId uint64) ([]domain.Device, error) {
	var devices []device

	err := dr.coll.
		Find(db.Cond{
			"organization_id": orgId,
			"deleted_date":    nil,
		}).
		All(&devices)

	if err != nil {
		return nil, err
	}

	return dr.mapModelToDomainCollection(devices), nil
}

func (dr deviceRepository) Save(d domain.Device) (domain.Device, error) {
	dev := dr.mapDomainToModel(d)

	now := time.Now()

	dev.CreatedDate = now
	dev.UpdatedDate = now
	dev.GUID = uuid.NewString()
	err := dr.coll.InsertReturning(&dev)
	if err != nil {
		return domain.Device{}, err
	}

	d = dr.mapModelToDomain(dev)

	return d, nil
}

func (dr deviceRepository) Update(d domain.Device) (domain.Device, error) {
	dev := dr.mapDomainToModel(d)

	dev.UpdatedDate = time.Now()

	err := dr.coll.
		Find(db.Cond{"id": d.Id, "deleted_date": nil}).
		Update(&dev)

	if err != nil {
		return domain.Device{}, err
	}

	d = dr.mapModelToDomain(dev)

	return d, nil
}

func (dr deviceRepository) Delete(id uint64) error {
	return dr.coll.
		Find(db.Cond{"id": id, "deleted_date": nil}).
		Update(map[string]interface{}{
			"deleted_date": time.Now(),
		})
}

func (dr deviceRepository) mapDomainToModel(d domain.Device) device {
	return device{
		Id:              d.Id,
		OrganizationId:  d.OrganizationId,
		RoomId:          d.RoomId,
		GUID:            d.GUID,
		InventoryNumber: d.InventoryNumber,
		SerialNumber:    d.SerialNumber,
		Characteristic:  d.Characteristic,
		Category:        string(d.Category),
		Units:           d.Units,
		PowerConsumtion: d.PowerConsumtion,
		CreatedDate:     d.CreatedDate,
		UpdatedDate:     d.UpdatedDate,
		DeletedDate:     d.DeletedDate,
	}
}

func (dr deviceRepository) mapModelToDomain(d device) domain.Device {
	return domain.Device{
		Id:              d.Id,
		OrganizationId:  d.OrganizationId,
		RoomId:          d.RoomId,
		GUID:            d.GUID,
		InventoryNumber: d.InventoryNumber,
		SerialNumber:    d.SerialNumber,
		Characteristic:  d.Characteristic,
		Category:        domain.DeviceCathegory(d.Category),
		Units:           d.Units,
		PowerConsumtion: d.PowerConsumtion,
		CreatedDate:     d.CreatedDate,
		UpdatedDate:     d.UpdatedDate,
		DeletedDate:     d.DeletedDate,
	}
}

func (dr deviceRepository) mapModelToDomainCollection(devices []device) []domain.Device {
	ds := make([]domain.Device, len(devices))

	for i := range devices {
		ds[i] = dr.mapModelToDomain(devices[i])
	}

	return ds
}

func (dr deviceRepository) FindByRoomId(roomId uint64) ([]domain.Device, error) {
	var devices []device

	err := dr.coll.
		Find(db.Cond{
			"room_id":      roomId,
			"deleted_date": nil,
		}).
		All(&devices)

	if err != nil {
		return nil, err
	}

	return dr.mapModelToDomainCollection(devices), nil
}
