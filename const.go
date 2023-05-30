package main

const (
	// DatetimeLayout for time strings from radiko
	DatetimeLayout = "20060102150405"
	// DefaultArea for radiko are
	DefaultArea = "JP13"
	// MaxRetryAttempts for BackOffDelay
	MaxRetryAttempts = 8
	// RetryDelaySecond for initial delay
	RetryDelaySecond = 60
	// DefaultInterval to fetch the programs
	DefaultInterval = "168h"
	// OutputDatetimeLayout for downloaded files
	OutputDatetimeLayout = "200601021504"
	// MaxConcurrents for doenloading programs
	MaxConcurrents = 64
	// MaxAttempts for downloading programs
	MaxAttempts = 4
	// TZTokyo for time location
	TZTokyo = "Asia/Tokyo"
)
