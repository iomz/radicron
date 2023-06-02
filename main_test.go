package main

import (
	"testing"
	"time"

	"github.com/yyoshiki41/radigo"
)

func TestConfig(t *testing.T) {
	configure("test/config-test.yml")

	if AreaID != "JP13" {
		t.Errorf("%v => want %v", AreaID, "JP13")
	}

	if FileFormat != radigo.AudioFormatAAC {
		t.Errorf("%v => want %v", FileFormat, radigo.AudioFormatAAC)
	}

	if InitialDelay != time.Minute {
		t.Errorf("%v => want %v", InitialDelay, time.Minute)
	}

	if Interval != "168h" {
		t.Errorf("%v => want %v", Interval, "168h")
	}

	if MaxConcurrency != 64 {
		t.Errorf("%v => want %v", MaxConcurrency, 64)
	}

	if MaxRetryAttempts != 8 {
		t.Errorf("%v => want %v", MaxRetryAttempts, 8)
	}
}
