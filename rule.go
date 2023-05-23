package main

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
	Title     string `mapstructure:"title"`      // required if keyword is unset
	Keyword   string `mapstructure:"keyword"`    // optional
	AreaID    string `mapstructure:"area-id"`    // optional
	StationID string `mapstructure:"station-id"` // optional
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

func (r *Rule) SetName(name string) {
	r.Name = name
}
