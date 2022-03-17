/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devicechain-io/dc-k8s/api/v1beta1"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Primary microservice implementation
type Microservice struct {
	StartTime time.Time

	TenantMicroservice v1beta1.TenantMicroservice

	lifecycle LifecycleManager
	shutdown  chan os.Signal
	done      chan bool
}

// Create a new microservice instance
func NewMicroservice(name string, callbacks LifecycleCallbacks) *Microservice {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ms := &Microservice{}
	ms.StartTime = time.Now()
	ms.lifecycle = NewLifecycleManager(name, ms, callbacks)
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
func banner() {
	fmt.Println(color.HiGreenString(`
    ____            _           ________          _     
   / __ \___ _   __(_)_______  / ____/ /_  ____ _(_)___ 
  / / / / _ \ | / / / ___/ _ \/ /   / __ \/ __  / / __ \
 / /_/ /  __/ |/ / / /__/  __/ /___/ / / / /_/ / / / / /
/_____/\___/|___/_/\___/\___/\____/_/ /_/\__,_/_/_/ /_/ 

`))
}

// Create microservice and initialize/start it.
func (ms *Microservice) Run() error {
	banner()
	log.Info().Msg("Creating new microservice and running intialization/startup...")

	go func() {
		err := ms.initializeAndStart()
		if err != nil {
			ms.done <- true
		}
	}()

	ms.waitForShutdown()
	return nil
}

// Issue initialize and start commands to microservice
func (ms *Microservice) initializeAndStart() error {
	err := ms.initialize(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to initialize microservice")
		return err
	}
	err = ms.start(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to start microservice")
		return err
	}
	return nil
}

// Issue stop and terminate commands to microservice
func (ms *Microservice) ShutDownNow() {
	err := ms.stop(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to stop microservice")
		ms.done <- true
		return
	}
	err = ms.terminate(context.Background())
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

// Initialize microservice
func (ms *Microservice) initialize(ctx context.Context) error {
	return ms.lifecycle.initialize(ctx)
}

// Initialize tenantmicroservice resource from k8s
func (ms *Microservice) initTenantMicroservice() error {
	tm, err := ms.getTenantMicroservice()
	if err != nil {
		return err
	}
	log.Info().Str("tenant", tm.Spec.TenantId).Str("microservice", tm.Spec.MicroserviceId).Msg("Found tenant microservice")
	ms.TenantMicroservice = *tm
	return nil
}

// Initialize microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleInitialize(ctx context.Context) error {
	err := ms.initTenantMicroservice()
	return err
}

// Start microservice
func (ms *Microservice) start(ctx context.Context) error {
	return ms.lifecycle.start(ctx)
}

// Start microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleStart(ctx context.Context) error {
	log.Info().Msg("Microservice started.")
	return nil
}

// Stop microservice
func (ms *Microservice) stop(ctx context.Context) error {
	return ms.lifecycle.stop(ctx)
}

// Stop microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleStop(ctx context.Context) error {
	log.Info().Msg("Microservice stopped.")
	return nil
}

// Stop microservice
func (ms *Microservice) terminate(ctx context.Context) error {
	return ms.lifecycle.terminate(ctx)
}

// Stop microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleTerminate(ctx context.Context) error {
	log.Info().Msg("Microservice terminated.")
	return nil
}
