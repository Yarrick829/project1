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

type OrganizationController struct {
	orgService app.OrganizationService
}

func NewOrganizationController(os app.OrganizationService) OrganizationController {
	return OrganizationController{
		orgService: os,
	}
}

func (c OrganizationController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		org, err := requests.Bind(r,
			requests.OrganizationRequest{},
			domain.Organization{})
		if err != nil {
			log.Printf("OrganizationController.Save(request.Bind): %s", err)
			BadRequest(w, err)
			return
		}
		user := r.Context().Value(UserKey).(domain.User)
		org.UserId = user.Id

		org, err = c.orgService.Save(org)
		if err != nil {
			log.Printf("OrganizationController.Save(c.orgService.Save): %s", err)
			InternalServerError(w, err)
			return
		}
		orgDto := resources.OrganizationDto{}
		orgDto = orgDto.DomainToDto(org)
		Success(w, orgDto)
	}
}

func (c OrganizationController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)

		orgs, err := c.orgService.FindList(user.Id)
		if err != nil {
			log.Printf("OrganizationController.FindList(c.orgService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.OrganizationDto{}.DomainToDtoCollection(orgs))
	}
}

func (c OrganizationController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.OrganizationDto{}.DomainToDto(org))
	}
}

func (c OrganizationController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		newOrg, err := requests.Bind(r,
			requests.OrganizationRequest{},
			domain.Organization{})
		if err != nil {
			log.Printf("OrganizationController.Update(request.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		org.Name = newOrg.Name
		org.Description = newOrg.Description
		org.City = newOrg.City
		org.Address = newOrg.Address
		org.Lat = newOrg.Lat
		org.Lon = newOrg.Lon

		org, err = c.orgService.Update(org)
		if err != nil {
			log.Printf("OrganizationController.Update(c.orgService.Update): %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.OrganizationDto{}.DomainToDto(org))
	}
}

func (c OrganizationController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := c.orgService.Delete(org.Id)

		if err != nil {
			log.Printf("OrganizationController.Delete(c.orgService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}
		noContent(w)
	}
}
