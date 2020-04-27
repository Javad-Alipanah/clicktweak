package model

// Report represents single url analytics
type Report struct {
	Id string `json:"Id"`

	Clicks Stat `json:"clicks"`

	Visitors Stat `json:"visitors"`
}

// Stat contains statistics in detail
type Stat struct {
	Total      int            `json:"total"`
	PerBrowser map[string]int `json:"per_browser"`
	PerDevice  map[string]int `json:"per_device"`
}
