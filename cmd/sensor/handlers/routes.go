package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/service"
	"github.com/gorilla/mux"
)

type Application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	sensors    service.SensorMetadataService
	jwtPubKey  []byte
	ParseToken func(token string, pubKey []byte) (*TokenClaims, error)
}

func (app Application) ErrorLog() *log.Logger {
	return app.errorLog
}

// TODO move to a common pkg folder
// TokenClaims contains a basic jwt claim structure
type TokenClaims struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
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
		errorLog:   errLog,
		infoLog:    infoLog,
		sensors:    srv,
		jwtPubKey:  []byte(os.Getenv("tls.crt")),
		ParseToken: ParseJWTToken,
	}, nil

}

func (app *Application) requireAuthentication(fn http.HandlerFunc, allowedRoles []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		var authToken string

		if strings.HasPrefix(strings.ToLower(authHeader), "token ") {
			authToken = authHeader[len("token "):]
		} else {
			query := r.URL.Query()
			authToken = query.Get("token")
		}

		if authToken == "" {
			app.jsonErrorReturn(w, errors.New("Token is required"), http.StatusUnauthorized)
			return
		}
		claims, err := app.ParseToken(authToken, app.jwtPubKey)
		if err != nil {
			app.jsonErrorReturn(w, errors.New("Token is invalid"), http.StatusBadRequest)
			return
		}
		if len(allowedRoles) > 0 && !slices.Contains(allowedRoles, claims.Role) {
			app.jsonErrorReturn(w, errors.New("This user can't perform this function"), http.StatusForbidden)
			return
		}
		fn(w, r)
	}
}

func ParseJWTToken(token string, pubKey []byte) (*TokenClaims, error) {
	jwtToken, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		a, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
		return a, err
	})
	if err != nil {
		return nil, err
	}
	if jwtToken.Valid {
		var tokenClaims *TokenClaims
		err = mapstructure.Decode(jwtToken.Claims, &tokenClaims)
		if err != nil {
			return nil, fmt.Errorf("could not parse token claims: %v", err)
		}
		return tokenClaims, nil
	}
	return nil, errors.New("token has expired")
}

func (app *Application) Routes() *mux.Router {
	// Register handler functions.
	r := mux.NewRouter()
	r.HandleFunc("/nearest/{lat}/{lon}", app.findNearest).Methods(http.MethodGet)
	r.HandleFunc("/{id}", app.findByID).Methods(http.MethodGet)
	r.HandleFunc("/", app.requireAuthentication(app.insert, []string{"ADMIN"})).Methods(http.MethodPost)
	r.HandleFunc("/sensor", app.insertWithLocationName).Methods(http.MethodPost)
	r.HandleFunc("/{id}", app.requireAuthentication(app.delete, []string{"ADMIN"})).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", app.requireAuthentication(app.update, []string{"ADMIN"})).Methods(http.MethodPut)
	r.HandleFunc("/by-name/{name}", app.findByName).Methods(http.MethodGet)
	r.HandleFunc("/nearest-by-name/{location}", app.findNearestByLocatioName).Methods(http.MethodGet)

	return r
}
