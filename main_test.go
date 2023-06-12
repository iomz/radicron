package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/yyoshiki41/go-radiko"
	"github.com/yyoshiki41/radigo"
)

func TestConfig(t *testing.T) {
	var err error
	Location, err = time.LoadLocation(TZTokyo)
	if err != nil {
		t.Error(err)
	}
	client, err := radiko.New("")
	if err != nil {
		log.Fatal(err)
	}
	ck := ContextKey("asset")
	asset, err := NewAsset(client)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), ck, asset)
	rules, err := reload(ctx, "test/config-test.yml")
	if err != nil {
		t.Error(err)
	}

	if asset.OutputFormat != radigo.AudioFormatAAC {
		t.Errorf("%v => want %v", asset.OutputFormat, radigo.AudioFormatAAC)
	}

	if len(rules) == 0 {
		t.Error("error parsing the rules")
	}
}
