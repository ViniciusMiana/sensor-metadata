package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Points for improvement:
1 - Think further on the duplicate information Location and GeoJson
2 - Separate data objects and mongo store in different files
3 - Have a common mongo.Database object for all stores in the same microservice
4 - Structure errors
5 - Increase test coverage
6 - Add information on insert/update dates, any relevant change history
*/

const sensorCollectionName = "sensorMetadata"

// Sensor represents a sensor with meta-data
type Sensor struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Tags     []string           `bson:"tags"`
	Location *Location          `bson:"location"`
	GeoJson  *GeoJson           `bson:"geoJson"`
}

// Sensor represents a location with lat and lon
type Location struct {
	Lat float64 `bson:"lat"`
	Lon float64 `bson:"lon"`
}

// GeoJson is the mongo gis data format
type GeoJson struct {
	Point       string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

func (l Location) toDatabase() *GeoJson {
	// Note that mongo wants lon first then lat
	return &GeoJson{
		Point:       "Point",
		Coordinates: []float64{l.Lon, l.Lat},
	}
}

func (s *Sensor) prepareForDatabase() {
	if s.Location != nil {
		s.GeoJson = s.Location.toDatabase()
	}
}

// SensorStore represents the public interface of the sensorStore
type SensorStore interface {
	Add(ctx context.Context, sensor Sensor) (primitive.ObjectID, error)
	Update(ctx context.Context, sensor Sensor) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*Sensor, error)
	FindByName(ctx context.Context, name string) (*Sensor, error)
	FindNearest(ctx context.Context, location Location) (*Sensor, error)
}

type sensorStore struct {
	client   *mongo.Client
	database *mongo.Database
	sensors  *mongo.Collection
}

// NewSensorStore creates a new sensor store
func NewSensorStore(uri, databaseName string) (*sensorStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	database := client.Database(databaseName)
	sensors := database.Collection(sensorCollectionName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: nil,
		},
		{
			Keys:    bson.M{"geoJson": "2dsphere"},
			Options: options.Index().SetSphereVersion(2),
		},
	}
	_, err = sensors.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}
	return &sensorStore{client: client, database: database, sensors: sensors}, nil
}

// Add adds a new sensor to the store
func (store *sensorStore) Add(ctx context.Context, sensor Sensor) (primitive.ObjectID, error) {
	sensor.prepareForDatabase()
	res, err := store.sensors.InsertOne(ctx, sensor)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

// Update updates an existing sensor in the store
func (store *sensorStore) Update(ctx context.Context, sensor Sensor) error {
	sensor.prepareForDatabase()
	if sensor.ID == primitive.NilObjectID {
		return errors.New("Sensor ID can't be nil")
	}
	filter := bson.M{"_id": sensor.ID}
	update := bson.M{"$set": sensor}
	_, err := store.sensors.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a sensor from the store
func (store *sensorStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := store.sensors.DeleteOne(ctx, filter)
	return err
}

// FindByID finds a sensor by its ID
func (store *sensorStore) FindByID(ctx context.Context, id primitive.ObjectID) (*Sensor, error) {
	filter := bson.M{"_id": id}
	var result Sensor
	err := store.sensors.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindByName finds a sensor by its name
func (store *sensorStore) FindByName(ctx context.Context, name string) (*Sensor, error) {
	filter := bson.M{"name": name}
	var result Sensor
	err := store.sensors.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindNearest finds the sensor nearest to a location
func (store *sensorStore) FindNearest(ctx context.Context, location Location) (*Sensor, error) {
	loc := location.toDatabase()
	filter := bson.M{"geoJson": bson.M{"$near": bson.M{"$geometry": loc}}}
	var result Sensor
	err := store.sensors.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
