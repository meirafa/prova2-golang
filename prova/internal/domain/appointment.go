package domain

type Appointment struct {
	Id              int    `json:"id"`
	Description     string `json:"description" binding:"required"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	IdDentist       string `json:"id_dentist" binding:"required"`
	IdPatient       string `json:"id_patient" binding:"required"`
}
