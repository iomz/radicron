package main

import (
	"log"
	"strings"
	"time"

	"github.com/yyoshiki41/go-radiko"
)

type Rules []*Rule

func (rs Rules) HasRuleWithoutStationID() bool {
	for _, r := range rs {
		if !r.HasStationID() {
			return true
		}
	}
	return false
}

func (rs Rules) HasRuleFor(stationID string) bool {
	for _, r := range rs {
		if r.StationID == stationID {
			return true
		}
	}
	return false
}

type Rule struct {
	Name      string `mapstructure:"name"`       // required
	Window    string `mapstructure:"window"`     // optional
	Title     string `mapstructure:"title"`      // required if pfm and keyword are unset
	Pfm       string `mapstructure:"pfm"`        // optional
	Keyword   string `mapstructure:"keyword"`    // optional
	AreaID    string `mapstructure:"area-id"`    // optional
	StationID string `mapstructure:"station-id"` // optional
}

// Match returns true if the rule matches the program
func (r *Rule) Match(stationID string, p radiko.Prog) bool {
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
	if r.HasStationID() && r.StationID != stationID {
		return false // skip mismatching rules for stationID
	}
	if r.HasTitle() && strings.Contains(p.Title, r.Title) {
		log.Printf("rule[%s] matched with title: '%s'", r.Name, p.Title)
		return true
	} else if r.HasPfm() && strings.Contains(p.Pfm, r.Pfm) {
		log.Printf("rule[%s] matched with pfm: '%s'", r.Name, p.Pfm)
		return true
	} else if r.HasKeyword() {
		// TODO: search for tags
		//for _, tag := range p.Tags
		if strings.Contains(p.Title, r.Keyword) {
			log.Printf("rule[%s] matched with title: '%s'", r.Name, p.Title)
			return true
		} else if strings.Contains(p.SubTitle, r.Keyword) {
			log.Printf("rule[%s] matched with sub-title: '%s'", r.Name, p.SubTitle)
			return true
		} else if strings.Contains(p.Desc, r.Keyword) {
			log.Printf("rule[%s] matched with desc: '%s'", r.Name, p.Desc)
			return true
		} else if strings.Contains(p.Pfm, r.Keyword) {
			log.Printf("rule[%s] matched with pfm: '%s'", r.Name, p.Pfm)
			return true
		} else if strings.Contains(p.Info, r.Keyword) {
			log.Printf("rule[%s] matched with info: '%s'", r.Name, p.Info)
			return true
		}
	}
	// both title and keyword are empty or not found
	return false
}

func (r *Rule) HasPfm() bool {
	return len(r.Pfm) != 0
}

func (r *Rule) HasKeyword() bool {
	return len(r.Keyword) != 0
}

func (r *Rule) HasStationID() bool {
	if len(r.StationID) == 0 ||
		r.StationID == "*" {
		return false
	}
	return true
}

func (r *Rule) HasTitle() bool {
	return len(r.Title) != 0
}

func (r *Rule) HasWindow() bool {
	return len(r.Window) != 0
}

func (r *Rule) SetName(name string) {
	r.Name = name
}
