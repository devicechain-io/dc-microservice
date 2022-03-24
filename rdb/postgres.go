/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package rdb

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Compute non-database connection URL for querying/creating database.
func (rdb *RdbManager) computePostgresRootUrl(pgconfig *PostgresConfig) string {
	config := rdb.Microservice.InstanceConfiguration.Persistence.Rdb
	hostname := fmt.Sprintf("%v", config.Configuration["hostname"])
	port := fmt.Sprintf("%v", config.Configuration["port"])
	username := fmt.Sprintf("%v", config.Configuration["username"])
	password := fmt.Sprintf("%v", config.Configuration["password"])
	return fmt.Sprintf("postgres://%s:%s@%s:%s/postgres", username, password, hostname, port)
}

// Assure that database is created before connecting to it.
func (rdb *RdbManager) assurePostgresDatabase(pgconfig *PostgresConfig) error {
	log.Info().Str("database", rdb.Microservice.TenantId).Msg("Verifying that tenant database exists.")
	url := rdb.computePostgresRootUrl(pgconfig)
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
func (rdb *RdbManager) computePostgresTenantDatabaseUrl(pg *PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", pg.Username, pg.Password, pg.Hostname, pg.Port, rdb.Microservice.TenantId)
}

// Assure that functional area schema is created before connecting to it.
func (rdb *RdbManager) assurePostgresSchema(pgconfig *PostgresConfig) error {
	log.Info().Str("schema", rdb.Microservice.FunctionalArea).Msg("Verifying that schema exists.")
	url := rdb.computePostgresTenantDatabaseUrl(pgconfig)
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
		log.Info().Str("schema", rdb.Microservice.FunctionalArea).Msg("Successfully created schema.")
	}

	return nil
}

// Compute DSN for connecting to database.
func (rdb *RdbManager) computePostgresDsn(pg *PostgresConfig) string {
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable",
		pg.Username, pg.Password, pg.Hostname, rdb.Microservice.TenantId, pg.Port)
	log.Info().Str("username", pg.Username).Str("password", pg.Password).Str("hostname", pg.Hostname).
		Int32("port", pg.Port).Msg("Initializing database connectivity")
	return dsn
}

// Boostrap a postgres database/schema.
func (rdb *RdbManager) bootstrapPostgres(pgconf *PostgresConfig) error {
	// Verify/create tenant database.
	err := rdb.assurePostgresDatabase(pgconf)
	if err != nil {
		return err
	}

	// Verify/create functional area schema.
	err = rdb.assurePostgresSchema(pgconf)
	if err != nil {
		return err
	}

	return nil
}

// Initialize a postgres database.
func (rdb *RdbManager) initializePostgres() error {
	pgconf, err := convertToPostgresConfig(rdb.Microservice.InstanceConfiguration.Persistence.Rdb)
	if err != nil {
		return err
	}

	// Bootstrap the postgres database/schema if needed.
	err = rdb.bootstrapPostgres(pgconf)
	if err != nil {
		return err
	}

	// Connect to database using params from instance configuration.
	dsn := rdb.computePostgresDsn(pgconf)
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

	sqldb, err := rdb.Database.DB()
	if err != nil {
		return err
	}
	sqldb.SetMaxIdleConns(1)
	sqldb.SetMaxOpenConns(int(pgconf.MaxConnections))
	sqldb.SetConnMaxLifetime(time.Hour)
	log.Info().
		Int32("max_connections", pgconf.MaxConnections).Msg("Created connection pool.")

	return nil
}
