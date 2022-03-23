/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package rdb

import (
	"context"
	"fmt"

	"github.com/devicechain-io/dc-microservice/core"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Manages lifecycle of relational database interactions.
type RdbManager struct {
	Microservice *core.Microservice
	Database     *gorm.DB

	lifecycle core.LifecycleManager
}

// Create a new rdb manager.
func NewRdbManager(ms *core.Microservice, callbacks core.LifecycleCallbacks) *RdbManager {
	rdb := &RdbManager{
		Microservice: ms,
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

// Lifecycle callback that runs initialization logic.
func (rdb *RdbManager) ExecuteInitialize(context.Context) error {
	config := rdb.Microservice.InstanceConfiguration.Persistence.Rdb
	hostname := fmt.Sprintf("%v", config.Configuration["hostname"])
	port := fmt.Sprintf("%v", config.Configuration["port"])
	username := fmt.Sprintf("%v", config.Configuration["username"])
	password := fmt.Sprintf("%v", config.Configuration["password"])
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
		username, password, hostname, rdb.Microservice.TenantId, port)
	log.Info().Str("username", username).Str("password", password).Str("hostname", hostname).
		Str("port", port).Msg("Initializing database connectivity")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.", rdb.Microservice.FunctionalArea),
			SingularTable: false,
		}})
	if err != nil {
		return err
	}

	rdb.Database = db
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
