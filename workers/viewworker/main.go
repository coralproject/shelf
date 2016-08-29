package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/coralproject/shelf/internal/wire"
	validator "gopkg.in/bluesuncorp/validator.v8"
)

// Payload includes indications of views that need to
// be executed, along with information about how to
// connect to the graph and mongodb.
type Payload struct {
	Context        string            `json:"context"`
	CayleyHost     string            `json:"cayley_host" validate:"required,min=2"`
	CayleyDB       string            `json:"cayley_db" validate:"required,min=2"`
	CayleyUser     string            `json:"cayley_user"`
	CayleyPassword string            `json:"cayley_password"`
	MongoHost      string            `json:"mongo_host" validate:"required,min=2"`
	MongoDB        string            `json:"mongo_db" validate:"required,min=2"`
	MongoUser      string            `json:"mongo_user"`
	MongoPassword  string            `json:"mongo_password"`
	LoggingLevel   int               `json:"logging_level"`
	Views          []wire.ViewParams `json:"views" validate:"required,min=1"`
}

var (
	// validate is used to perform model field validation.
	validate *validator.Validate

	// payload is the current payload passed into the worker.
	payload Payload
)

// Validate checks the Payload value for consistency.
func (p *Payload) Validate() error {

	// Validate the database fields.
	if err := validate.Struct(p); err != nil {
		return err
	}

	// Validate the views.
	for _, view := range p.Views {
		if view.ResultsCollection == "" {
			return fmt.Errorf("View parameters should include a results collection")
		}
	}
	return nil
}

func init() {

	// Import the payload to the view worker.
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
	validate = validator.New(&validator.Config{TagName: "validate"})
	if err := payload.Validate(); err != nil {
		log.Error(payload.Context, "init", err, "Validating the payload")
		os.Exit(1)
	}

	// Initialize MongoDB.
	cfg := mongo.Config{
		Host:     payload.MongoHost,
		AuthDB:   payload.MongoDB,
		DB:       payload.MongoDB,
		User:     payload.MongoUser,
		Password: payload.MongoPassword,
		Timeout:  25 * time.Second,
	}
	if err := db.RegMasterSession(payload.Context, payload.MongoDB, cfg); err != nil {
		log.Error(payload.Context, "init", err, "Initializing MongoDB")
		os.Exit(1)
	}
}

func main() {
	log.Dev(payload.Context, "main", "Started")

	// Connect to MongoDB.
	mgoDB, err := db.NewMGO(payload.Context, payload.MongoDB)
	if err != nil {
		log.Error(payload.Context, "Mongo", err, "Completed")
		os.Exit(1)
	}
	defer mgoDB.CloseMGO(payload.Context)

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

	// Execute the Views.
	for _, view := range payload.Views {
		if _, err := wire.Execute(payload.Context, mgoDB, store, &view); err != nil {
			log.Error(payload.Context, "main", err, "Completed")
			os.Exit(1)
		}
	}

	log.Dev(payload.Context, "main", "Completed")
}
