/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package kafka

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/devicechain-io/dc-microservice/core"
	"github.com/rs/zerolog/log"
	kafka "github.com/segmentio/kafka-go"
)

// Manages lifecycle of kafka interactions.
type KafkaManager struct {
	Microservice *core.Microservice

	oncreate  func(*KafkaManager) error
	readers   []*kafka.Reader
	writers   []*kafka.Writer
	lifecycle core.LifecycleManager
}

// Create a new kafka manager.
func NewKafkaManager(ms *core.Microservice, callbacks core.LifecycleCallbacks,
	oncreate func(*KafkaManager) error) *KafkaManager {
	kmgr := &KafkaManager{
		Microservice: ms,
	}

	kmgr.readers = make([]*kafka.Reader, 0)
	kmgr.writers = make([]*kafka.Writer, 0)
	kmgr.oncreate = oncreate

	// Create lifecycle manager.
	kfkaname := fmt.Sprintf("%s-%s", ms.FunctionalArea, "kafka")
	kmgr.lifecycle = core.NewLifecycleManager(kfkaname, kmgr, callbacks)
	return kmgr
}

// Get the kafka brokers url.
func (kmgr *KafkaManager) KafkaBrokersUrl() string {
	cfg := kmgr.Microservice.InstanceConfiguration.Infrastructure.Kafka
	return fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)
}

// Build topic name specific to instance/tenant.
func (kmgr *KafkaManager) NewScopedTopic(topic string) string {
	return fmt.Sprintf("%s.%s.%s", kmgr.Microservice.InstanceId, kmgr.Microservice.TenantId, topic)
}

// Build consumer group name specific to instant/tenant/microservice.
func (kmgr *KafkaManager) NewScopedConsumerGroup(topic string) string {
	return fmt.Sprintf("%s.%s.group-%s-%s", kmgr.Microservice.InstanceId, kmgr.Microservice.TenantId,
		kmgr.Microservice.FunctionalArea, topic)
}

// Create a topic if it doesn't already exist.
func (kmgr *KafkaManager) ValidateTopic(topic string) error {
	cfg := kmgr.Microservice.InstanceConfiguration.Infrastructure.Kafka
	conn, err := kafka.Dial("tcp", kmgr.KafkaBrokersUrl())
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     int(cfg.DefaultTopicPartitions),
			ReplicationFactor: int(cfg.DefaultTopicReplicationFactor),
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}

// Create a new kafka reader.
func (kmgr *KafkaManager) NewReader(groupId string, topic string) (*kafka.Reader, error) {
	err := kmgr.ValidateTopic(topic)
	if err != nil {
		return nil, err
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kmgr.KafkaBrokersUrl()},
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Info().Msg(fmt.Sprintf("Added new kafka reader on group '%s' for topic '%s'", groupId, topic))
	kmgr.readers = append(kmgr.readers, reader)
	return reader, nil
}

// Create a new kafka writer.
func (kmgr *KafkaManager) NewWriter(topic string, batchSize int, batchTimeout time.Duration) (*kafka.Writer, error) {
	err := kmgr.ValidateTopic(topic)
	if err != nil {
		return nil, err
	}
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kmgr.KafkaBrokersUrl()),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    batchSize,
		BatchTimeout: batchTimeout,
	}

	log.Info().Msg(fmt.Sprintf("Added new kafka writer for topic '%s'", topic))
	kmgr.writers = append(kmgr.writers, writer)
	return writer, nil
}

// Initialize component.
func (kmgr *KafkaManager) Initialize(ctx context.Context) error {
	return kmgr.lifecycle.Initialize(ctx)
}

// Lifecycle callback that runs initialization logic.
func (kmgr *KafkaManager) ExecuteInitialize(context.Context) error {
	url := kmgr.KafkaBrokersUrl()
	conn, err := kafka.Dial("tcp", url)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Info().Msg(fmt.Sprintf("Verified connectivity to kafka at '%s'", url))
	return nil
}

// Start component.
func (kmgr *KafkaManager) Start(ctx context.Context) error {
	return kmgr.lifecycle.Start(ctx)
}

// Lifecycle callback that runs startup logic.
func (kmgr *KafkaManager) ExecuteStart(context.Context) error {
	err := kmgr.oncreate(kmgr)
	if err != nil {
		return err
	}
	log.Info().Msg("Kafka component creation completed successfully.")
	return nil
}

// Stop component.
func (kmgr *KafkaManager) Stop(ctx context.Context) error {
	return kmgr.lifecycle.Stop(ctx)
}

// Lifecycle callback that runs shutdown logic.
func (kmgr *KafkaManager) ExecuteStop(context.Context) error {
	log.Info().Msg("Shutting down kafka writers.")
	for _, writer := range kmgr.writers {
		err := writer.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing kafka writer.")
		}
	}
	log.Info().Msg("Shutting down kafka readers.")
	for _, reader := range kmgr.readers {
		err := reader.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing kafka reader.")
		}
	}
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
