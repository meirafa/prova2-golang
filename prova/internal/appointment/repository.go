package appointment

import (
	"errors"
	"github.com/meirafa/prova2-golang/internal/domain"
	"github.com/meirafa/prova2-golang/pkg/store"
)

var table = store.AP

type Repository interface {
	//GetAll retorna todas consultas (appointment)
	GetAll() (interface{}, error)
	//GetById retorna uma consulta (appointment) por id
	GetByID(entityId int) (interface{}, error)
	// GetByDocumentPatient busca uma consulta pelo documento do paciente
	GetByDocumentPatient(Document string) (interface{}, error)
	// Create cria uma nova consulta
	Create(a domain.Appointment) (interface{}, error)
	//Update atualiza uma consulta
	Update(entityId int, a domain.Appointment) (interface{}, error)
	//Delete exclui uma consulta
	Delete(entityId int) error
}

type repository struct {
	store store.ApStore
}

//NewRepositoryAppointment cria um novo repositório
func NewRepository(store store.ApStore) Repository {
	return &repository{store}
}

func (r *repository) GetAll() (interface{}, error) {
	return r.store.GetAll(table)
}

func (r *repository) GetByID(entityId int) (interface{}, error) {
	return r.store.GetByID(entityId, table)
}

func (r *repository) GetByDocumentPatient(Document string) (interface{}, error) {
	return r.store.GetAllAppointmentsByPatientIdentify(Document)
}

func (r *repository) Create(a domain.Appointment) (interface{}, error) {
	return r.store.Save(a, table)
}

func (r *repository) Update(entityId int, a domain.Appointment) (interface{}, error) {
	aInterface, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	appointments, ok := aInterface.([]domain.AppointmentDTO)
	if !ok {
		return nil, err
	}

	for _, appointment := range appointments {
		if appointment.Id == entityId {
			return r.store.Update(entityId, a, table)
		}
	}
	return nil, errors.New("appointment not found")
}

func (r *repository) Delete(entityId int) error {
	return r.store.Delete(entityId, table)
}
