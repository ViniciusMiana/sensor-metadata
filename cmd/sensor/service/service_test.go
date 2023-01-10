package service

import (
	"context"
	"testing"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/db"
	dbMock "github.com/ViniciusMiana/sensor-metadata/mocks/sensor/db"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFindByID(t *testing.T) {
	mockSensor := dbMock.NewSensorStore(t)
	ctx := context.Background()
	var err error
	sensor := db.Sensor{
		Name: "Sensor 1",
		Tags: []string{"Tag1", "Tag2"},
		Location: &db.Location{
			Lat: 55,
			Lon: 44,
		},
	}
	service := sensorMetadataService{
		sensorStore: mockSensor,
	}
	mockSensor.On("FindByID", ctx, mock.Anything).Return(&sensor, nil).Once()
	defer mockSensor.AssertExpectations(t)
	result, err := service.FindByID(ctx, primitive.NewObjectID().Hex())
	require.NoError(t, err)
	dbResult, err := result.ToDatabase()
	require.NoError(t, err)
	require.Equal(t, sensor, *dbResult)
}

func TestFindByName(t *testing.T) {
	mockSensor := dbMock.NewSensorStore(t)
	ctx := context.Background()
	var err error
	sensor := db.Sensor{
		Name: "Sensor 1",
		Tags: []string{"Tag1", "Tag2"},
		Location: &db.Location{
			Lat: 55,
			Lon: 44,
		},
	}
	service := sensorMetadataService{
		sensorStore: mockSensor,
	}
	mockSensor.On("FindByName", ctx, sensor.Name).Return(&sensor, nil).Once()
	defer mockSensor.AssertExpectations(t)
	result, err := service.FindByName(ctx, "Sensor 1")
	require.NoError(t, err)
	dbResult, err := result.ToDatabase()
	require.NoError(t, err)
	require.Equal(t, sensor, *dbResult)
}

func TestFindNearest(t *testing.T) {
	mockSensor := dbMock.NewSensorStore(t)
	ctx := context.Background()
	var err error
	sensor := db.Sensor{
		Name: "Sensor 1",
		Tags: []string{"Tag1", "Tag2"},
		Location: &db.Location{
			Lat: 55,
			Lon: 44,
		},
	}
	service := sensorMetadataService{
		sensorStore: mockSensor,
	}
	mockSensor.On("FindNearest", ctx, db.Location{
		Lat: 1,
		Lon: 2,
	}).Return(&sensor, nil).Once()
	defer mockSensor.AssertExpectations(t)
	result, err := service.FindNearest(ctx, "1", "2")
	require.NoError(t, err)
	dbResult, err := result.ToDatabase()
	require.NoError(t, err)
	require.Equal(t, sensor, *dbResult)
}
