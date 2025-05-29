package internal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store struct {
	mongoClient   *mongo.Client
	locationsColl *mongo.Collection
	measuresColl  *mongo.Collection
}

func NewStore(mongoURI string) (*Store, error) {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := mongoClient.Database("atmo-check")
	locationsColl := db.Collection("locations")
	measuresColl := db.Collection("measurements")

	return &Store{
		mongoClient:   mongoClient,
		locationsColl: locationsColl,
		measuresColl:  measuresColl,
	}, nil
}

func (s *Store) StoreLocation(ctx context.Context, location location) error {
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

func (s *Store) GetLocations(ctx context.Context) ([]location, error) {
	cursor, err := s.locationsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer cursor.Close(ctx)

	var locations []location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, fmt.Errorf("failed to decode locations: %w", err)
	}
	return locations, nil
}

func (s *Store) StoreMeasurement(ctx context.Context, m measurement) error {
	t, err := time.Parse(time.RFC3339, m.DateTime.Utc)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}
	m.Timestamp = t

	_, err = s.measuresColl.InsertOne(ctx, bson.M{
		"timestamp": m.Timestamp,
		"metadata": bson.M{
			"location_id": m.LocationsId,
			"parameter":   m.Parameter,
		},
		"value":       m.Value,
		"coordinates": m.Coordinates,
		"sensorsId":   m.SensorsId,
	})
	return err
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
