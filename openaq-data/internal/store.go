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

func (s *Store) StoreLocation(ctx context.Context, location Location) error {
	_, err := s.locationsColl.UpdateOne(ctx,
		bson.M{"_id": location.ID},
		bson.M{"$set": location},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to store location: %w", err)
	}
	return nil
}

func (s *Store) GetLocations(ctx context.Context) ([]Location, error) {
	cursor, err := s.locationsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer cursor.Close(ctx)

	var locations []Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, fmt.Errorf("failed to decode locations: %w", err)
	}
	return locations, nil
}

func (s *Store) StoreMeasurement(ctx context.Context, m Measurement) error {
	_, err := s.measuresColl.InsertOne(ctx, m)
	return err
}

func (s *Store) GetMeasurementsByLocation(ctx context.Context, locationID int32, limit int64) ([]Measurement, error) {
	filter := bson.M{"locationId": locationID}
	opts := options.Find().SetSort(
		bson.D{{
			Key:   "timestamp",
			Value: -1,
		}},
	)
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := s.measuresColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var measurements []Measurement
	if err := cursor.All(ctx, &measurements); err != nil {
		return nil, err
	}
	return measurements, nil
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
