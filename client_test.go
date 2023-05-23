package main

import (
	"context"
	"testing"
)

func TestNewRadikoClient(t *testing.T) {
	areaID := "JP13"
	client, err := NewRadikoClient(context.Background(), areaID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if client.AreaID() != "OUT" && // the GitHub Action runner is outside Japan
		client.AreaID() != areaID {
		t.Errorf("client was configured with a wrong AreaID %v", client.AreaID())
	}
}
