package service

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ViniciusMiana/sensor-metadata/cmd/sensor/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Points for improvement: */
// TODO 1 - Do we really need this layer?
// TODO 2 - Separate data objects and service in different files
// TODO 3 - Structure errors
// TODO 4 - Increase test coverage

// Location represents a location with lat and lon
type Location struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// SensorMetadata represents a sensor metadata DTO
type SensorMetadata struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name"`
	Location *Location `json:"location,omitempty"`
	Tags     []string  `json:"tags"`
}

// SensorMetadataWithLocationName represents a sensor metadata DTO
type SensorMetadataWithLocationName struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name"`
	Location string   `json:"location,omitempty"`
	Tags     []string `json:"tags"`
}

// ToDatabase converts sensor meta-data to the database format
func (s SensorMetadata) ToDatabase() (*db.Sensor, error) {
	mObj := db.Sensor{
		Name:     s.Name,
		Tags:     s.Tags,
		Location: nil,
	}
	if s.Location != nil {
		lat, err := strconv.ParseFloat(s.Location.Lat, 64)
		if err != nil {
			return nil, err
		}
		lon, err := strconv.ParseFloat(s.Location.Lon, 64)
		if err != nil {
			return nil, err
		}

		mObj.Location = &db.Location{
			Lat: lat,
			Lon: lon,
		}
	}
	if s.ID != "" {
		oid, err := primitive.ObjectIDFromHex(s.ID)
		if err != nil {
			return nil, err
		}
		mObj.ID = oid
	}
	return &mObj, nil
}

// FromDatabaseToSensorMetadata converts the mongo datq structure to the DTO
func FromDatabaseToSensorMetadata(mobj db.Sensor) *SensorMetadata {
	sensor := SensorMetadata{
		Name: mobj.Name,
		Tags: mobj.Tags,
	}
	if mobj.ID != primitive.NilObjectID {
		sensor.ID = mobj.ID.Hex()
	}
	if mobj.Location != nil {
		sensor.Location = &Location{
			Lat: fmt.Sprintf("%f", mobj.Location.Lat),
			Lon: fmt.Sprintf("%f", mobj.Location.Lon),
		}
	}
	return &sensor
}

// SensorMetadataService is the interface to the provided services
type SensorMetadataService interface {
	FindByName(ctx context.Context, name string) (sensor *SensorMetadata, err error)
	FindByID(ctx context.Context, id string) (sensor *SensorMetadata, err error)
	Add(ctx context.Context, sensor SensorMetadata) (id string, err error)
	AddWithLocationName(ctx context.Context, sensor SensorMetadataWithLocationName) (id string, err error)
	Update(ctx context.Context, sensor SensorMetadata) (err error)
	Delete(ctx context.Context, id string) (err error)
	FindNearest(ctx context.Context, lat, lon string) (sensor *SensorMetadata, err error)
	FindNearestByLocatioName(ctx context.Context, location string) (sensor *SensorMetadata, err error)
}

type sensorMetadataService struct {
	sensorStore db.SensorStore
	mapBox      MapBox
}

func NewSensorMetadataService(uri, databaseName string) (*sensorMetadataService, error) {
	apiKey := os.Getenv("API_KEY")
	ss, err := db.NewSensorStore(uri, databaseName)
	if err != nil {
		return nil, err
	}
	return &sensorMetadataService{
		sensorStore: ss,
		mapBox:      NewMapBox(apiKey),
	}, nil
}

func (s sensorMetadataService) FindByName(ctx context.Context, name string) (sensor *SensorMetadata, err error) {
	sensorMongo, err := s.sensorStore.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return FromDatabaseToSensorMetadata(*sensorMongo), nil
}

func (s sensorMetadataService) FindByID(ctx context.Context, id string) (sensor *SensorMetadata, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	sensorMongo, err := s.sensorStore.FindByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	return FromDatabaseToSensorMetadata(*sensorMongo), nil
}

func (s sensorMetadataService) Add(ctx context.Context, sensor SensorMetadata) (id string, err error) {
	sensorMongo, err := sensor.ToDatabase()
	if err != nil {
		return "", err
	}
	oid, err := s.sensorStore.Add(ctx, *sensorMongo)
	if err != nil {
		return "", err
	}
	return oid.Hex(), nil
}

func (s sensorMetadataService) AddWithLocationName(ctx context.Context, sensor SensorMetadataWithLocationName) (id string, err error) {
	loc, err := s.mapBox.FindLatLon(sensor.Location)
	if err != nil {
		return "", err
	}
	return s.Add(ctx, SensorMetadata{
		ID:       sensor.ID,
		Name:     sensor.Name,
		Location: loc,
		Tags:     sensor.Tags,
	})
}

func (s sensorMetadataService) FindNearestByLocatioName(ctx context.Context, location string) (sensor *SensorMetadata, err error) {
	loc, err := s.mapBox.FindLatLon(location)
	if err != nil {
		return nil, err
	}
	return s.FindNearest(ctx, loc.Lat, loc.Lon)

}

func (s sensorMetadataService) Update(ctx context.Context, sensor SensorMetadata) error {
	sensorMongo, err := sensor.ToDatabase()
	if err != nil {
		return err
	}
	err = s.sensorStore.Update(ctx, *sensorMongo)
	return err
}

func (s sensorMetadataService) FindNearest(ctx context.Context, lat, lon string) (sensor *SensorMetadata, err error) {
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return nil, err
	}
	lonF, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return nil, err
	}
	loc := db.Location{
		Lat: latF,
		Lon: lonF,
	}
	sensorMongo, err := s.sensorStore.FindNearest(ctx, loc)
	if err != nil {
		return nil, err
	}
	return FromDatabaseToSensorMetadata(*sensorMongo), nil
}

func (s sensorMetadataService) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.sensorStore.Delete(ctx, oid)
}
