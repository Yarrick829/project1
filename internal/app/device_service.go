package app

import (
	"errors"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type deviceService struct {
	deviceRepo      database.DeviceRepository
	measurementRepo database.MeasurementRepository
}

type DeviceService interface {
	Save(d domain.Device) (domain.Device, error)
	Update(d domain.Device) (domain.Device, error)
	Delete(id uint64) error
	Find(dId uint64) (interface{}, error)
	FindList(orgId uint64) ([]domain.Device, error)
}

func NewDeviceService(dr database.DeviceRepository, mr database.MeasurementRepository) DeviceService {
	return deviceService{
		deviceRepo:      dr,
		measurementRepo: mr,
	}
}

func (ds deviceService) Save(d domain.Device) (domain.Device, error) {

	err := ds.validateDevice(d)
	if err != nil {
		log.Printf("deviceService.Save(validateDevice): %s", err)
		return domain.Device{}, err
	}

	d, err = ds.deviceRepo.Save(d)
	if err != nil {
		log.Printf("deviceService.Save(ds.deviceRepo.Save): %s", err)
		return domain.Device{}, err
	}

	return d, nil
}

func (ds deviceService) Update(d domain.Device) (domain.Device, error) {

	err := ds.validateDevice(d)
	if err != nil {
		log.Printf("deviceService.Update(validateDevice): %s", err)
		return domain.Device{}, err
	}

	d, err = ds.deviceRepo.Update(d)
	if err != nil {
		log.Printf("deviceService.Update(ds.deviceRepo.Update): %s", err)
		return domain.Device{}, err
	}

	return d, nil
}

func (ds deviceService) Delete(id uint64) error {

	err := ds.deviceRepo.Delete(id)
	if err != nil {
		log.Printf("deviceService.Delete(ds.deviceRepo.Delete): %s", err)
		return err
	}

	return nil
}

func (ds deviceService) Find(dId uint64) (interface{}, error) {

	d, err := ds.deviceRepo.FindByDeviceId(dId)
	if err != nil {
		return nil, err
	}

	d.Measurements, err =
		ds.measurementRepo.FindByDeviceId(d.Id)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func (ds deviceService) validateDevice(d domain.Device) error {

	switch d.Category {

	case domain.Sensor:
		if d.Units == nil || *d.Units == "" {
			return errors.New("units is required for sensor")
		}

	case domain.Actuator:
		if d.PowerConsumtion == nil {
			return errors.New("power consumption is required for actuator")
		}

	default:
		return errors.New("invalid device category")
	}

	return nil
}

func (ds deviceService) FindList(orgId uint64) ([]domain.Device, error) {

	devices, err := ds.deviceRepo.FindList(orgId)
	if err != nil {
		log.Printf("deviceService.FindList(ds.deviceRepo.FindList): %s", err)
		return nil, err
	}

	return devices, nil
}
