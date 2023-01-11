package test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/exp/errors/fmt"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/handlers"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const nilString = "NIL"

func startServer(t *testing.T) *httptest.Server {
	mongoURI := "mongodb://localhost:27017"
	mongoDBName := "sensors" + primitive.NewObjectID().Hex()

	// Initialize a new instance of application containing the dependencies.
	app, err := handlers.NewApplication(mongoURI, mongoDBName)
	require.NoError(t, err)
	app.ParseToken = ParseTestToken
	srv := httptest.NewServer(app.Routes())
	return srv
}

// NewTestAuthWriter creates TestAuthWriter that injects the passed role and error in the token
func NewTestAuthWriter(role, err string) runtime.ClientAuthInfoWriter {
	taw := TestAuthWriter{}
	if err == "" {
		taw.err = nilString
	} else {
		taw.err = err
	}
	if role == "" {
		taw.role = nilString
	} else {
		taw.role = role
	}
	return taw
}

// TestAuthWriter will implement authentication on test server
type TestAuthWriter struct {
	role string
	err  string
}

// ParseTestToken parses the test token it should contain two strings role and error or NIL.
func ParseTestToken(token string, _ []byte) (*handlers.TokenClaims, error) {
	fields := strings.Fields(token)
	role := ""
	var err error
	if fields[0] != nilString {
		role = fields[0]
	}
	if fields[1] != nilString {
		err = errors.New(fields[1])
	}
	return &handlers.TokenClaims{
		UserName: "TEST",
		Role:     role,
	}, err
}

// AuthenticateRequest implements the ClientAuthInfoWriter interface
func (taw TestAuthWriter) AuthenticateRequest(cr runtime.ClientRequest, reg strfmt.Registry) error {
	return cr.SetHeaderParam("Authorization", fmt.Sprintf("token %s %s", taw.role, taw.err))
}
