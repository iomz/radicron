package radicron

import "testing"

var ruletests = []struct {
	in  *Rule
	out bool
}{
	{
		&Rule{"Name", "Title", []string{"sun"}, "Keyword", "Pfm", "StationID", "Window"},
		true,
	},
	{
		&Rule{"", "", []string{}, "", "", "", ""},
		false,
	},
}

func TestHasDoW(t *testing.T) {
	for _, tt := range ruletests {
		if tt.in.HasDoW() != tt.out {
			t.Errorf("(%v).HasDoW => %v, want %v", tt.in, tt.in.HasDoW(), tt.out)
		}
	}
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
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "FMT", "Window"},
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "TBS", "Window"},
			},
			"FMT",
			true,
		},
		{
			Rules{
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "FMT", "Window"},
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "TBS", "Window"},
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
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "", "Window"},
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "TBS", "Window"},
			},
			true,
		},
		{
			Rules{
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "FMT", "Window"},
				&Rule{"Name", "Title", []string{}, "Keyword", "Pfm", "TBS", "Window"},
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
