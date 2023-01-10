package test

import (
	"strings"
	"testing"

	"github.com/ViniciusMiana/sensor-metadata/test/client/client"
	"github.com/ViniciusMiana/sensor-metadata/test/client/client/sensor"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSensor(t *testing.T) {

	srv := startServer(t)
	serverAddress := strings.TrimPrefix(srv.URL, "http://")

	unauthenticatedTransport := httptransport.New(serverAddress, "/", []string{"http"})
	serviceClient := client.New(unauthenticatedTransport, strfmt.Default)
	// TODO implement authentication and closing of servers testAuthWriter := NewTestAuthWriter()
	//	_, closeServer := tKit.StartTestServer(ctx, t)
	//	defer closeServer()

	t.Run("error, not found", func(t *testing.T) {
		id := primitive.NewObjectID()
		params := sensor.NewFindByIDParams().WithID(id.Hex())

		_, err := serviceClient.Sensor.FindByID(params)
		require.EqualError(t, err, "[GET /{id}][400] findByIdBadRequest  &{Message:mongo: no documents in result}")
	})

}
