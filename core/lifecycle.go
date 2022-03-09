/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package core

import (
	"context"
	"errors"
)

type LifecycleState int64

// Enumeration of lifecycle states
//go:generate stringer -type=LifecycleState
const (
	Uninitialized LifecycleState = iota
	Initializing
	Initialized
	Starting
	Started
	Stopping
	Stopped
	Terminating
	Terminated
)

// Common lifecycle concept for components
type LifecycleComponent interface {
	// Initialize component. Happens once on startup.
	lifecycleInitialize(context.Context) error

	// Start component. May happen on startup or after stop.
	lifecycleStart(context.Context) error

	// Stop a started component.
	lifecycleStop(context.Context) error

	// Terminate component.
	lifecycleTerminate(context.Context) error
}

type LifecycleManager struct {
	Component LifecycleComponent
	State     LifecycleState
}

func NewLifecycleManager(component LifecycleComponent) *LifecycleManager {
	mgr := &LifecycleManager{Component: component, State: Uninitialized}
	return mgr
}

// Handle component initialization
func (mgr *LifecycleManager) initialize(ctx context.Context) error {
	if mgr.State != Uninitialized {
		return errors.New("attempting to initialize component that is already initialized")
	}
	prev := mgr.State
	mgr.State = Initializing
	err := mgr.Component.lifecycleInitialize(ctx)
	if err != nil {
		mgr.State = prev
		return err
	}
	mgr.State = Initialized
	return nil
}

// Handle component startup
func (mgr *LifecycleManager) start(ctx context.Context) error {
	if mgr.State == Uninitialized {
		return errors.New("attempting to start an uninitialized component")
	}
	if mgr.State == Starting {
		return errors.New("attempting to start a component that is already starting")
	}
	if mgr.State == Started {
		return errors.New("attempting to start a component that is already started")
	}
	if mgr.State == Stopping {
		return errors.New("attempting to start a component that is stopping")
	}
	if mgr.State == Terminating {
		return errors.New("attempting to start a component that is terminating")
	}
	if mgr.State == Terminated {
		return errors.New("attempting to start a component that is terminated")
	}
	prev := mgr.State
	mgr.State = Starting
	err := mgr.Component.lifecycleStart(ctx)
	if err != nil {
		mgr.State = prev
		return err
	}
	mgr.State = Started
	return nil
}

// Handle component shutdown
func (mgr *LifecycleManager) stop(ctx context.Context) error {
	if mgr.State == Uninitialized {
		return errors.New("attempting to stop an uninitialized component")
	}
	if mgr.State == Starting {
		return errors.New("attempting to stop a component that is partially started")
	}
	if mgr.State == Stopping {
		return errors.New("attempting to stop a component that is already stopping")
	}
	if mgr.State == Stopped {
		return errors.New("attempting to stop a component that is already stopped")
	}
	if mgr.State == Terminating {
		return errors.New("attempting to stop a component that is terminating")
	}
	if mgr.State == Terminated {
		return errors.New("attempting to stop a component that is terminated")
	}
	prev := mgr.State
	mgr.State = Stopping
	err := mgr.Component.lifecycleStop(ctx)
	if err != nil {
		mgr.State = prev
		return err
	}
	mgr.State = Stopped
	return nil
}

// Handle component termination
func (mgr *LifecycleManager) terminate(ctx context.Context) error {
	if mgr.State == Uninitialized {
		return errors.New("attempting to terminate component that is not initialized")
	}
	if mgr.State != Stopped {
		return errors.New("attempting to terminate component that is not stopped")
	}
	prev := mgr.State
	mgr.State = Terminating
	err := mgr.Component.lifecycleTerminate(ctx)
	if err != nil {
		mgr.State = prev
		return err
	}
	mgr.State = Terminated
	return nil
}
