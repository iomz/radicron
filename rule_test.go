package main

import "testing"

var ruletests = []struct {
	in  *Rule
	out bool
}{
	{
		&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "StationID"},
		true,
	},
	{
		&Rule{"", "", "", "", "", "", ""},
		false,
	},
}

func TestHasPfm(t *testing.T) {
	for _, tt := range ruletests {
		if tt.in.HasPfm() != tt.out {
			t.Errorf("(%v).HasPfm  => %v, want %v", tt.in, tt.in.HasPfm(), tt.out)
		}
	}
}

func TestHasKeyword(t *testing.T) {
	for _, tt := range ruletests {
		if tt.in.HasKeyword() != tt.out {
			t.Errorf("(%v).HasKeyword  => %v, want %v", tt.in, tt.in.HasKeyword(), tt.out)
		}
	}
}

func TestHasStationID(t *testing.T) {
	for _, tt := range ruletests {
		if tt.in.HasStationID() != tt.out {
			t.Errorf("(%v).HasStationID  => %v, want %v", tt.in, tt.in.HasStationID(), tt.out)
		}
	}
}

func TestHasTitle(t *testing.T) {
	for _, tt := range ruletests {
		if tt.in.HasTitle() != tt.out {
			t.Errorf("(%v).HasTitle  => %v, want %v", tt.in, tt.in.HasTitle(), tt.out)
		}
	}
}

func TestHasWindow(t *testing.T) {
	for _, tt := range ruletests {
		if tt.in.HasWindow() != tt.out {
			t.Errorf("(%v).HasWindow  => %v, want %v", tt.in, tt.in.HasWindow(), tt.out)
		}
	}
}

func TestSetName(t *testing.T) {
	r := &Rule{}
	r.SetName("Name")
	if r.Name != "Name" {
		t.Errorf("(%v)  => %v, want %v", r, r.Name, "Name")
	}
}

func TestHasRuleFor(t *testing.T) {
	var rulestests = []struct {
		in  Rules
		sid string
		out bool
	}{
		{
			Rules{
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "FMT"},
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "TBS"},
			},
			"FMT",
			true,
		},
		{
			Rules{
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "FMT"},
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "TBS"},
			},
			"MBS",
			false,
		},
	}
	for _, tt := range rulestests {
		res := tt.in.HasRuleForStationID(tt.sid)
		if tt.out != res {
			t.Errorf("(%v).HasRuleFor(\"%s\") => %v, want %v", tt.in, tt.sid, res, tt.out)
		}
	}
}

func TestHasRuleWithoutStationID(t *testing.T) {
	var rulestests = []struct {
		in  Rules
		out bool
	}{
		{
			Rules{
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", ""},
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "TBS"},
			},
			true,
		},
		{
			Rules{
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "FMT"},
				&Rule{"Name", "Window", "Title", "Pfm", "Keyword", "AreaID", "TBS"},
			},
			false,
		},
	}
	for _, tt := range rulestests {
		res := tt.in.HasRuleWithoutStationID()
		if tt.out != res {
			t.Errorf("(%v).HasRuleWithoutStationID() => %v, want %v", tt.in, res, tt.out)
		}
	}
}
