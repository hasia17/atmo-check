//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
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

	// Load OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("./api/openapi.yaml")
	if err != nil {
		return fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	// Validate the spec
	if err := doc.Validate(context.Background()); err != nil {
		return fmt.Errorf("invalid OpenAPI spec: %w", err)
	}

	// Configure code generation for types only
	config := codegen.Configuration{
		PackageName: "api",
		Generate: codegen.GenerateOptions{
			Models: true,
			// TODO: check fiber servers
		},
	}

	// Generate the code
	code, err := codegen.Generate(doc, config)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Write to file
	outFile, err := os.Create("api/types.go")
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.WriteString(outFile, code)
	return err
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
