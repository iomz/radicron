package radicron

import (
	"log"
	"strings"
	"time"
)

type Rules []*Rule

func (rs Rules) HasMatch(stationID string, p *Prog) bool {
	for _, r := range rs {
		if r.Match(stationID, p) {
			return true
		}
	}
	return false
}

func (rs Rules) HasRuleWithoutStationID() bool {
	for _, r := range rs {
		if !r.HasStationID() {
			return true
		}
	}
	return false
}

func (rs Rules) HasRuleForStationID(stationID string) bool {
	for _, r := range rs {
		if r.StationID == stationID {
			return true
		}
	}
	return false
}

type Rule struct {
	Name      string   `mapstructure:"name"`       // required
	Title     string   `mapstructure:"title"`      // required if pfm and keyword are unset
	DoW       []string `mapstructure:"dow"`        // optional
	Keyword   string   `mapstructure:"keyword"`    // optional
	Pfm       string   `mapstructure:"pfm"`        // optional
	StationID string   `mapstructure:"station-id"` // optional
	Window    string   `mapstructure:"window"`     // optional
}

// Match returns true if the rule matches the program
// 1. check the Window filter
// 2. check the DoW filter
// 3. check the StationID filter
// TODO: reduce the complexity
func (r *Rule) Match(stationID string, p *Prog) bool { //nolint:gocyclo
	// 1. check Window
	if r.HasWindow() {
		startTime, err := time.ParseInLocation(DatetimeLayout, p.Ft, Location)
		if err != nil {
			log.Printf("invalid start time format '%s': %s", p.Ft, err)
			return false
		}
		fetchWindow, err := time.ParseDuration(r.Window)
		if err != nil {
			log.Printf("parsing [%s].past failed: %v (using 24h)", r.Name, err)
			fetchWindow = time.Hour * 24
		}
		if startTime.Add(fetchWindow).Before(CurrentTime) {
			return false // skip before the fetch window
		}
	}
	// 2. check DoW
	if r.HasDoW() {
		dow := map[string]time.Weekday{
			"sun": time.Sunday,
			"mon": time.Monday,
			"tue": time.Tuesday,
			"wed": time.Wednesday,
			"thu": time.Thursday,
			"fri": time.Friday,
			"sat": time.Saturday,
		}
		st, _ := time.ParseInLocation(DatetimeLayout, p.Ft, Location)
		dowMatch := false
		for _, d := range r.DoW {
			if st.Weekday() == dow[strings.ToLower(d)] {
				dowMatch = true
			}
		}
		if !dowMatch {
			return false
		}
	}
	// 3. check StationID
	if r.HasStationID() && r.StationID != stationID {
		return false // skip mismatching rules for stationID
	}
	// 4. Match
	if r.HasTitle() && strings.Contains(p.Title, r.Title) {
		log.Printf("rule[%s] matched with title: '%s'", r.Name, p.Title)
		return true
	} else if r.HasPfm() && strings.Contains(p.Pfm, r.Pfm) {
		log.Printf("rule[%s] matched with pfm: '%s'", r.Name, p.Pfm)
		return true
	} else if r.HasKeyword() {
		if strings.Contains(p.Title, r.Keyword) {
			log.Printf("rule[%s] matched with title: '%s'", r.Name, p.Title)
			return true
		} else if strings.Contains(p.Pfm, r.Keyword) {
			log.Printf("rule[%s] matched with pfm: '%s'", r.Name, p.Pfm)
			return true
		} else if strings.Contains(p.Info, r.Keyword) {
			log.Printf("rule[%s] matched with info: \n%s", r.Name, strings.ReplaceAll(p.Info, "\n", ""))
			return true
		} else if strings.Contains(p.Desc, r.Keyword) {
			log.Printf("rule[%s] matched with desc: '%s'", r.Name, strings.ReplaceAll(p.Desc, "\n", ""))
			return true
		}
		for _, tag := range p.Tags {
			if strings.Contains(tag, r.Keyword) {
				log.Printf("rule[%s] matched with tag: '%s'", r.Name, tag)
				return true
			}
		}
	}
	// both title and keyword are empty or not found
	return false
}

func (r *Rule) HasDoW() bool {
	return len(r.DoW) > 0
}

func (r *Rule) HasPfm() bool {
	return r.Pfm != ""
}

func (r *Rule) HasKeyword() bool {
	return r.Keyword != ""
}

func (r *Rule) HasStationID() bool {
	if r.StationID == "" ||
		r.StationID == "*" {
		return false
	}
	return true
}

func (r *Rule) HasTitle() bool {
	return r.Title != ""
}

func (r *Rule) HasWindow() bool {
	return r.Window != ""
}

func (r *Rule) SetName(name string) {
	r.Name = name
}
