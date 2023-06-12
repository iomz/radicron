package radicron

import (
	"log"
	"time"
)

var (
	CurrentTime time.Time
	Location    *time.Location
)

func init() {
	var err error

	Location, err = time.LoadLocation(TZTokyo)
	if err != nil {
		log.Fatal(err)
	}
}
