package db

import (
	"clicktweak/internal/pkg/model"
)

// Log is an abstraction for log database
type Log interface {
	// GetStats returns the url statistics in the given time range
	//
	// returns (report, nil) on success an (nil, report) on failure
	GetStats(id, from, until string) (*model.Report, error)

	// Save batch inserts the array of log entries to the database
	Save(log []*model.Log, len int) error
}
