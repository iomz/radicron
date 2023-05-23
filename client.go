package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yyoshiki41/go-radiko"
)

// NewRadikoClient initializes the radiko client
// reimplemented from https://github.com/yyoshiki41/radigo/blob/main/client.go
func NewRadikoClient(ctx context.Context, areaID string) (*radiko.Client, error) {
	var client *radiko.Client
	var currentAreaID string
	var err error

	// check the current area
	currentAreaID, err = radiko.AreaID()
	if err != nil {
		return nil, err
	}

	client, err = radiko.New("")
	if err != nil {
		return nil, err
	}

	// When currentAreaID is not the same as the given areaID
	// we cannot download any programs
	// TODO: login with a premium account
	if areaID != "" && areaID != currentAreaID {
		client.SetAreaID(areaID)
		return nil, fmt.Errorf(
			"the specified area-id (%s) differs from your location's area-id (%s)",
			areaID,
			currentAreaID,
		)
	}

	// authorize the token
	_, err = client.AuthorizeToken(ctx)
	if err != nil {
		log.Println("failed to get auth_token")
		return nil, err
	}

	return client, nil
}
