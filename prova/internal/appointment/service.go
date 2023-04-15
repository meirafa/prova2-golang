package appointment

import (
	"errors"

	"github.com/meirafa/prova2-golang/internal/domain"
)

type Service interface {
	//GetAll retorna todas consulta (appointment)
	GetAll() ([]domain.AppointmentDTO, error)
	//GetById retorna uma consulta (appointment) por id
	GetByID(id int) (domain.AppointmentDTO, error)
	// GetByDocumentPatient busca uma consulta pelo documento do paciente
	GetByDocumentPatient(Document string) ([]domain.AppointmentDTO, error)
	// Create cria uma nova consulta
	Create(a domain.Appointment) (domain.AppointmentDTO, error)
	//Update atualiza uma consulta
	Update(id int, a domain.Appointment) (domain.AppointmentDTO, error)
	//Delete exclui uma consulta
	Delete(id int) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetAll() ([]domain.AppointmentDTO, error) {
	list, err := s.r.GetAll()
	if err != nil {
		return nil, err
	}
	appointments, ok := list.([]domain.AppointmentDTO)
	if !ok {
		return nil, errors.New("an error has occurred while trying to fetch data from database")
	}
	return appointments, nil
}

func (s *service) GetByID(id int) (domain.AppointmentDTO, error) {
	aInterface, err := s.r.GetByID(id)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}
	appointment, ok := aInterface.(domain.AppointmentDTO)
	if !ok {
		return appointment, errors.New("an appointment with id provided has not been found")
	}
	return appointment, nil
}

func (s *service) GetByDocumentPatient(Document string) ([]domain.AppointmentDTO, error) {
	list, err := s.r.GetByDocumentPatient(Document)
	if err != nil {
		return nil, err
	}
	appointments, ok := list.([]domain.AppointmentDTO)
	if !ok {
		return nil, errors.New("an error occurred while trying to fetch data from database")
	}
	return appointments, nil
}

func (s *service) Create(a domain.Appointment) (domain.AppointmentDTO, error) {
	aSavedInterface, err := s.r.Create(a)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}
	apSaved, ok := aSavedInterface.(domain.AppointmentDTO)
	if ok {
		return apSaved, nil
	}

	return domain.AppointmentDTO{}, errors.New("failed to save new appointment")
}

func (s *service) Update(id int, a domain.Appointment) (domain.AppointmentDTO, error) {
	aUpdate, err := s.GetByID(id)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}

	if a.Description == "" {
		a.Description = aUpdate.Description
	}
	if a.AppointmentDate == "" {
		a.AppointmentDate = aUpdate.AppointmentDate
	}
	if a.IdDentist == "" {
		a.IdDentist = aUpdate.IdDentist
	}
	if a.IdPatient == "" {
		a.IdPatient = aUpdate.IdPatient
	}
	a.Id = aUpdate.Id

	updated, err := s.r.Update(id, a)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}
	response, ok := updated.(domain.AppointmentDTO)
	if !ok {
		return domain.AppointmentDTO{}, errors.New("failed to update appointment")
	}

	return response, nil
}

func (s *service) Delete(id int) error {
	return s.r.Delete(id)
}
