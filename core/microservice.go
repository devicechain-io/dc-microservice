/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devicechain-io/dc-microservice/config"
	"github.com/olekukonko/tablewriter"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Primary microservice implementation
type Microservice struct {
	StartTime time.Time

	// Passed from environment
	TenantId         string
	TenantName       string
	MicroserviceId   string
	MicroserviceName string
	FunctionalArea   string

	// Configuration content
	InstanceConfiguration        config.InstanceConfiguration
	MicroserviceConfigurationRaw []byte

	// Internal lifeycle processing
	lifecycle LifecycleManager
	shutdown  chan os.Signal
	done      chan bool
}

// Create a new microservice instance
func NewMicroservice(callbacks LifecycleCallbacks) *Microservice {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ms := &Microservice{}
	ms.StartTime = time.Now()
	ms.TenantId = os.Getenv(ENV_TENANT_ID)
	ms.TenantName = os.Getenv(ENV_TENANT_NAME)
	ms.MicroserviceId = os.Getenv(ENV_MICROSERVICE_ID)
	ms.MicroserviceName = os.Getenv(ENV_MICROSERVICE_NAME)
	ms.FunctionalArea = os.Getenv(ENV_MS_FUNCTIONAL_AREA)

	// Create lifecycle manager and channels for tracking shutdown.
	ms.lifecycle = NewLifecycleManager(ms.FunctionalArea, ms, callbacks)
	ms.done = make(chan bool, 1)
	ms.shutdown = make(chan os.Signal, 1)

	// Hook interrupt and terminate signals for graceful shutdown
	signal.Notify(ms.shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Async handle shutdown on signals
	go func() {
		sig := <-ms.shutdown
		fmt.Println()
		log.Warn().Msgf("Received signal '%v'. Shutting down gracefully...", sig)
		ms.ShutDownNow()
	}()

	return ms
}

// Prints a banner to the console
func (ms *Microservice) Banner() {
	fmt.Println(color.HiGreenString(`
    ____            _           ________          _     
   / __ \___ _   __(_)_______  / ____/ /_  ____ _(_)___ 
  / / / / _ \ | / / / ___/ _ \/ /   / __ \/ __  / / __ \
 / /_/ /  __/ |/ / / /__/  __/ /___/ / / / /_/ / / / / /
/_____/\___/|___/_/\___/\___/\____/_/ /_/\__,_/_/_/ /_/ 

`))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	data := [][]string{
		{"Tenant", fmt.Sprintf("%s (%s)", ms.TenantName, ms.TenantId)},
		{"Microservice", fmt.Sprintf("%s (%s)", ms.MicroserviceName, ms.MicroserviceId)},
	}
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	fmt.Println()
}

// Create microservice and initialize/start it.
func (ms *Microservice) Run() error {
	ms.Banner()
	log.Info().Msg("Creating new microservice and running intialization/startup...")

	go func() {
		err := ms.InitializeAndStart()
		if err != nil {
			ms.done <- true
		}
	}()

	ms.waitForShutdown()
	return nil
}

// Issue initialize and start commands to microservice
func (ms *Microservice) InitializeAndStart() error {
	err := ms.Initialize(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to initialize microservice")
		return err
	}
	err = ms.Start(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to start microservice")
		return err
	}
	return nil
}

// Issue stop and terminate commands to microservice
func (ms *Microservice) ShutDownNow() {
	err := ms.Stop(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to stop microservice")
		ms.done <- true
		return
	}
	err = ms.Terminate(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to terminate microservice")
		ms.done <- true
		return
	}

	ms.done <- true
}

// Wait for microservice to shut down
func (ms *Microservice) waitForShutdown() {
	<-ms.done
}

// Reloads instance configuration from configmap volume mapping
func (ms *Microservice) ReloadInstanceConfiguration() error {
	bytes, err := os.ReadFile("/etc/dci-config/instance")
	if err != nil {
		return err
	}
	config := &config.InstanceConfiguration{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	ms.InstanceConfiguration = *config
	return nil
}

// Reloads microservice configuration from configmap volume mapping
func (ms *Microservice) ReloadMicroserviceConfiguration() error {
	fa, found := os.LookupEnv(ENV_MS_FUNCTIONAL_AREA)
	if !found {
		return fmt.Errorf("environment variable for functional area (%s) not set", ENV_MS_FUNCTIONAL_AREA)
	}

	bytes, err := os.ReadFile(fmt.Sprintf("/etc/dct-config/%s", fa))
	if err != nil {
		return err
	}
	ms.MicroserviceConfigurationRaw = bytes
	return nil
}

// Initialize microservice
func (ms *Microservice) Initialize(ctx context.Context) error {
	return ms.lifecycle.Initialize(ctx)
}

// Initialize microservice (as called by lifecycle manager)
func (ms *Microservice) ExecuteInitialize(ctx context.Context) error {
	err := ms.ReloadInstanceConfiguration()
	if err != nil {
		return err
	}
	log.Info().Msg("Successfully loaded instance configuration.")

	err = ms.ReloadMicroserviceConfiguration()
	if err != nil {
		return err
	}
	log.Info().Msg("Successfully loaded microservice configuration.")

	return err
}

// Start microservice
func (ms *Microservice) Start(ctx context.Context) error {
	return ms.lifecycle.Start(ctx)
}

// Start microservice (as called by lifecycle manager)
func (ms *Microservice) ExecuteStart(ctx context.Context) error {
	log.Info().Msg("Microservice started.")
	return nil
}

// Stop microservice
func (ms *Microservice) Stop(ctx context.Context) error {
	return ms.lifecycle.Stop(ctx)
}

// Stop microservice (as called by lifecycle manager)
func (ms *Microservice) ExecuteStop(ctx context.Context) error {
	log.Info().Msg("Microservice stopped.")
	return nil
}

// Terminate microservice
func (ms *Microservice) Terminate(ctx context.Context) error {
	return ms.lifecycle.Terminate(ctx)
}

// Terminate microservice (as called by lifecycle manager)
func (ms *Microservice) ExecuteTerminate(ctx context.Context) error {
	log.Info().Msg("Microservice terminated.")
	return nil
}
