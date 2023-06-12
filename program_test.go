package radicron

import (
	"embed"
	"strings"
	"testing"
)

var (
	//go:embed test/weekly-program-test.xml
	WeeklyProgramTestXML embed.FS
)

func TestWeeklyProgramUnmarshal(t *testing.T) {
	xmlFile, err := WeeklyProgramTestXML.Open("test/weekly-program-test.xml")
	if err != nil {
		t.Error(err)
	}
	progs, err := decodeWeeklyProgram(xmlFile)
	if err != nil {
		t.Error(err)
	}
	if len(progs) != 1 {
		t.Errorf("unmarshal failed: %v", progs)
	}

	p := progs[0]
	var got, want string

	got = p.StationID
	want = "FMT"
	if got != want {
		t.Errorf("p.StationID => %v, want %v", got, want)
	}

	got = p.Ft
	want = "20230605130000"
	if got != want {
		t.Errorf("p.Ft => %v, want %v", got, want)
	}

	got = p.To
	want = "20230605145500"
	if got != want {
		t.Errorf("p.To => %v, want %v", got, want)
	}

	got = p.Title
	want = "山崎怜奈の誰かに話したかったこと。"
	if got != want {
		t.Errorf("p.Title => %v, want %v", got, want)
	}

	/*
		got = p.Desc
		want = ""
		if got != want {
			t.Errorf("p.Desc => %v, want %v", got, want)
		}

		got = p.Info
		want = ""
		if got != want {
			t.Errorf("p.Info => %v, want %v", got, want)
		}
	*/

	got = p.Pfm
	want = "山崎怜奈"
	if got != want {
		t.Errorf("p.Pfm => %v, want %v", got, want)
	}

	got = p.Genre.Personality
	want = "タレント"
	if got != want {
		t.Errorf("p.Genre.Personality => %v, want %v", got, want)
	}

	got = p.Genre.Program
	want = "トーク"
	if got != want {
		t.Errorf("p.Genre.Program => %v, want %v", got, want)
	}

	got = strings.Join(p.Tags, ",")
	want = "山崎怜奈,音楽との出会いが楽しめる,作業がはかどる,気分転換におすすめ,学生におすすめ"
	if got != want {
		t.Errorf("p.Tags => %v, want %v", got, want)
	}
}
