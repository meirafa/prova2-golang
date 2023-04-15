package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/meirafa/prova2-golang/internal/appointment"
	"github.com/meirafa/prova2-golang/internal/domain"
	"github.com/meirafa/prova2-golang/pkg/web"
	"net/http"
	"strconv"
)

type appointmentHandler struct {
	s appointment.Service
}

// NewAppointmentHandler cria um novo controller de agendamentos
func NewAppointmentHandler(s appointment.Service) *appointmentHandler {
	return &appointmentHandler{
		s: s,
	}
}

// GetAll retorna todas consultas (appointments)
func (h *appointmentHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		if response == nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "was not found appointments registered")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

//GetById retorna uma consulta (appointment) por id
func (h *appointmentHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		response, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// GetByDocumentPatient busca uma consulta pelo documento do paciente
func (h *appointmentHandler) GetByDocumentPatient() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("document")
		_, err := strconv.Atoi(idParam)

		if err != nil {
			web.BadResponse(ctx, 400, "error", "invalid document")
			return
		}

		response, err := h.s.GetByDocumentPatient(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Post cria uma nova consulta
func (h *appointmentHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var appointment domain.Appointment
		err := ctx.ShouldBindJSON(&appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid appointment data, please verify field(s): "+err.Error())
			return
		}

		isValid, err := isEmptyAppointment(&appointment)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		response, err := h.s.Create(appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

//Put atualiza uma consulta
func (h *appointmentHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		_, err = h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, 404, "error", "appointment not found")
			return
		}
		if err != nil {
			web.BadResponse(ctx, 409, "error", err.Error())
			return
		}

		var appointment domain.Appointment
		err = ctx.ShouldBindJSON(&appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid appointment data, verify the fields and try again")
			return
		}

		isValid, err := isEmptyAppointment(&appointment)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		response, err := h.s.Update(id, appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// Patch atualiza uma consulta ou algum de seus campos
func (h *appointmentHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Description     string `json:"description,omitempty"`
		AppointmentDate string `json:"appointment_date,omitempty"`
		IdDentist       string `json:"id_dentist,omitempty"`
		IdPatient       string `json:"id_patient,omitempty"`
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
			web.BadResponse(ctx, 404, "error", "appointment not found")
			return
		}

		if err := ctx.ShouldBindJSON(&r); err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid request")
			return
		}
		update := domain.Appointment{
			Description:     r.Description,
			AppointmentDate: r.AppointmentDate,
			IdDentist:       r.IdDentist,
			IdPatient:       r.IdPatient,
		}

		response, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

//Delete exclui uma consulta
func (h *appointmentHandler) Delete() gin.HandlerFunc {
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
		web.DeleteResponse(ctx, http.StatusOK, "appointment removed")
	}
}

// isEmptyAppointment valida se os campos não estão vazios
func isEmptyAppointment(appointment *domain.Appointment) (bool, error) {
	if appointment.Description == "" || appointment.IdDentist == "" || appointment.AppointmentDate == "" || appointment.IdPatient == "" {
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
