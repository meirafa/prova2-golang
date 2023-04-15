package store

import (
	"database/sql"
	"log"

	"github.com/meirafa/prova2-golang/internal/domain"
)

// ApStore - Define o contrato para ApStore que é composto de uma interface Store.
type ApStore interface {
	Store
	GetAllAppointmentsByPatientIdentify(identifyNumber string) ([]domain.AppointmentDTO, error)
	GetAllAppointmentsByDentistsLicense(Registration string) ([]domain.AppointmentDTO, error)
	GetAllAppointmentsByDateTimeInterval(startDateTime, endDateTime string) ([]domain.Appointment, error)
}

// NewSQLAp - Inicializa interface ApStore
func NewSQLAp() ApStore {
	database, err := sql.Open("mysql", "user:password@/my_db")
	if err != nil {
		panic(err)
	}
	return &appointmentStore{
		sqlStore: &sqlStore{db: database},
		db:       database,
	}
}

type appointmentStore struct {
	*sqlStore
	db *sql.DB
}

// GetAllAppointmentsByPatientIdentify - retorna uma lista de todas as consultas feitas por um paciente através do seu número de identidade
func (sa *appointmentStore) GetAllAppointmentsByPatientIdentify(identifyNumber string) ([]domain.AppointmentDTO, error) {
	var appointment domain.AppointmentDTO
	var appointments []domain.AppointmentDTO

	query := "SELECT a.id, a.description, DATE_FORMAT(a.appointment_date,'%d/%m/%Y %H:%i') appointment_date,a.id_dentist,a.id_patient,d.id,d.surname,d.name,d.registration,p.id,p.surname,p.name,p.document,DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') created_at FROM appointments a INNER JOIN dentists d on a.id_dentist = d.registration INNER JOIN patients p on a.id_patient = p.document WHERE a.id_patient = ? ORDER BY a.appointment_date"
	rows, err := sa.db.Query(query, identifyNumber)
	if err != nil {
		return appointments, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&appointment.Id,
			&appointment.Description,
			&appointment.AppointmentDate,
			&appointment.IdDentist,
			&appointment.IdPatient,
			&appointment.Dentist.Id,
			&appointment.Dentist.Surname,
			&appointment.Dentist.Name,
			&appointment.Dentist.Registration,
			&appointment.Patient.Id,
			&appointment.Patient.Surname,
			&appointment.Patient.Name,
			&appointment.Patient.Document,
			&appointment.Patient.CreatedAt); err != nil {
			return appointments, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

// GetAllAppointmentsByDentistsLicense - retorna uma lista de todas as consultas feitas por um dentista através do seu número de licença
func (sa *appointmentStore) GetAllAppointmentsByDentistsLicense(Registration string) ([]domain.AppointmentDTO, error) {
	var appointment domain.AppointmentDTO
	var appointments []domain.AppointmentDTO

	query := "SELECT a.id, a.description, DATE_FORMAT(a.appointment_date,'%d/%m/%Y %H:%i') appointment_date,a.id_dentist,a.id_patient,d.id,d.surname,d.name,d.registration,p.id,p.surname,p.name,p.document,DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') created_at FROM appointments a INNER JOIN dentists d on a.id_dentist = d.registration INNER JOIN patients p on a.id_patient = p.document WHERE a.id_dentist = ? ORDER BY a.appointment_date"
	rows, err := sa.db.Query(query, Registration)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		if err := rows.Scan(
			&appointment.Id,
			&appointment.Description,
			&appointment.AppointmentDate,
			&appointment.IdDentist,
			&appointment.IdPatient,
			&appointment.Dentist.Id,
			&appointment.Dentist.Surname,
			&appointment.Dentist.Name,
			&appointment.Dentist.Registration,
			&appointment.Patient.Id,
			&appointment.Patient.Surname,
			&appointment.Patient.Name,
			&appointment.Patient.Document,
			&appointment.Patient.CreatedAt); err != nil {
			return appointments, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

// GetAllAppointmentsByDateTimeInterval - retorna uma lista de todos os compromissos durante um intervalo de data e hora. Usado principalmente para validar se uma data está disponível.
func (sa *appointmentStore) GetAllAppointmentsByDateTimeInterval(startDateTime, endDateTime string) ([]domain.Appointment, error) {
	var appointment domain.Appointment
	var appointments []domain.Appointment
	rows, err := sa.db.Query("SELECT * FROM appointments WHERE appointment_date BETWEEN ? AND ?", startDateTime, endDateTime)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		if err := rows.Scan(
			&appointment.Id,
			&appointment.Description,
			&appointment.AppointmentDate,
			&appointment.IdDentist,
			&appointment.IdPatient); err != nil {
			return appointments, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}
