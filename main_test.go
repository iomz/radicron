package main

import (
	"testing"

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

	if Interval != "168h" {
		t.Errorf("%v => want %v", Interval, "168h")
	}
}
