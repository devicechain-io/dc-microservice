/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// Primary microservice implementation
type Microservice struct {
	StartTime time.Time

	lifecycle *LifecycleManager
	shutdown  chan os.Signal
	done      chan bool
}

// Create a new microservice instance
func NewMicroservice() *Microservice {
	ms := &Microservice{}
	ms.StartTime = time.Now()
	ms.lifecycle = NewLifecycleManager(ms)
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

// Issue stop and terminate commands to microservice
func (ms *Microservice) ShutDownNow() {
	err := ms.Stop(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("unable to stop microservice")
		ms.done <- true
		return
	}
	err = ms.Terminate(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("unable to terminate microservice")
		ms.done <- true
		return
	}

	ms.done <- true
}

// Wait for microservice to shut down
func (ms *Microservice) WaitForShutdown() {
	<-ms.done
}

// Initialize microservice
func (ms *Microservice) Initialize(ctx context.Context) error {
	return ms.lifecycle.initialize(ctx)
}

// Initialize microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleInitialize(ctx context.Context) error {
	log.Info().Msg("Microservice initialized.")
	return nil
}

// Start microservice
func (ms *Microservice) Start(ctx context.Context) error {
	return ms.lifecycle.start(ctx)
}

// Start microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleStart(ctx context.Context) error {
	log.Info().Msg("Microservice started.")
	return nil
}

// Stop microservice
func (ms *Microservice) Stop(ctx context.Context) error {
	return ms.lifecycle.stop(ctx)
}

// Stop microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleStop(ctx context.Context) error {
	log.Info().Msg("Microservice stopped.")
	return nil
}

// Stop microservice
func (ms *Microservice) Terminate(ctx context.Context) error {
	return ms.lifecycle.terminate(ctx)
}

// Stop microservice (as called by lifecycle manager)
func (ms *Microservice) lifecycleTerminate(ctx context.Context) error {
	log.Info().Msg("Microservice terminated.")
	err := errors.New("this is a test error")
	log.Error().Err(err).Msg("outter error description!")
	return nil
}
