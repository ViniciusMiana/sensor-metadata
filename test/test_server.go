package test

import (
	"net/http/httptest"
	"testing"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/handlers"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func startServer(t *testing.T) *httptest.Server {
	mongoURI := "mongodb://localhost:27017"
	mongoDBName := "sensors" + primitive.NewObjectID().Hex()

	// Initialize a new instance of application containing the dependencies.
	app, err := handlers.NewApplication(mongoURI, mongoDBName)
	require.NoError(t, err)

	srv := httptest.NewServer(app.Routes())
	return srv
}

func NewTestAuthWriter() runtime.ClientAuthInfoWriter {
	return &TestAuthWriter{}
}

// TestAuthWriter will implement authentication on test server
type TestAuthWriter struct {
}

// AuthenticateRequest implements the ClientAuthInfoWriter interface
func (taw TestAuthWriter) AuthenticateRequest(cr runtime.ClientRequest, reg strfmt.Registry) error {
	// TODO implement proper jwt token
	return cr.SetHeaderParam("TOKEN", "TOKEN")
}
