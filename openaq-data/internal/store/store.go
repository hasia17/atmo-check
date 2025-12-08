package store

import (
	"context"
	"fmt"
	"openaq-data/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store struct {
	mongoClient  *mongo.Client
	stationsColl *mongo.Collection
	measuresColl *mongo.Collection
}

func New(mongoURI string) (*Store, error) {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := mongoClient.Database("openaq")
	stationsColl := db.Collection("stations")
	measuresColl := db.Collection("measurements")

	return &Store{
		mongoClient:  mongoClient,
		stationsColl: stationsColl,
		measuresColl: measuresColl,
	}, nil
}

func (s *Store) StoreStations(ctx context.Context, stations []types.Station) error {
	for _, station := range stations {
		if err := s.storeStation(ctx, station); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) storeStation(ctx context.Context, station types.Station) error {
	_, err := s.stationsColl.UpdateOne(ctx,
		bson.M{"_id": station.Id},
		bson.M{"$set": station},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to store station: %w", err)
	}
	return nil
}

func (s *Store) GetStations(ctx context.Context) ([]types.Station, error) {
	cursor, err := s.stationsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stations: %w", err)
	}
	defer cursor.Close(ctx)

	var stations []types.Station
	if err := cursor.All(ctx, &stations); err != nil {
		return nil, fmt.Errorf("failed to decode stations: %w", err)
	}
	return stations, nil
}

func (s *Store) GetStationByID(ctx context.Context, id int32) (*types.Station, error) {
	var station types.Station
	err := s.stationsColl.FindOne(ctx, bson.M{"_id": id}).Decode(&station)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch station by id: %w", err)
	}
	return &station, nil
}

func (s *Store) GetParametersByStationID(ctx context.Context, id int32) ([]types.Parameter, error) {
	var station types.Station
	err := s.stationsColl.FindOne(ctx, bson.M{"_id": id}).Decode(&station)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch station for parameters: %w", err)
	}
	return *station.Parameters, nil
}

func (s *Store) StoreMeasurements(ctx context.Context, m []types.Measurement) error {
	for _, measure := range m {
		_, err := s.measuresColl.InsertOne(ctx, measure)
		if err != nil {
			return fmt.Errorf("failed to store measurement: %w", err)
		}
	}
	return nil
}

func (s *Store) GetMeasurementsByStation(ctx context.Context, stationID int32, limit int64) ([]types.Measurement, error) {
	filter := bson.M{"stationId": stationID}
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

	var measurements []types.Measurement
	if err := cursor.All(ctx, &measurements); err != nil {
		return nil, err
	}
	return measurements, nil
}

func (s *Store) GetLatestMeasurementsByStation(ctx context.Context, stationID int32) ([]types.Measurement, error) {
	filter := bson.M{"stationid": stationID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := s.measuresColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	latest := make(map[int32]types.Measurement) // parameterID -> Measurement
	for cursor.Next(ctx) {
		var m types.Measurement
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		if _, exists := latest[*m.SensorId]; !exists {
			latest[*m.SensorId] = m
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	result := make([]types.Measurement, 0, len(latest))
	for _, m := range latest {
		result = append(result, m)
	}
	return result, nil
}

func (s *Store) DeleteMeasurementsForStation(ctx context.Context, stationID int32) error {
	_, err := s.measuresColl.DeleteMany(ctx, bson.M{"stationId": stationID})
	if err != nil {
		return fmt.Errorf("failed to delete measurements for station: %w", err)
	}
	return nil
}

func (s *Store) HasData(ctx context.Context) (bool, error) {
	stationCount, err := s.stationsColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, fmt.Errorf("failed to count stations: %w", err)
	}
	if stationCount == 0 {
		return false, nil
	}

	measureCount, err := s.measuresColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, fmt.Errorf("failed to count measurements: %w", err)
	}
	return measureCount > 0, nil
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
