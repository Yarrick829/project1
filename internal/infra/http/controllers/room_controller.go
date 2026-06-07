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

type RoomController struct {
	RoomService app.RoomService
}

func NewRoomController(rs app.RoomService) RoomController {
	return RoomController{
		RoomService: rs,
	}
}

func (rc RoomController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rm, err := requests.Bind(r,
			requests.RoomRequest{},
			domain.Room{})
		if err != nil {
			log.Printf("RoomController.Save(request.Bind): %s", err)
			BadRequest(w, err)
			return
		}
		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if org.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		rm.OrganizationId = org.Id

		rm, err = rc.RoomService.Save(rm)
		if err != nil {
			log.Printf("RoomController.Save(rc.RoomService.Save): %s", err)
			InternalServerError(w, err)
			return
		}
		rmDto := resources.RoomDto{}
		rmDto = rmDto.DomainToDto(rm)
		Success(w, rmDto)
	}
}

func (rc RoomController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		rm := r.Context().Value(RoomKey).(domain.Room)

		if org.Id != rm.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.RoomDto{}.DomainToDto(rm))
	}
}

func (rc RoomController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)
		rm := r.Context().Value(RoomKey).(domain.Room)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}
		if org.Id != rm.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}
		newRm, err := requests.Bind(r,
			requests.RoomRequest{},
			domain.Room{})
		if err != nil {
			log.Printf("RoomController.Update(request.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		rm.Name = newRm.Name
		rm.Description = newRm.Description

		rm, err = rc.RoomService.Update(rm)
		if err != nil {
			log.Printf("OrganizationController.Update(c.orgService.Update): %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.RoomDto{}.DomainToDto(rm))
	}
}

func (rc RoomController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		rm := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != rm.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := rc.RoomService.Delete(rm.Id)

		if err != nil {
			log.Printf("RoomController.Delete(c.orgService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}
		noContent(w)
	}
}

func (rc RoomController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		rooms, err := rc.RoomService.FindList(uint(org.Id))
		if err != nil {
			log.Printf("RoomController.FindList(rc.RoomService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.RoomDto{}.DomainToDtoCollection(rooms))
	}
}
