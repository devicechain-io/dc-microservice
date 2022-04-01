/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package kafka

import (
	"context"
	"fmt"

	"github.com/devicechain-io/dc-microservice/core"
	kafka "github.com/segmentio/kafka-go"
)

// Manages lifecycle of kafka interactions.
type KafkaManager struct {
	Microservice *core.Microservice

	lifecycle core.LifecycleManager
}

// Create a new kafka manager.
func NewKafkaManager(ms *core.Microservice, callbacks core.LifecycleCallbacks) *KafkaManager {
	kmgr := &KafkaManager{
		Microservice: ms,
	}
	// Create lifecycle manager.
	kfkaname := fmt.Sprintf("%s-%s", ms.FunctionalArea, "kafka")
	kmgr.lifecycle = core.NewLifecycleManager(kfkaname, kmgr, callbacks)
	return kmgr
}

// Initialize component.
func (kmgr *KafkaManager) Initialize(ctx context.Context) error {
	return kmgr.lifecycle.Initialize(ctx)
}

// Lifecycle callback that runs initialization logic.
func (kmgr *KafkaManager) ExecuteInitialize(context.Context) error {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	return nil
}

// Start component.
func (kmgr *KafkaManager) Start(ctx context.Context) error {
	return kmgr.lifecycle.Start(ctx)
}

// Lifecycle callback that runs startup logic.
func (kmgr *KafkaManager) ExecuteStart(context.Context) error {
	return nil
}

// Stop component.
func (kmgr *KafkaManager) Stop(ctx context.Context) error {
	return kmgr.lifecycle.Stop(ctx)
}

// Lifecycle callback that runs shutdown logic.
func (rdb *KafkaManager) ExecuteStop(context.Context) error {
	return nil
}

// Terminate component.
func (kmgr *KafkaManager) Terminate(ctx context.Context) error {
	return kmgr.lifecycle.Terminate(ctx)
}

// Lifecycle callback that runs termination logic.
func (kmgr *KafkaManager) ExecuteTerminate(context.Context) error {
	return nil
}
