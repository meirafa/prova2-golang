package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meirafa/prova2-golang/internal/domain"
	"github.com/meirafa/prova2-golang/internal/patient"
	"github.com/meirafa/prova2-golang/pkg/web"
)

type patientHandler struct {
	s patient.Service
}

// NewAppointmentHandler cria um novo controller de paciente
func NewPatientHandler(s patient.Service) *patientHandler {
	return &patientHandler{
		s: s,
	}
}

//GetAll retorna todos os pacientes (patient) cadastrados
func (h *patientHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		patients, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, patients)
	}
}

//GetByID retorna um paciente (patient) por id
func (h *patientHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		patient, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "patient not found")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, patient)
	}
}

// Post insere um novo paciente
func (h *patientHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var patient domain.Patient
		err := ctx.ShouldBindJSON(&patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient")
			return
		}

		isValid, err := isEmptyPatient(&patient)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Create(patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		web.ResponseOK(ctx, http.StatusCreated, response)
	}
}

//Put atualiza um paciente
func (h *patientHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient id provided")
			return
		}
		_, err = h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, 404, "error", "patient not found")
			return
		}
		if err != nil {
			web.BadResponse(ctx, 409, "error", err.Error())
			return
		}

		var patient domain.Patient
		err = ctx.ShouldBindJSON(&patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient data")
		}

		isValid, err := isEmptyPatient(&patient)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Update(id, patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusConflict, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Patch atualiza um paciente ou algum de seus campos
func (h *patientHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Surname   string `json:"surname,omitempty"`
		Name      string `json:"name,omitempty"`
		Document  string `json:"document,omitempty"`
		CreatedAt string `json:"created_at,omitempty"`
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
			web.BadResponse(ctx, 404, "error", "patient not found")
			return
		}
		if err := ctx.ShouldBindJSON(&r); err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid request")
			return
		}
		update := domain.Patient{
			Surname:   r.Surname,
			Name:      r.Name,
			Document:  r.Document,
			CreatedAt: r.CreatedAt,
		}
		response, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

//Delete exclui um paciente
func (h *patientHandler) Delete() gin.HandlerFunc {
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
		web.DeleteResponse(ctx, http.StatusOK, "patient deleted")
	}
}

// isEmptyPatient valida se os campos não estão vazios
func isEmptyPatient(patient *domain.Patient) (bool, error) {
	if patient.Surname == "" || patient.Name == "" || patient.CreatedAt == "" || patient.Document == "" {
		return false, errors.New("patient fields can't be empty")
	}
	return true, nil
}
