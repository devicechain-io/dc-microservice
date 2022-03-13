/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package config

import "fmt"

// Redis configuration parameters
type RedisConfiguration struct {
	Hostname        string
	Port            int32
	NodeCount       int32
	MasterGroupName string
}

// Kafka configuration parameters
type KafkaConfiguration struct {
	Hostname                      string
	Port                          uint32
	DefaultTopicPartitions        uint32
	DefaultTopicReplicationFactor uint32
}

// Prometheus metrics configuration
type MetricsConfiguration struct {
	Enabled  bool
	HttpPort int32
}

// gRPC connectivity configuration
type GrpcConfiguration struct {
	MaxRetryCount         uint32
	InitialBackoffSeconds uint32
	MaxBackoffSeconds     uint32
	BackoffMultiplier     float32
	ResolveFQDN           bool
}

// Infrastructure configuration section
type InfrastructureConfiguration struct {
	Redis   RedisConfiguration
	Kafka   KafkaConfiguration
	Metrics MetricsConfiguration
	Grpc    GrpcConfiguration
}

// Relational database configuration
type RdbConfiguration struct {
	Type          string
	Configuration map[string]interface{}
}

// Time series database configuration
type TsdbConfiguration struct {
	Type          string
	Configuration map[string]interface{}
}

// Configuration of persistence stores
type PersistenceConfiguration struct {
	Rdb  RdbConfiguration
	Tsdb TsdbConfiguration
}

// Instance-level configuration settings
type InstanceConfiguration struct {
	Infrastructure InfrastructureConfiguration
	Persistence    PersistenceConfiguration
}

// Information required to create a resource file
type ConfigurationResource struct {
	Name    string
	Content []byte
}

type ConfigurationResourceProvider interface {
	GetConfigurationResources() ([]ConfigurationResource, error)
}

// Get instance configuration CRs that should be created in tooling
func GetConfigurationResources() ([]ConfigurationResource, error) {
	resources := make([]ConfigurationResource, 0)

	name := "dcic-default"
	config := NewDefaultInstanceConfiguration()
	content, err := GenerateInstanceConfig(name, config)
	if err != nil {
		return nil, err
	}
	dcidefault := ConfigurationResource{
		Name:    fmt.Sprintf("%s_%s", "core.devicechain.io", name),
		Content: content,
	}

	resources = append(resources, dcidefault)
	return resources, nil
}

// Creates the default instance configuration
func NewDefaultInstanceConfiguration() *InstanceConfiguration {
	return &InstanceConfiguration{
		Infrastructure: InfrastructureConfiguration{
			Redis: RedisConfiguration{
				Hostname:        "dc-infrastructure-redis-ha-announce",
				Port:            26379,
				NodeCount:       3,
				MasterGroupName: "devicechain",
			},
			Kafka: KafkaConfiguration{
				Hostname:                      "dc-kafka-kafka-bootstrap",
				Port:                          9092,
				DefaultTopicPartitions:        8,
				DefaultTopicReplicationFactor: 3,
			},
			Metrics: MetricsConfiguration{
				Enabled:  true,
				HttpPort: 9090,
			},
			Grpc: GrpcConfiguration{
				MaxRetryCount:         6,
				InitialBackoffSeconds: 10,
				MaxBackoffSeconds:     600,
				BackoffMultiplier:     1.5,
				ResolveFQDN:           false,
			},
		},
		Persistence: PersistenceConfiguration{
			Rdb: RdbConfiguration{
				Type: "postgres95",
				Configuration: map[string]interface{}{
					"hostname":       "dc-postgresql.sitewhere-system",
					"port":           5432,
					"maxConnections": 5,
					"username":       "syncope",
					"password":       "syncope",
				},
			},
			Tsdb: TsdbConfiguration{
				Type: "influxdb",
				Configuration: map[string]interface{}{
					"hostname":     "dc-infrastructure-influxdb.sitewhere-system",
					"port":         8086,
					"databaseName": "tenant_${tenant.id}",
				},
			},
		},
	}
}
