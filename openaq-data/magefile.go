//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Build() error {
	fmt.Println("Building openaq-data...")
	cmd := exec.Command("go", "build", "-o", "openaq-data", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test() error {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Clean() error {
	fmt.Println("Cleaning up...")
	return os.Remove("openaq-data")
}

func OapiGenerate() error {
	fmt.Println("Generating Go code from OpenAPI spec...")

	spec := "api.yaml"
	output := filepath.Join("gen", "api.gen.go")
	pkg := "gen"
	generate := "types,client,server"
	module := "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest"

	cmd := exec.Command("go", "run", module,
		"-generate", generate,
		"-o", output,
		"-package", pkg,
		spec,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DropDB() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27018"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Dropping database 'atmo-check'...")
	if err := client.Database("atmo-check").Drop(ctx); err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	fmt.Println("Database dropped.")
	return nil
}
