package patient

import (
	"errors"

	"github.com/meirafa/prova2-golang/internal/domain"
)

type Service interface {
	GetAll() ([]domain.Patient, error)
	GetByID(id int) (domain.Patient, error)
	Create(p domain.Patient) (domain.Patient, error)
	Update(id int, p domain.Patient) (domain.Patient, error)
	Delete(id int) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetAll() ([]domain.Patient, error) {
	list, err := s.r.GetAll()
	if err != nil {
		return nil, err
	}
	patients, ok := list.([]domain.Patient)
	if !ok {
		return nil, errors.New("an error occurred while trying to fetch data from db")
	}
	return patients, nil
}

func (s *service) GetByID(id int) (domain.Patient, error) {
	pInterface, err := s.r.GetByID(id)
	if err != nil {
		return domain.Patient{}, err
	}
	patient, ok := pInterface.(domain.Patient)
	if !ok {
		return patient, errors.New("not found a patient with provided id")
	}
	return patient, nil
}

func (s *service) Create(p domain.Patient) (domain.Patient, error) {
	pSavedInterface, err := s.r.Create(p)
	if err != nil {
		return domain.Patient{}, err
	}

	patientSaved, ok := pSavedInterface.(domain.Patient)
	if ok {
		return patientSaved, nil
	}

	return domain.Patient{}, errors.New("failed to save a new patient at db")
}

func (s *service) Update(id int, p domain.Patient) (domain.Patient, error) {
	pdb, err := s.GetByID(id)
	if err != nil {
		return domain.Patient{}, err
	}

	if p.Surname == "" {
		p.Surname = pdb.Surname
	}
	if p.Name == "" {
		p.Name = pdb.Name
	}
	if p.Document == "" {
		p.Document = pdb.Document
	}
	if p.CreatedAt == "" {
		p.CreatedAt = pdb.CreatedAt
	}
	p.Id = pdb.Id
	pUpdated, err := s.r.Update(id, p)
	if err != nil {
		return domain.Patient{}, err
	}
	patientUpdated, ok := pUpdated.(domain.Patient)
	if ok {
		return patientUpdated, nil
	}

	return domain.Patient{}, errors.New("failed to update the patient")
}

func (s *service) Delete(id int) error {
	return s.r.Delete(id)
}
