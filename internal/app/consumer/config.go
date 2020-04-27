package consumer

import "time"

// Config is singleton global instance of consumer config
var Config config

type config struct {
	// App is the config of the consumer
	APP struct {
		// Batch sets the number of elements per insert
		Batch int

		// Timeout for batch insert
		Timeout time.Duration
	}

	// Kafka is the properties of kafka service
	Kafka struct {
		// Topic to consume from
		Topic string

		// Host is address of kafka
		Host string

		// Group is the consumer group id
		Group string
	}

	// LogDB is the properties of logs database
	LogDB struct {
		// Dialect is database name ex: postgres, clickhouse, etc.
		Dialect string

		// ConnectionString is the database uri
		ConnectionString string `mapstructure:"connection_string"`
	} `mapstructure:"logdb"`

	// Log is the logging config
	Log struct {
		// Level configures logging level
		Level string
	}
}
