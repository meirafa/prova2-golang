package dentist

import (
	"errors"

	"github.com/meirafa/prova2-golang/internal/domain"
)

type Service interface {
	//GetAll retorna todos os dentistas (dentist) cadastrados
	GetAll() ([]domain.Dentist, error)
	//GetByID retorna um dentista (dentist) por id
	GetByID(id int) (interface{}, error)
	// Create insere um novo dentista
	Create(d domain.Dentist) (domain.Dentist, error)
	//Update atualiza um dentista
	Update(id int, d domain.Dentist) (domain.Dentist, error)
	//Delete exclui um dentista
	Delete(id int) error
}

type service struct {
	r Repository
}

// NewService cria um novo serviço
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetAll() ([]domain.Dentist, error) {
	list, err := s.r.GetAll()
	if err != nil {
		return nil, err
	}

	dentists, ok := list.([]domain.Dentist)
	if !ok {
		return nil, errors.New("dentist with id provided has not been found")
	}
	return dentists, nil
}

func (s *service) GetByID(id int) (interface{}, error) {
	dInterface, err := s.r.GetByID(id)
	if err != nil {
		return nil, err
	}
	dentist, ok := dInterface.(domain.Dentist)
	if !ok {
		return dentist, errors.New("an error occurred while trying to fetch data from db")
	}
	return dentist, nil
}

func (s *service) Create(d domain.Dentist) (domain.Dentist, error) {
	dSavedInterface, err := s.r.Create(d)
	if err != nil {
		return domain.Dentist{}, err
	}

	dentistSaved, ok := dSavedInterface.(domain.Dentist)
	if ok {
		return dentistSaved, nil
	}

	return domain.Dentist{}, errors.New("failed to save a new dentist in db")
}

func (s *service) Update(id int, d domain.Dentist) (domain.Dentist, error) {
	dUpdatedInterface, err := s.r.Update(id, d)

	if err != nil {

		return domain.Dentist{}, err
	}
	dentistUpdated, ok := dUpdatedInterface.(domain.Dentist)
	if ok {
		return dentistUpdated, nil
	}

	return domain.Dentist{}, errors.New("failed to update the dentist entry")
}

func (s *service) Delete(id int) error {
	return s.r.Delete(id)

}
