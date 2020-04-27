package analyzer

// Config is singleton global instance of analyzer config
var Config config

type config struct {
	// App is the config of core service
	App struct {
		// Listen is of type host:port for http server to listen on
		Listen string

		// Secret is the key for JWT generation
		Secret string
	}

	// Database is the config of core database
	Database struct {
		// ConnectionString is the connStr of the database
		ConnectionString string `mapstructure:"connection_string"`

		// Dialect represents database type; ex: postgres, etc.
		Dialect string
	}

	// LogDB is the config of logs database
	LogDB struct {
		// ConnectionString is the connStr of the database
		ConnectionString string `mapstructure:"connection_string"`

		// Dialect represents database type; ex: postgres, etc.
		Dialect string
	}

	// Log is the logging config
	Log struct {
		// Level configures logging level
		Level string
	}
}
