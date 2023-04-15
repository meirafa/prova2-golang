package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/meirafa/prova2-golang/cmd/server/handler"
	"github.com/meirafa/prova2-golang/internal/appointment"
	"github.com/meirafa/prova2-golang/internal/dentist"
	"github.com/meirafa/prova2-golang/internal/patient"
	"github.com/meirafa/prova2-golang/pkg/store"
)

func main() {
	// 	DB INITIALIZATION
	db, err := sql.Open("mysql", "user:password@/my_db")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStore := store.NewSQLStore()
	apStore := store.NewSQLAp()

	appRepo := appointment.NewRepository(apStore)
	appService := appointment.NewService(appRepo)
	appHandler := handler.NewAppointmentHandler(appService)

	dentistRepo := dentist.NewRepository(sqlStore)
	dentistService := dentist.NewService(dentistRepo)

	dentistHandler := handler.NewDentistHandler(dentistService)

	patientRepo := patient.NewRepository(sqlStore)
	patientService := patient.NewService(patientRepo)
	patientHandler := handler.NewPatientHandler(patientService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	api := r.Group("/api/")
	{
		appointments := api.Group("/appointments")
		{
			appointments.GET("", appHandler.GetAll())
			appointments.GET(":id", appHandler.GetByID())
			appointments.GET("/patient/:document", appHandler.GetByDocumentPatient())

			appointments.POST("", appHandler.Post())
			appointments.PUT(":id", appHandler.Put())
			appointments.PATCH(":id", appHandler.Patch())
			appointments.DELETE(":id", appHandler.Delete())
		}
		dentists := api.Group("/dentists")
		{
			dentists.GET("", dentistHandler.GetAll())
			dentists.GET(":id", dentistHandler.GetByID())

			dentists.POST("", dentistHandler.Post())
			dentists.PUT(":id", dentistHandler.Put())
			dentists.PATCH(":id", dentistHandler.Patch())
			dentists.DELETE(":id", dentistHandler.Delete())
		}
		patients := api.Group("/patients")
		{
			patients.GET("", patientHandler.GetAll())
			patients.GET(":id", patientHandler.GetByID())

			patients.POST("", patientHandler.Post())
			patients.PUT(":id", patientHandler.Put())
			patients.PATCH(":id", patientHandler.Patch())
			patients.DELETE(":id", patientHandler.Delete())
		}
	}

	r.Run(":8083")
}
