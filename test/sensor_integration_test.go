package test

import (
	"strings"
	"testing"

	"github.com/ViniciusMiana/sensor-metadata/test/client/models"

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

	t.Run("error, not found", func(t *testing.T) {
		id := primitive.NewObjectID()
		params := sensor.NewFindByIDParams().WithID(id.Hex())

		_, err := serviceClient.Sensor.FindByID(params)
		require.EqualError(t, err, "[GET /{id}][400] findByIdBadRequest  &{Message:mongo: no documents in result}")
	})
	t.Run("auth tests", func(t *testing.T) {
		params := sensor.NewCreateSensorParams().WithSensorMetaDataCreationRequest(
			&models.SensorMetadata{
				Name: "XXX",
			},
		)
		_, err := serviceClient.Sensor.CreateSensor(params, NewTestAuthWriter("USER", ""))
		require.Error(t, err)
		require.EqualError(t, err, "[POST /][403] createSensorForbidden  &{Message:This user can't perform this function}")
		_, err = serviceClient.Sensor.CreateSensor(params, NewTestAuthWriter("ADMIN", ""))
		require.NoError(t, err)
		_, err = serviceClient.Sensor.CreateSensor(params, NewTestAuthWriter("", "Something wrong"))
		require.Error(t, err)
		require.EqualError(t, err, "[POST /][400] createSensorBadRequest  &{Message:Token is invalid}")

	})

}
