package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ViniciusMiana/sensor-metadata/cmd/authenticator/db"
)

// TODO remove dependency from DB having a separate DTO

func (app *Application) register(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	err = app.service.Register(user)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusInternalServerError)
		return
	}
	app.emptyReturn(w, http.StatusCreated)
}

func (app *Application) login(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusBadRequest)
		return
	}
	token, err := app.service.Login(user)
	if err != nil {
		app.jsonErrorReturn(w, err, http.StatusUnauthorized)
		return
	}
	app.jsonReturn(w, http.StatusOK, token)
}
