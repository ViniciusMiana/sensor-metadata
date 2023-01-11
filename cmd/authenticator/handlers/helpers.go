package handlers

import (
	"encoding/json"
	"net/http"
)

// Juca is a structure to return errors in json format
type Error struct {
	Message string `json:"message"`
}

func (app Application) jsonErrorReturn(w http.ResponseWriter, err error, httpStatus int) {
	res := Error{
		Message: err.Error(),
	}
	app.jsonReturn(w, httpStatus, res)
}

func (app Application) jsonReturn(w http.ResponseWriter, statusCode int, jsonObject interface{}) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(jsonObject)
	if err != nil {
		app.errorLog.Printf("could not encode json return: %s %+v", err.Error(), jsonObject)
	}
	app.infoLog.Printf("return json response %+v", jsonObject)
}

func (app Application) emptyReturn(w http.ResponseWriter, statusCode int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	app.infoLog.Printf("return empty %d", statusCode)
}
