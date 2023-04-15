package domain

type AppointmentDTO struct {
	Appointment
	Dentist Dentist `json:"dentist" binding:"required"`
	Patient Patient `json:"patient" binding:"required"`
}
