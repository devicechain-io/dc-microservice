/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package config

// Redis configuration parameters
type RedisConfiguration struct {
	Hostname string
	Port     int32
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

// Keycloak connectivity configuration
type KeycloakConfiguration struct {
	Hostname string
	Port     uint32
}

// Infrastructure configuration section
type InfrastructureConfiguration struct {
	Redis    RedisConfiguration
	Kafka    KafkaConfiguration
	Metrics  MetricsConfiguration
	Keycloak KeycloakConfiguration
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

// Creates the default instance configuration
func NewDefaultInstanceConfiguration() *InstanceConfiguration {
	return &InstanceConfiguration{
		Infrastructure: InfrastructureConfiguration{
			Redis: RedisConfiguration{
				Hostname: "dc-redis-master.dc-system",
				Port:     6379,
			},
			Kafka: KafkaConfiguration{
				Hostname:                      "dc-kafka-kafka-bootstrap.dc-system",
				Port:                          9092,
				DefaultTopicPartitions:        4,
				DefaultTopicReplicationFactor: 1,
			},
			Metrics: MetricsConfiguration{
				Enabled:  true,
				HttpPort: 9090,
			},
			Keycloak: KeycloakConfiguration{
				Hostname: "dc-keycloak.dc-system",
				Port:     8080,
			},
		},
		Persistence: PersistenceConfiguration{
			Rdb: RdbConfiguration{
				Type: "postgres95",
				Configuration: map[string]interface{}{
					"hostname":       "dc-postgresql.dc-system",
					"port":           5432,
					"maxConnections": 5,
					"username":       "devicechain",
					"password":       "devicechain",
				},
			},
			Tsdb: TsdbConfiguration{
				Type: "influxdb",
				Configuration: map[string]interface{}{
					"hostname":     "dc-influxdb.dc-system",
					"port":         8086,
					"databaseName": "tenant_${tenant.id}",
				},
			},
		},
	}
}
