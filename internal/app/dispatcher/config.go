package dispatcher

// Config is singleton global instance of dispatcher config
var Config config

type config struct {
	// App is the config of core service
	App struct {
		// Listen is of type host:port for http server to listen on
		Listen string
	}

	// Forwarder is the config of log forwarder
	Forwarder struct {
		// Url is the address of log forwarder
		Url string

		// Workers is the number of workers for log handling
		Workers int

		// ChannelSize is the size of log channel
		ChannelSize int `mapstructure:"channel_size"`
	}

	// Database is the config of core database
	Database struct {
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
