package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/coralproject/shelf/internal/wire"
	validator "gopkg.in/bluesuncorp/validator.v8"
)

// Payload includes indications of relationships that need to
// be imported into a graph, along with information about how to
// connect to the graph.
type Payload struct {
	Context        string            `json:"context"`
	CayleyHost     string            `json:"cayley_host" validate:"required,min=2"`
	CayleyDB       string            `json:"cayley_db" validate:"required,min=2"`
	CayleyUser     string            `json:"cayley_user"`
	CayleyPassword string            `json:"cayley_password"`
	LoggingLevel   int               `json:"logging_level"`
	AddQuadData    []wire.QuadParams `json:"add_quad_data"`
	RemoveQuadData []wire.QuadParams `json:"remove_quad_data"`
}

var (
	// validate is used to perform model field validation.
	validate *validator.Validate

	// payload is the current payload passed into the worker.
	payload Payload
)

// Validate checks the Payload value for consistency.
func (p *Payload) Validate() error {

	// Validate the database parameters.
	if err := validate.Struct(p); err != nil {
		return err
	}

	// Validate the quad data.
	if p.AddQuadData == nil && p.RemoveQuadData == nil {
		return fmt.Errorf("No quad data provided in the payload.")
	}

	return nil
}

func init() {

	// Import the payload to the graph worker.
	payloadString, present := os.LookupEnv("PAYLOAD")
	if !present {
		fmt.Println("The PAYLOAD environmental var. is not defined, exiting.")
		os.Exit(1)
	}

	// Unmarshal the payload.
	var payload Payload
	if err := json.Unmarshal([]byte(payloadString), &payload); err != nil {
		fmt.Println("Unable to unmarshal the payload, exiting.")
		os.Exit(1)
	}

	// Initialize the log system.
	logLevel := func() int {
		return payload.LoggingLevel
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	// Validate the payload.
	validate = validator.New(&validator.Config{TagName: "Validating the payload"})
	if err := payload.Validate(); err != nil {
		log.Error(payload.Context, "init", err, "Completed")
		os.Exit(1)
	}
}

func main() {
	log.Dev(payload.Context, "main", "Started")

	// Connect to Cayley.
	opts := make(map[string]interface{})
	opts["database_name"] = payload.CayleyDB
	opts["username"] = payload.CayleyUser
	opts["password"] = payload.CayleyPassword
	store, err := cayley.NewGraph("mongo", payload.CayleyHost, opts)
	if err != nil {
		log.Error(payload.Context, "main", err, "Completed")
		os.Exit(1)
	}

	// Add relationships.
	if payload.AddQuadData != nil {
		if err := wire.AddToGraph(payload.Context, store, payload.AddQuadData); err != nil {
			log.Error(payload.Context, "main", err, "Completed")
			os.Exit(1)
		}
	}

	// Remove relationships.
	if payload.RemoveQuadData != nil {
		if err := wire.RemoveFromGraph(payload.Context, store, payload.RemoveQuadData); err != nil {
			log.Error(payload.Context, "main", err, "Completed")
			os.Exit(1)
		}
	}

	log.Dev(payload.Context, "main", "Completed")
}
