package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/service"
	"github.com/gorilla/mux"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	sensors  service.SensorMetadataService
}

func (app Application) ErrorLog() *log.Logger {
	return app.errorLog
}

func NewApplication(uri, databaseName string) (*Application, error) {
	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.Println("Starting application")
	srv, err := service.NewSensorMetadataService(uri, databaseName)
	if err != nil {
		return nil, err
	}
	return &Application{
		errorLog: errLog,
		infoLog:  infoLog,
		sensors:  srv,
	}, nil

}

func (app *Application) Routes() *mux.Router {
	// Register handler functions.
	r := mux.NewRouter()
	r.HandleFunc("/nearest/{lat}/{lon}", app.findNearest).Methods(http.MethodGet)
	r.HandleFunc("/{id}", app.findByID).Methods(http.MethodGet)
	r.HandleFunc("/", app.insert).Methods(http.MethodPost)
	r.HandleFunc("/{id}", app.delete).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", app.update).Methods(http.MethodPut)
	r.HandleFunc("/by-name/{name}", app.findByName).Methods(http.MethodGet)

	return r
}
