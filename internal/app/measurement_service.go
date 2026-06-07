package app

import (
	"errors"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type measurementService struct {
	measurementRepo database.MeasurementRepository
}

type MeasurementService interface {
	Save(m domain.Measurement) (domain.Measurement, error)
	Update(m domain.Measurement) (domain.Measurement, error)
	Delete(id uint64) error
	Find(mId uint64) (interface{}, error)
	FindList(deviceId uint64) ([]domain.Measurement, error)
}

func NewMeasurementService(mr database.MeasurementRepository) MeasurementService {
	return measurementService{
		measurementRepo: mr,
	}
}

func (ms measurementService) Save(m domain.Measurement) (domain.Measurement, error) {

	err := ms.validateMeasurement(m)
	if err != nil {
		log.Printf("measurementService.Save(validateMeasurement): %s", err)
		return domain.Measurement{}, err
	}

	m, err = ms.measurementRepo.Save(m)
	if err != nil {
		log.Printf("measurementService.Save(ms.measurementRepo.Save): %s", err)
		return domain.Measurement{}, err
	}

	return m, nil
}

func (ms measurementService) Update(m domain.Measurement) (domain.Measurement, error) {

	err := ms.validateMeasurement(m)
	if err != nil {
		log.Printf("measurementService.Update(validateMeasurement): %s", err)
		return domain.Measurement{}, err
	}

	m, err = ms.measurementRepo.Update(m)
	if err != nil {
		log.Printf("measurementService.Update(ms.measurementRepo.Update): %s", err)
		return domain.Measurement{}, err
	}

	return m, nil
}

func (ms measurementService) Delete(id uint64) error {

	err := ms.measurementRepo.Delete(id)
	if err != nil {
		log.Printf("measurementService.Delete(ms.measurementRepo.Delete): %s", err)
		return err
	}

	return nil
}

func (ms measurementService) Find(mId uint64) (interface{}, error) {

	m, err := ms.measurementRepo.FindByMeasurementId(mId)
	if err != nil {
		log.Printf("measurementService.Find(ms.measurementRepo.FindByMeasurementId): %s", err)
		return nil, err
	}

	return m, nil
}

func (ms measurementService) validateMeasurement(m domain.Measurement) error {

	if m.DeviceId == 0 {
		return errors.New("device id is required")
	}

	if m.RoomId == 0 {
		return errors.New("room id is required")
	}

	return nil
}

func (ms measurementService) FindList(deviceId uint64) ([]domain.Measurement, error) {

	measurements, err := ms.measurementRepo.FindByDeviceId(deviceId)
	if err != nil {
		log.Printf(
			"measurementService.FindList(ms.measurementRepo.FindByDeviceId): %s",
			err,
		)
		return nil, err
	}

	return measurements, nil
}
