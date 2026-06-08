package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type EventController struct {
	EventService app.EventService
}

func NewEventController(es app.EventService) EventController {
	return EventController{
		EventService: es,
	}
}

func (ec EventController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		e, err := requests.Bind(
			r,
			requests.EventRequest{},
			domain.Event{},
		)

		if err != nil {
			log.Printf("EventController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if org.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		e, err = ec.EventService.Save(e)
		if err != nil {
			log.Printf("EventController.Save(ec.EventService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		eDto := resources.EventDto{}
		eDto = eDto.DomainToDto(e)

		Success(w, eDto)
	}
}

func (ec EventController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		e := r.Context().Value(EventKey).(domain.Event)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.EventDto{}.DomainToDto(e))
	}
}

func (ec EventController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		e := r.Context().Value(EventKey).(domain.Event)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		newE, err := requests.Bind(
			r,
			requests.EventRequest{},
			domain.Event{},
		)

		if err != nil {
			log.Printf("EventController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		e.DeviceId = newE.DeviceId
		e.RoomId = newE.RoomId
		e.Action = newE.Action

		e, err = ec.EventService.Update(e)
		if err != nil {
			log.Printf("EventController.Update(ec.EventService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.EventDto{}.DomainToDto(e))
	}
}

func (ec EventController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		e := r.Context().Value(EventKey).(domain.Event)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := ec.EventService.Delete(e.Id)
		if err != nil {
			log.Printf("EventController.Delete(ec.EventService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

func (ec EventController) FindList() http.HandlerFunc {
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

		pagination := domain.Pagination{
			Page:         1,
			CountPerPage: 20,
		}

		filters := database.EventFilters{
			DeviceId: device.Id,
		}

		events, err := ec.EventService.FindList(
			pagination,
			filters,
		)

		if err != nil {
			log.Printf(
				"EventController.FindList(ec.EventService.FindList): %s",
				err,
			)
			InternalServerError(w, err)
			return
		}

		Success(
			w,
			resources.EventsDto{}.DomainPaginationToDto(events),
		)
	}
}
