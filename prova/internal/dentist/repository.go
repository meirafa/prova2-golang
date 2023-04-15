package dentist

import (
	"errors"
	"log"

	"github.com/meirafa/prova2-golang/internal/domain"
	"github.com/meirafa/prova2-golang/pkg/store"
)

var table = "dentists"

type Repository interface {
	//GetAll retorna todos os dentistas (dentist) cadastrados
	GetAll() (interface{}, error)
	//GetByID retorna um dentista (dentist) por id
	GetByID(id int) (interface{}, error)
	// Create insere um novo dentista
	Create(d domain.Dentist) (interface{}, error)
	//Update atualiza um dentista
	Update(id int, d domain.Dentist) (interface{}, error)
	//Delete exclui um dentista
	Delete(id int) error
}

type repository struct {
	store store.Store
}

//NewRepository cria um novo repositório
func NewRepository(store store.Store) Repository {
	return &repository{store}
}

func (r *repository) GetAll() (interface{}, error) {
	return r.store.GetAll(table)
}

func (r *repository) GetByID(id int) (interface{}, error) {
	return r.store.GetByID(id, table)
}

func (r *repository) Create(d domain.Dentist) (interface{}, error) {
	if !r.validateRegistration(d.Registration) {
		return nil, errors.New("license number already exists on database")
	}
	return r.store.Save(d, table)
}

func (r *repository) Update(id int, d domain.Dentist) (interface{}, error) {
	var dentists []domain.Dentist
	dentistsInterface, err := r.GetAll()
	if err != nil {
		log.Fatalln("error while trying to fetch data from db")
		return nil, err
	}
	dentists, ok := dentistsInterface.([]domain.Dentist)
	if !ok {

		return nil, err
	}

	for _, dentist := range dentists {
		if dentist.Id == id {

			if !r.validateRegistration(d.Registration) && d.Registration != dentist.Registration {
				return nil, errors.New("license number already exists")
			}

			return r.store.Update(id, d, table)
		}
	}
	return nil, errors.New("dentist not found")
}

func (r *repository) Delete(id int) error {

	return r.store.Delete(id, table)
}

func (r *repository) validateRegistration(Registration string) bool {

	var dentists []domain.Dentist
	dentistsInterface, err := r.GetAll()
	if err != nil {
		return false
	}
	dentists, ok := dentistsInterface.([]domain.Dentist)
	if !ok {
		return false
	}

	for _, dentist := range dentists {

		if dentist.Registration == Registration {
			return false
		}
	}
	return true
}
