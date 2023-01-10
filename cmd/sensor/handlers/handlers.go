package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/service"
	"github.com/gorilla/mux"
)

func (app *Application) findByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	m, err := app.sensors.FindByID(ctx, id)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	app.jsonReturn(w, http.StatusOK, m)
}

func (app *Application) findByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["name"]
	m, err := app.sensors.FindByName(ctx, id)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	app.jsonReturn(w, http.StatusOK, m)
}

func (app *Application) findNearest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	lat := vars["lat"]
	lon := vars["lon"]
	m, err := app.sensors.FindNearest(ctx, lat, lon)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	app.jsonReturn(w, http.StatusOK, m)
}

func (app *Application) insert(w http.ResponseWriter, r *http.Request) {
	var sensor service.SensorMetadata
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	id, err := app.sensors.Add(ctx, sensor)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusInternalServerError)
		return
	}
	app.jsonReturn(w, http.StatusCreated, ID{ID: id})
}

func (app *Application) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	err := app.sensors.Delete(ctx, id)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	app.emptyReturn(w, http.StatusNoContent)
}

func (app *Application) update(w http.ResponseWriter, r *http.Request) {
	var sensor service.SensorMetadata
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	if id != sensor.ID {
		err = errors.New("id in url and sensor don't match")
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	err = app.sensors.Update(ctx, sensor)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusInternalServerError)
		return
	}
	app.emptyReturn(w, http.StatusNoContent)
}
