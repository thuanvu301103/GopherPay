package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/thuanvu301103/auth-service/internal/auth"
)

func main() {
	// Initialize the loader for PostgreSQL
	// Add all models that need to be tracked as tables here
	stmts, err := gormschema.New("postgres").Load(&auth.User{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load GORM schema: %v\n", err)
		os.Exit(1)
	}
	// Output the schema to stdout for Atlas to read
	io.WriteString(os.Stdout, stmts)
}
