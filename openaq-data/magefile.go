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
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	openAPISpecPath = "openapi.yaml"
)

type Test mg.Namespace

func (t Test) All() error {
	return sh.RunV("go", "test", "./...")
}

func (t Test) Races() error {
	return sh.RunV("go", "test", "./...", "--race")
}

type genConfig struct {
	packageName    string
	outputFilePath string
	genOpts        codegen.GenerateOptions
}

type Gen mg.Namespace

func (g Gen) Types() error {
	return g.generate(genConfig{
		packageName:    "api",
		outputFilePath: "internal/api/types.go",
		genOpts: codegen.GenerateOptions{
			Models: true,
		},
	})
}

func (g Gen) Api() error {
	return g.generate(genConfig{
		packageName:    "api",
		outputFilePath: "internal/api/server.go",
		genOpts: codegen.GenerateOptions{
			StdHTTPServer: true,
		},
	})
}

func (g Gen) generate(opts genConfig) error {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(openAPISpecPath)
	if err != nil {
		return fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	if err := doc.Validate(context.Background()); err != nil {
		return fmt.Errorf("invalid OpenAPI spec: %w", err)
	}

	config := codegen.Configuration{
		PackageName: opts.packageName,
		Generate:    opts.genOpts,
	}

	code, err := codegen.Generate(doc, config)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	outFile, err := os.Create(opts.outputFilePath)
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

func Build() error {
	fmt.Println("Building openaq-data...")
	cmd := exec.Command("go", "build", "-o", "openaq-data", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
