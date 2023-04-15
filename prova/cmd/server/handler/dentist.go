package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meirafa/prova2-golang/internal/dentist"
	"github.com/meirafa/prova2-golang/internal/domain"
	"github.com/meirafa/prova2-golang/pkg/web"
)

type dentistHandler struct {
	s dentist.Service
}

// NewAppointmentHandler cria um novo controller de dentista
func NewDentistHandler(s dentist.Service) *dentistHandler {
	return &dentistHandler{
		s: s,
	}
}

//GetAll retorna todos os dentistas (dentist) cadastrados
func (h *dentistHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

//GetByID retorna um dentista (dentist) por id
func (h *dentistHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		response, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "dentist not found")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Post insere um novo dentista
func (h *dentistHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dentist domain.Dentist
		err := ctx.ShouldBindJSON(&dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid dentist")
			return
		}

		isValid, err := isEmptyDentist(&dentist)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Create(dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusCreated, response)
	}
}

//Put atualiza um dentista
func (h *dentistHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id")
			return
		}

		_, err = h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, 404, "error", "dentist not found")
			return
		}
		if err != nil {
			web.BadResponse(ctx, 409, "error", err.Error())
			return
		}

		var dentist domain.Dentist
		err = ctx.ShouldBindJSON(&dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid dentist data")
		}

		isValid, err := isEmptyDentist(&dentist)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Update(id, dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusConflict, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Patch atualiza um dentista ou algum de seus campos
func (h *dentistHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Surname      string `json:"surname,omitempty"`
		Name         string `json:"name,omitempty"`
		Registration string `json:"registration,omitempty"`
	}
	return func(ctx *gin.Context) {
		var r Request
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		_, err = h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, 404, "error", "dentist not found")
			return
		}

		if err := ctx.ShouldBindJSON(&r); err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid request")
			return
		}
		update := domain.Dentist{
			Surname:      r.Surname,
			Name:         r.Name,
			Registration: r.Registration,
		}

		updated, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, updated)
	}
}

//Delete exclui um dentista
func (h *dentistHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.DeleteResponse(ctx, http.StatusOK, "dentist deleted")
	}
}

// isEmptyDentist valida se os campos não estão vazios
func isEmptyDentist(dentist *domain.Dentist) (bool, error) {
	if dentist.Surname == "" || dentist.Name == "" || dentist.Registration == "" {
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
