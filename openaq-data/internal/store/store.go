package store

import (
	"context"
	"fmt"
	"openaq-data/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TODO: maybe split this into multiple files
type Storer interface {
	StoreLocations(ctx context.Context, locations []models.Location) error
	GetLocations(ctx context.Context) ([]models.Location, error)
	GetLocationByID(ctx context.Context, id int32) (*models.Location, error)

	StoreMeasurements(ctx context.Context, m []models.Measurement) error
	GetMeasurementsByLocation(ctx context.Context, locationId int32) ([]models.Measurement, error)
	DeleteMeasurementsForLocation(ctx context.Context, locationID int32) error

	StoreParameters(ctx context.Context, parameters []models.Parameter) error
	GetParameters(ctx context.Context) ([]models.Parameter, error)

	Close() error
}
type Store struct {
	mongoClient    *mongo.Client
	locationsColl  *mongo.Collection
	measuresColl   *mongo.Collection
	parametersColl *mongo.Collection
}

func New(mongoURI string) (*Store, error) {
	mongoClient, err := configureClient(context.Background(), mongoURI)
	if err != nil {
		return nil, err
	}

	db := mongoClient.Database("openaq")
	locationsColl := db.Collection("locations")
	measuresColl := db.Collection("measurements")
	parametersColl := db.Collection("parameters")

	return &Store{
		mongoClient:    mongoClient,
		locationsColl:  locationsColl,
		measuresColl:   measuresColl,
		parametersColl: parametersColl,
	}, nil
}

func configureClient(ctx context.Context, mongoURI string) (*mongo.Client, error) {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return mongoClient, err
}

func (s *Store) StoreLocations(ctx context.Context, locations []models.Location) error {
	for _, l := range locations {
		if err := s.storeLocation(ctx, l); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) storeLocation(ctx context.Context, location models.Location) error {
	_, err := s.locationsColl.UpdateOne(ctx,
		bson.M{"_id": location.Id},
		bson.M{"$set": location},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to store location: %w", err)
	}
	return nil
}

func (s *Store) GetLocations(ctx context.Context) ([]models.Location, error) {
	cursor, err := s.locationsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer cursor.Close(ctx)

	var locations []models.Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, fmt.Errorf("failed to decode locations: %w", err)
	}
	return locations, nil
}

func (s *Store) GetLocationByID(ctx context.Context, id int32) (*models.Location, error) {
	var loc models.Location
	err := s.locationsColl.FindOne(ctx, bson.M{"_id": id}).Decode(&loc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location by id: %w", err)
	}
	return &loc, nil
}

func (s *Store) StoreMeasurements(ctx context.Context, m []models.Measurement) error {
	for _, measure := range m {
		_, err := s.measuresColl.InsertOne(ctx, measure)
		if err != nil {
			return fmt.Errorf("failed to store measurement: %w", err)
		}
	}
	return nil
}

func (s *Store) GetMeasurementsByLocation(ctx context.Context, locationId int32) ([]models.Measurement, error) {
	filter := bson.M{"locationsId": locationId}
	opts := options.Find().SetSort(
		bson.D{{
			Key:   "timestamp",
			Value: -1,
		}},
	)
	cursor, err := s.measuresColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var measurements []models.Measurement
	if err := cursor.All(ctx, &measurements); err != nil {
		return nil, err
	}
	return measurements, nil
}

func (s *Store) DeleteMeasurementsForLocation(ctx context.Context, locationID int32) error {
	_, err := s.measuresColl.DeleteMany(ctx, bson.M{"locationsId": locationID})
	if err != nil {
		return fmt.Errorf("failed to delete measurements for location: %w", err)
	}
	return nil
}

func (s *Store) GetParameters(ctx context.Context) ([]models.Parameter, error) {
	cursor, err := s.parametersColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parameters: %w", err)
	}
	defer cursor.Close(ctx)

	var parameters []models.Parameter
	if err := cursor.All(ctx, &parameters); err != nil {
		return nil, fmt.Errorf("failed to decode parameters: %w", err)
	}
	return parameters, nil
}

func (s *Store) StoreParameters(ctx context.Context, parameters []models.Parameter) error {
	for _, parameter := range parameters {
		if err := s.storeParameter(ctx, parameter); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) storeParameter(ctx context.Context, parameter models.Parameter) error {
	_, err := s.parametersColl.UpdateOne(ctx,
		bson.M{"_id": parameter.Id},
		bson.M{"$set": parameter},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to store parameter: %w", err)
	}
	return nil
}

func (s *Store) Close() error {
	if s.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.mongoClient.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
		}
	}
	return nil
}
