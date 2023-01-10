package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/handlers"
)

func main() {

	// Define command-line flags
	serverAddr := flag.String("serverAddr", "", "HTTP server network address")
	serverPort := flag.Int("serverPort", 4000, "HTTP server network port")
	mongoURI := flag.String("mongoURI", "mongodb://localhost:27017", "Database hostname url")
	mongoDBName := flag.String("mongoDBName", "sensors", "Database name")
	flag.Parse()

	// Initialize a new instance of application containing the dependencies.
	app, err := handlers.NewApplication(*mongoURI, *mongoDBName)
	if err != nil {
		panic(err)
	}
	// Initialize a new http.Server struct.
	serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	srv := &http.Server{
		Addr:         serverURI,
		ErrorLog:     app.ErrorLog(),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServe()
	app.ErrorLog().Fatal(err)
}
