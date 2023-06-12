package main

type Schedules []*Schedule

func (ss Schedules) HasDuplicate(prog Prog) bool {
	for _, s := range ss {
		// for the same title
		if s.Prog.Title == prog.Title {
			if s.Prog.StationID != prog.StationID {
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
	Prog Prog
}
