package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type MeasurementController struct {
	MeasurementService app.MeasurementService
}

func NewMeasurementController(ms app.MeasurementService) MeasurementController {
	return MeasurementController{
		MeasurementService: ms,
	}
}

func (mc MeasurementController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m, err := requests.Bind(
			r,
			requests.MeasurementRequest{},
			domain.Measurement{},
		)

		if err != nil {
			log.Printf("MeasurementController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if org.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		m, err = mc.MeasurementService.Save(m)
		if err != nil {
			log.Printf("MeasurementController.Save(mc.MeasurementService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		mDto := resources.MeasurementDto{}
		mDto = mDto.DomainToDto(m)

		Success(w, mDto)
	}
}

func (mc MeasurementController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		m := r.Context().Value(MeasurementKey).(domain.Measurement)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.MeasurementDto{}.DomainToDto(m))
	}
}

func (mc MeasurementController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		m := r.Context().Value(MeasurementKey).(domain.Measurement)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		newM, err := requests.Bind(
			r,
			requests.MeasurementRequest{},
			domain.Measurement{},
		)

		if err != nil {
			log.Printf("MeasurementController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		m.DeviceId = newM.DeviceId
		m.RoomId = newM.RoomId
		m.Value = newM.Value

		m, err = mc.MeasurementService.Update(m)
		if err != nil {
			log.Printf("MeasurementController.Update(mc.MeasurementService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.MeasurementDto{}.DomainToDto(m))
	}
}

func (mc MeasurementController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		m := r.Context().Value(MeasurementKey).(domain.Measurement)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := mc.MeasurementService.Delete(m.Id)
		if err != nil {
			log.Printf("MeasurementController.Delete(mc.MeasurementService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

func (mc MeasurementController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		device := r.Context().Value(DeviceKey).(domain.Device)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		if org.Id != device.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		measurements, err := mc.MeasurementService.FindList(device.Id)
		if err != nil {
			log.Printf(
				"MeasurementController.FindList(mc.MeasurementService.FindList): %s",
				err,
			)
			InternalServerError(w, err)
			return
		}

		Success(
			w,
			resources.MeasurementDto{}.DomainToDtoCollection(measurements),
		)
	}
}
