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

type DeviceController struct {
	DeviceService app.DeviceService
}

func NewDeviceController(ds app.DeviceService) DeviceController {
	return DeviceController{
		DeviceService: ds,
	}
}

func (dc DeviceController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		d, err := requests.Bind(r,
			requests.DeviceRequest{},
			domain.Device{})

		if err != nil {
			log.Printf("DeviceController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if org.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		d.OrganizationId = org.Id

		d, err = dc.DeviceService.Save(d)
		if err != nil {
			log.Printf("DeviceController.Save(dc.DeviceService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		dDto := resources.DeviceDto{}
		dDto = dDto.DomainToDto(d)

		Success(w, dDto)
	}
}

func (dc DeviceController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		d := r.Context().Value(DeviceKey).(domain.Device)

		if org.Id != d.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDto(d))
	}
}

func (dc DeviceController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		d := r.Context().Value(DeviceKey).(domain.Device)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		if org.Id != d.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		newD, err := requests.Bind(r,
			requests.DeviceRequest{},
			domain.Device{})

		if err != nil {
			log.Printf("DeviceController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		d.RoomId = newD.RoomId
		d.InventoryNumber = newD.InventoryNumber
		d.SerialNumber = newD.SerialNumber
		d.Characteristic = newD.Characteristic
		d.Category = newD.Category
		d.Units = newD.Units
		d.PowerConsumtion = newD.PowerConsumtion

		d, err = dc.DeviceService.Update(d)
		if err != nil {
			log.Printf("DeviceController.Update(dc.DeviceService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDto(d))
	}
}

func (dc DeviceController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		d := r.Context().Value(DeviceKey).(domain.Device)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		if org.Id != d.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := dc.DeviceService.Delete(d.Id)
		if err != nil {
			log.Printf("DeviceController.Delete(dc.DeviceService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

func (dc DeviceController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		devices, err := dc.DeviceService.FindList(org.Id)
		if err != nil {
			log.Printf("DeviceController.FindList(dc.DeviceService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDtoCollection(devices))
	}
}
