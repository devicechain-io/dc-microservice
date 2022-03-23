/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package rdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/devicechain-io/dc-microservice/core"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	pgx "github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Manages lifecycle of relational database interactions.
type RdbManager struct {
	Microservice *core.Microservice
	Database     *gorm.DB
	Migrations   []*gormigrate.Migration

	lifecycle core.LifecycleManager
}

// Create a new rdb manager.
func NewRdbManager(ms *core.Microservice, callbacks core.LifecycleCallbacks,
	migrations []*gormigrate.Migration) *RdbManager {
	rdb := &RdbManager{
		Microservice: ms,
		Migrations:   migrations,
	}
	// Create lifecycle manager and channels for tracking shutdown.
	rdbname := fmt.Sprintf("%s-%s", ms.FunctionalArea, "rdb")
	rdb.lifecycle = core.NewLifecycleManager(rdbname, rdb, callbacks)
	return rdb
}

// Initialize component.
func (rdb *RdbManager) Initialize(ctx context.Context) error {
	return rdb.lifecycle.Initialize(ctx)
}

// Compute non-database connection URL for querying/creating database.
func (rdb *RdbManager) computeRootUrl() string {
	config := rdb.Microservice.InstanceConfiguration.Persistence.Rdb
	hostname := fmt.Sprintf("%v", config.Configuration["hostname"])
	port := fmt.Sprintf("%v", config.Configuration["port"])
	username := fmt.Sprintf("%v", config.Configuration["username"])
	password := fmt.Sprintf("%v", config.Configuration["password"])
	return fmt.Sprintf("postgres://%s:%s@%s:%s/postgres", username, password, hostname, port)
}

// Assure that database is created before connecting to it.
func (rdb *RdbManager) assurePostgresDatabase() error {
	log.Info().Str("database", rdb.Microservice.TenantId).Msg("Verifying that tenant database exists.")
	url := rdb.computeRootUrl()
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	// List all databases
	found := false
	result := conn.PgConn().ExecParams(context.Background(), "SELECT datname FROM pg_database WHERE datistemplate = false", [][]byte{}, nil, nil, nil)
	for result.NextRow() {
		currdb := string(result.Values()[0])
		if rdb.Microservice.TenantId == currdb {
			log.Info().Msg("Found existing tenant database.")
			found = true
		}
	}
	_, err = result.Close()
	if err != nil {
		return err
	}

	if !found {
		// Create tenant database.
		log.Info().Msg("Database was not found. Creating...")
		result := conn.PgConn().ExecParams(context.Background(), fmt.Sprintf("CREATE DATABASE %s", rdb.Microservice.TenantId),
			[][]byte{}, nil, nil, nil)
		_, err := result.Close()
		if err != nil {
			return err
		}
		log.Info().Str("database", rdb.Microservice.TenantId).Msg("Successfully created tenant database.")
	}

	return nil
}

// Compute non-database connection URL for querying/creating database.
func (rdb *RdbManager) computeDatabaseUrl() string {
	config := rdb.Microservice.InstanceConfiguration.Persistence.Rdb
	hostname := fmt.Sprintf("%v", config.Configuration["hostname"])
	port := fmt.Sprintf("%v", config.Configuration["port"])
	username := fmt.Sprintf("%v", config.Configuration["username"])
	password := fmt.Sprintf("%v", config.Configuration["password"])
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, hostname, port, rdb.Microservice.TenantId)
}

// Assure that functional area schema is created before connecting to it.
func (rdb *RdbManager) assurePostgresSchema() error {
	log.Info().Str("schema", rdb.Microservice.FunctionalArea).Msg("Verifying that schema exists.")
	url := rdb.computeDatabaseUrl()
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	// List all databases
	found := false
	result := conn.PgConn().ExecParams(context.Background(), "SELECT schema_name FROM information_schema.schemata", [][]byte{}, nil, nil, nil)
	for result.NextRow() {
		currsch := string(result.Values()[0])
		if rdb.Microservice.FunctionalArea == currsch {
			log.Info().Msg("Found existing schema for functional area.")
			found = true
		}
	}
	_, err = result.Close()
	if err != nil {
		return err
	}

	if !found {
		// Create functional area schema.
		log.Info().Msg("Schema was not found. Creating...")
		result := conn.PgConn().ExecParams(context.Background(), fmt.Sprintf("CREATE SCHEMA \"%s\"", rdb.Microservice.FunctionalArea),
			[][]byte{}, nil, nil, nil)
		_, err := result.Close()
		if err != nil {
			return err
		}
		log.Info().Str("database", rdb.Microservice.FunctionalArea).Msg("Successfully created schema.")
	}

	return nil
}

// Compute DSN for connecting to database.
func (rdb *RdbManager) computeDsn() string {
	config := rdb.Microservice.InstanceConfiguration.Persistence.Rdb
	hostname := fmt.Sprintf("%v", config.Configuration["hostname"])
	port := fmt.Sprintf("%v", config.Configuration["port"])
	username := fmt.Sprintf("%v", config.Configuration["username"])
	password := fmt.Sprintf("%v", config.Configuration["password"])
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
		username, password, hostname, rdb.Microservice.TenantId, port)
	log.Info().Str("username", username).Str("password", password).Str("hostname", hostname).
		Str("port", port).Msg("Initializing database connectivity")
	return dsn
}

// Lifecycle callback that runs initialization logic.
func (rdb *RdbManager) ExecuteInitialize(context.Context) error {
	// Make sure database exists before interacting with it.
	dbtype := rdb.Microservice.InstanceConfiguration.Persistence.Rdb.Type
	if strings.HasPrefix(dbtype, "postgres") {
		// Verify/create tenant database.
		err := rdb.assurePostgresDatabase()
		if err != nil {
			return err
		}

		// Verify/create functional area schema.
		err = rdb.assurePostgresSchema()
		if err != nil {
			return err
		}

		// Connect to database using params from instance configuration.
		dsn := rdb.computeDsn()
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   fmt.Sprintf("%s.", rdb.Microservice.FunctionalArea),
				SingularTable: false,
			}})
		if err != nil {
			return err
		}

		rdb.Database = db
	} else {
		return fmt.Errorf("relational database %s not currently supported", dbtype)
	}

	// Run migrations in a database independent manner.
	m := gormigrate.New(rdb.Database, gormigrate.DefaultOptions, rdb.Migrations)
	if err := m.Migrate(); err != nil {
		return err
	}

	return nil
}

// Start component.
func (rdb *RdbManager) Start(ctx context.Context) error {
	return rdb.lifecycle.Start(ctx)
}

// Lifecycle callback that runs startup logic.
func (rdb *RdbManager) ExecuteStart(context.Context) error {
	return nil
}

// Stop component.
func (rdb *RdbManager) Stop(ctx context.Context) error {
	return rdb.lifecycle.Stop(ctx)
}

// Lifecycle callback that runs shutdown logic.
func (rdb *RdbManager) ExecuteStop(context.Context) error {
	return nil
}

// Terminate component.
func (rdb *RdbManager) Terminate(ctx context.Context) error {
	return rdb.lifecycle.Terminate(ctx)
}

// Lifecycle callback that runs termination logic.
func (rdb *RdbManager) ExecuteTerminate(context.Context) error {
	sqldb, err := rdb.Database.DB()
	if err != nil {
		return err
	}
	sqldb.Close()
	return nil
}
