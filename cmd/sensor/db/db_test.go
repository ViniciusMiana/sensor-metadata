package db

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddFindByIDSensor(t *testing.T) {
	var s SensorStore
	var err error
	inserted := Sensor{
		Name: "Sensor 1",
		Tags: []string{"Tag1", "Tag2"},
		Location: &Location{
			Lat: 55,
			Lon: 44,
		},
	}
	s, err = NewSensorStore(`mongodb://localhost:27017`, "sensors"+primitive.NewObjectID().Hex())
	require.NoError(t, err)
	defer func() {
		_, err = s.(*sensorStore).sensors.DeleteMany(context.Background(), bson.D{})
		require.NoError(t, err)
	}()
	ctx := context.Background()
	id, err := s.Add(ctx, inserted)
	require.NoError(t, err)
	require.NotNil(t, id)
	require.NotEqual(t, primitive.NilObjectID, id)
	sensor, err := s.FindByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, inserted.Name, sensor.Name)
	require.Equal(t, inserted.Tags, sensor.Tags)
	require.Equal(t, inserted.Location, sensor.Location)
	require.Equal(t, inserted.Location.Lat, sensor.GeoJson.Coordinates[1])
	require.Equal(t, inserted.Location.Lon, sensor.GeoJson.Coordinates[0])
}

func TestAddUpdateFindByNameSensor(t *testing.T) {
	var s SensorStore
	var err error
	inserted := Sensor{
		Name: "Sensor 1",
		Tags: []string{"Tag1", "Tag2"},
		Location: &Location{
			Lat: 55,
			Lon: 44,
		},
	}
	s, err = NewSensorStore(`mongodb://localhost:27017`, "sensors"+primitive.NewObjectID().Hex())
	require.NoError(t, err)
	defer func() {
		_, err = s.(*sensorStore).sensors.DeleteMany(context.Background(), bson.D{})
		require.NoError(t, err)
	}()
	ctx := context.Background()
	id, err := s.Add(ctx, inserted)
	require.NoError(t, err)
	require.NotNil(t, id)
	require.NotEqual(t, primitive.NilObjectID, id)
	sensor, err := s.FindByID(ctx, id)
	require.NoError(t, err)
	sensor.Tags = []string{"Tag3", "Tag4"}
	sensor.Name = "New Name"
	sensor.Location.Lon = 123
	sensor.Location.Lat = 3.14
	err = s.Update(ctx, *sensor)
	require.NoError(t, err)
	updated, err := s.FindByName(ctx, "New Name")
	require.NoError(t, err)
	require.Equal(t, sensor.Name, updated.Name)
	require.Equal(t, sensor.Tags, updated.Tags)
	require.Equal(t, sensor.Location, updated.Location)
	require.Equal(t, sensor.Location.Lat, updated.GeoJson.Coordinates[1])
	require.Equal(t, sensor.Location.Lon, updated.GeoJson.Coordinates[0])
	require.Equal(t, sensor.ID, updated.ID)
}

func TestAddFindDeleteSensor(t *testing.T) {
	var s SensorStore
	var err error
	inserted := Sensor{
		Name: "Sensor 1",
		Tags: []string{"Tag1", "Tag2"},
		Location: &Location{
			Lat: 55,
			Lon: 44,
		},
	}
	s, err = NewSensorStore(`mongodb://localhost:27017`, "sensors"+primitive.NewObjectID().Hex())
	require.NoError(t, err)
	defer func() {
		_, err = s.(*sensorStore).sensors.DeleteMany(context.Background(), bson.D{})
		require.NoError(t, err)
	}()
	ctx := context.Background()
	id, err := s.Add(ctx, inserted)
	require.NoError(t, err)
	require.NotNil(t, id)
	require.NotEqual(t, primitive.NilObjectID, id)
	_, err = s.FindByID(ctx, id)
	require.NoError(t, err)
	err = s.Delete(ctx, id)
	require.NoError(t, err)
	_, err = s.FindByID(ctx, id)
	require.Error(t, err)
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
}

func TestFindNearest(t *testing.T) {
	var s SensorStore
	var err error
	sensors := []Sensor{
		{
			Name: "Sensor Washington",
			Tags: []string{"Tag1", "Tag2"},
			Location: &Location{
				Lat: 38.9072,
				Lon: -77.0369,
			},
		},
		{
			Name: "Sensor NY",
			Tags: []string{"Tag1", "Tag2"},
			Location: &Location{
				Lat: 40.7128,
				Lon: -74.0060,
			},
		},
		{
			Name: "Sensor Atlanta",
			Tags: []string{"Tag1", "Tag2"},
			Location: &Location{
				Lat: 33.7488,
				Lon: -84.3877,
			},
		},
	}
	s, err = NewSensorStore(`mongodb://localhost:27017`, "sensors"+primitive.NewObjectID().Hex())
	require.NoError(t, err)
	defer func() {
		_, err = s.(*sensorStore).sensors.DeleteMany(context.Background(), bson.D{})
		require.NoError(t, err)
	}()
	ctx := context.Background()
	for i := range sensors {
		sensor := sensors[i]
		_, err = s.Add(ctx, sensor)
		require.NoError(t, err)
	}
	sensor, err := s.FindNearest(ctx, Location{Lat: 34, Lon: 45})
	require.NoError(t, err)
	require.Equal(t, sensors[1].Name, sensor.Name)
	sensor, err = s.FindNearest(ctx, Location{Lat: 34, Lon: -74})
	require.NoError(t, err)
	require.Equal(t, sensors[0].Name, sensor.Name)
	sensor, err = s.FindNearest(ctx, Location{Lat: 36.1627, Lon: -86.7816})
	require.NoError(t, err)
	require.Equal(t, sensors[2].Name, sensor.Name)

}
