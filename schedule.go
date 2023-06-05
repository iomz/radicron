package main

import "github.com/yyoshiki41/go-radiko"

type Schedules []*Schedule

func (ss Schedules) HasDuplicate(stationID string, prog radiko.Prog) bool {
	for _, s := range ss {
		// for the same title
		if s.Prog.Title == prog.Title {
			if s.StationID != stationID {
				// it has the same title at a different station
				return true
			} else if s.Prog.Ft == prog.Ft {
				// it already has the exact schedle
				return true
			}
		}
	}
	return false
}

type Schedule struct {
	StationID string
	Prog      radiko.Prog
}
