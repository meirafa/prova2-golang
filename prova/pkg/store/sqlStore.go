package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/meirafa/prova2-golang/internal/domain"
)

var (
	AP = "appointments"
	DE = "dentists"
	PE = "patients"
)

// NewSQLStore estabelece conexão com a interface
func NewSQLStore() Store {
	database, err := sql.Open("mysql", "user:password@/my_db")
	if err != nil {
		panic(err)
	}
	return &sqlStore{
		db: database,
	}
}

type sqlStore struct {
	db *sql.DB
}

// GetAll retorna todas as consultas da tabela selecionada
func (s *sqlStore) GetAll(tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxGetAllByTable(tableName, s)
	case DE:
		return auxGetAllByTable(tableName, s)
	case PE:
		return auxGetAllByTable(tableName, s)
	default:
		return nil, errors.New("an error occurred while trying to get data from db")
	}
}

// GetByID retorna as consultas da tabela selecionada por id
func (s *sqlStore) GetByID(id int, tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxGetByIDByTable(tableName, id, s)
	case DE:
		return auxGetByIDByTable(tableName, id, s)
	case PE:
		return auxGetByIDByTable(tableName, id, s)
	default:
		return nil, errors.New("an error occurred while trying to get data from db")
	}
}

// Save insere os novos dados na tabela selecionada
func (s *sqlStore) Save(entity interface{}, tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxSave(tableName, s, entity)
	case DE:
		return auxSave(tableName, s, entity)
	case PE:
		return auxSave(tableName, s, entity)
	default:
		return nil, errors.New("failed to start inserting your data")
	}
}

// Update atualiza as linhas da tabela seleconada
func (s *sqlStore) Update(entityID int, entity interface{}, tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxUpdate(tableName, s, entity, entityID)
	case DE:
		return auxUpdate(tableName, s, entity, entityID)
	case PE:
		return auxUpdate(tableName, s, entity, entityID)
	default:
		return nil, errors.New("failed to update entity data")
	}
}

// Delete exclui uma linha da tabela selecionada por ID
func (s *sqlStore) Delete(entityID int, tableName string) error {
	switch tableName {
	case AP:
		return auxDelete(tableName, s, entityID)
	case DE:
		return auxDelete(tableName, s, entityID)
	case PE:
		return auxDelete(tableName, s, entityID)
	default:
		return errors.New("failed to delete")
	}
}

// auxGetAllByTable - Função chamada por GetAll, aqui a tabela selecionada é validada e todas as consultas de seleção são feitas.
func auxGetAllByTable(tableName string, s *sqlStore) (interface{}, error) {
	var entities []struct{}

	switch tableName {
	case AP:
		Query := "SELECT a.id, a.description, DATE_FORMAT(a.appointment_date,'%d/%m/%Y %H:%i') appointment_date,a.id_dentist,a.id_patient,d.id,d.surname,d.name,d.registration,p.id,p.surname,p.name,p.document,DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') created_at FROM appointments a INNER JOIN dentists d on a.id_dentist = d.registration INNER JOIN patients p on a.id_patient = p.document ORDER BY a.appointment_date"
		rows, err := s.db.Query(Query)
		if err != nil {
			return entities, err
		}

		var appointment domain.AppointmentDTO
		var appointments []domain.AppointmentDTO

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
	case DE:
		rows, err := s.db.Query("SELECT * FROM dentists")
		if err != nil {
			return entities, err
		}

		var dentist domain.Dentist
		var dentists []domain.Dentist

		for rows.Next() {
			if err := rows.Scan(
				&dentist.Id,
				&dentist.Surname,
				&dentist.Name,
				&dentist.Registration); err != nil {
				return dentists, err
			}
			dentists = append(dentists, dentist)
		}
		return dentists, nil
	case PE:
		rows, err := s.db.Query("SELECT p.id, p.surname,p.name,p.document, DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') FROM patients p")
		if err != nil {
			return entities, err
		}

		var patient domain.Patient
		var patients []domain.Patient
		for rows.Next() {
			if err := rows.Scan(
				&patient.Id,
				&patient.Surname,
				&patient.Name,
				&patient.Document,
				&patient.CreatedAt); err != nil {
				return patients, err
			}
			patients = append(patients, patient)
		}
		return patients, nil
	default:
		return nil, errors.New("failed to load data from db at sqlStore.go file")
	}
}

// auxGetByIDByTable - Função chamada por GetByID, aqui a tabela selecionada é validada e todas as seleções * de *table_name* onde id = *entity_id* são feitas.
func auxGetByIDByTable(tableName string, entityID int, s *sqlStore) (interface{}, error) {
	var entity struct{}

	switch tableName {
	case AP:
		query := "SELECT a.id, a.description, DATE_FORMAT(a.appointment_date,'%d/%m/%Y %H:%i') appointment_date,a.id_dentist,a.id_patient,d.id,d.surname,d.name,d.registration,p.id,p.surname,p.name,p.document,DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') created_at FROM appointments a INNER JOIN dentists d on a.id_dentist = d.registration INNER JOIN patients p on a.id_patient = p.document WHERE a.id = ? ORDER BY a.appointment_date"
		rows, err := s.db.Query(query, entityID)
		if err != nil {
			return entity, err
		}
		defer rows.Close()

		var appointment domain.AppointmentDTO
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
				return appointment, err
			}
			return appointment, nil
		}
		if rows.Next() {
			return appointment, nil
		}
		return nil, err
	case DE:
		rows, err := s.db.Query("SELECT * FROM dentists WHERE id = ?", entityID)
		if err != nil {
			return entity, err
		}
		defer rows.Close()

		var dentist domain.Dentist
		for rows.Next() {
			if err = rows.Scan(
				&dentist.Id,
				&dentist.Surname,
				&dentist.Name,
				&dentist.Registration); err != nil {
				return dentist, err
			}
			return dentist, nil
		}
		if rows.Next() {
			return dentist, nil
		}
		return nil, err
	case PE:
		rows, err := s.db.Query("SELECT p.id, p.surname,p.name,p.document, DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') FROM patients p WHERE id = ?", entityID)
		if err != nil {
			return entity, err
		}
		defer rows.Close()

		var patient domain.Patient
		for rows.Next() {
			if err = rows.Scan(
				&patient.Id,
				&patient.Surname,
				&patient.Name,
				&patient.Document,
				&patient.CreatedAt); err != nil {
				return nil, err
			}
			return patient, nil
		}
		if rows.Next() {
			return patient, nil
		}
		return nil, err
	default:
		return nil, errors.New("failed to get by id from db")
	}
}

// auxSave - Função chamada por Save, aqui as inserções são feitas na tabela selecionada.
func auxSave(tableName string, s *sqlStore, entity interface{}) (interface{}, error) {
	switch tableName {
	case AP:
		log.Println("... inserting data into appointments table.")
		var appointment domain.Appointment
		appointment, ok := entity.(domain.Appointment)
		if ok {
			apAppointmentDateParsed, err := time.Parse("02/01/2006 15:04", appointment.AppointmentDate)
			if err != nil {
				log.Println("failed to convert datetimee")
				return nil, errors.New("failed to convert datetimee")
			}
			//
			log.Println(apAppointmentDateParsed.String())
			result, err := s.db.Exec("INSERT INTO appointments(DESCRIPTION, appointment_date, id_dentist, id_patient) VALUES(?,?,?,?)",
				appointment.Description,
				apAppointmentDateParsed,
				appointment.IdDentist,
				appointment.IdPatient)
			if err != nil {
				fmt.Println("inserting data failed :", err.Error())
				return nil, err
			}
			lastInsertedID, err := result.LastInsertId()
			if err != nil {
				fmt.Println("error trying to get id inserted:", err.Error())
				return nil, err
			}
			appointment.Id = int(lastInsertedID)
			log.Println("... INSERT operation was successfully")
			return s.GetByID(appointment.Id, AP)
		}
	case DE:
		var dentist domain.Dentist
		dentist, ok := entity.(domain.Dentist)
		if ok {
			result, err := s.db.Exec("INSERT INTO dentists(surname, name, registration) VALUES (?,?,?)",
				dentist.Surname,
				dentist.Name,
				dentist.Registration)
			if err != nil {
				fmt.Println("inserting data failed :", err.Error())
				return nil, err
			}
			lastInsertedID, err := result.LastInsertId()
			if err != nil {
				fmt.Println("error trying to get id inserted:", err.Error())
				return nil, err
			}
			dentist.Id = int(lastInsertedID)
			fmt.Println("dentist inserted at db:", dentist)
			return dentist, nil
		}
	case PE:
		var patient domain.Patient
		patient, ok := entity.(domain.Patient)
		if ok {
			patCreatedAtParsed, err := time.Parse("02/01/2006 15:04:05", patient.CreatedAt)
			if err != nil {
				return nil, errors.New("failed to convert patient created_at field")
			}
			result, err := s.db.Exec("INSERT INTO patients(surname, name, document, created_at) VALUES (?,?,?,?)",
				patient.Surname,
				patient.Name,
				patient.Document,
				patCreatedAtParsed)
			if err != nil {
				fmt.Println("inserting data failed :", err.Error())
				return nil, err
			}
			lastInsertedID, err := result.LastInsertId()
			if err != nil {
				fmt.Println("error trying to get id inserted:", err.Error())
				return nil, err
			}
			patient.Id = int(lastInsertedID)
			return patient, nil
		}
	default:
		return nil, errors.New("failed to insert data at database")
	}
	return nil, errors.New("failed to insert data at database")
}

// auxUpdate - Função chamada por Update, aqui as atualizações são feitas na tabela selecionada.
func auxUpdate(tableName string, s *sqlStore, entity interface{}, entityId int) (interface{}, error) {
	switch tableName {
	case AP:
		var appointment domain.Appointment
		appointment, ok := entity.(domain.Appointment)
		if ok {
			apAppointmentDateParsed, err := time.Parse("02/01/2006 15:04", appointment.AppointmentDate)
			if err != nil {
				log.Println(err.Error(), "\nDate parsed: ", apAppointmentDateParsed)
				return nil, errors.New("failed to convert datetime")
			}
			_, err = s.db.Exec("UPDATE appointments SET description = ?, appointment_date = ?, id_dentist = ?, id_patient = ? WHERE id = ?",
				appointment.Description,
				apAppointmentDateParsed,
				appointment.IdDentist,
				appointment.IdPatient,
				entityId)
			if err != nil {
				return nil, err
			}
			return s.GetByID(entityId, AP)
		}
	case DE:
		var dentist domain.Dentist
		dentist, ok := entity.(domain.Dentist)
		if ok {
			_, err := s.db.Exec("UPDATE dentists SET surname = ?, name = ?, registration = ? WHERE id = ?",
				dentist.Surname,
				dentist.Name,
				dentist.Registration,
				entityId)
			if err != nil {
				return nil, err
			}
			return dentist, nil
		}
	case PE:
		var patient domain.Patient
		patient, ok := entity.(domain.Patient)
		if ok {
			paCreatedAtParsed, err := time.Parse("02/01/2006 15:04", patient.CreatedAt)
			if err != nil {
				return nil, errors.New("failed to convert patient created_at field: " + patient.CreatedAt)
			}

			_, err = s.db.Exec("UPDATE patients SET surname = ?, name = ?, document = ?, created_at = ? WHERE id = ?",
				patient.Surname,
				patient.Name,
				patient.Document,
				paCreatedAtParsed,
				entityId)
			if err != nil {
				return nil, err
			}
			return patient, nil
		}
	default:
		return nil, errors.New("failed to update data into database")
	}
	return nil, errors.New("failed to update data into database")
}

// auxDelete - Função chamada por Delete, aqui as deleções são feitas na tabela selecionada.
func auxDelete(tableName string, s *sqlStore, entityID int) error {
	switch tableName {
	case AP:
		result, err := s.db.Exec("DELETE FROM appointments WHERE id =?", entityID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if count == 0 {
			return errors.New("entity not found at database")
		}
		return nil
	case DE:
		result, err := s.db.Exec("DELETE FROM dentists WHERE id =?", entityID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if count == 0 {
			return errors.New("entity not found at database")
		}
		return nil
	case PE:
		result, err := s.db.Exec("DELETE FROM patients WHERE id =?", entityID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if count == 0 {
			return errors.New("entity not found at database")
		}
		return nil
	default:
		return errors.New("failed to delete row")
	}
}
