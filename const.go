package main

const (
	// DatetimeLayout for time strings from radiko
	DatetimeLayout = "20060102150405"
	// DefaultArea for radiko are
	DefaultArea = "JP13"
	// RetryDelaySecond for initial delay
	DefaultInitialDelaySeconds = 60
	// DefaultInterval to fetch the programs
	DefaultInterval = "168h"
	// DefaultMaxConcurrents
	DefaultMaxConcurrency = 64
	// MaxRetryAttempts for BackOffDelay
	DefaultMaxRetryAttempts = 8
	// OutputDatetimeLayout for downloaded files
	OutputDatetimeLayout = "200601021504"
	// TZTokyo for time location
	TZTokyo = "Asia/Tokyo"
)
