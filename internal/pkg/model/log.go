package model

// Log is the data structure for access log of short urls
type Log struct {
	// Id is the accessed short url's ID
	Id string `db:"id" json:"id"`

	// CreatedAt is the time in which the access occurred
	CreatedAt string `db:"created_at" json:"@timestamp"`

	// Device from which the url requested
	Device string `db:"device" json:"device"`

	// Browser from which user requested the url
	Browser string `db:"browser" json:"browser"`

	// RemoteAddr is the IP of the client
	RemoteAddr string `db:"remote_address" json:"remote_addr"`
}
