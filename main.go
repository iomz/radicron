package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"github.com/yyoshiki41/go-radiko"
	"github.com/yyoshiki41/radigo"
)

var (
	location *time.Location // store the current location
	stations []string       // store the available stations
)

func run(ctx context.Context, client *radiko.Client, interval string) {
	// save the current time
	currentTime := time.Now().In(location)

	// re-fetch the rules
	rules := Rules{}
	for name := range viper.GetStringMap("rules") {
		rule := &Rule{}
		err := viper.UnmarshalKey(fmt.Sprintf("rules.%s", name), rule)
		if err != nil {
			log.Fatal(err)
		}
		rule.SetName(name)
		rules = append(rules, rule)
	}
	// create the wait group for downloading
	var wg sync.WaitGroup
	for _, stationID := range stations {
		if !rules.HasRuleWithoutStationID() && // search all stations
			!rules.HasRuleFor(stationID) { // search this station
			continue
		}

		// Fetch the weekly program
		log.Printf("fetching the %s program", stationID)
		weeklyPrograms, err := client.GetWeeklyPrograms(ctx, stationID)
		if err != nil {
			log.Printf("failed to fetch the %s program: %v", stationID, err)
			continue
		}

		// cycle through the programs to record
		for _, r := range rules {
			if r.HasStationID() && r.StationID != stationID {
				continue // skip unspecified stations
			}
			log.Printf("searching for [%s] title='%s' keyword='%s'", r.Name, r.Title, r.Keyword)

			for _, p := range weeklyPrograms[0].Progs.Progs {
				var start, title string
				if r.HasTitle() && strings.Contains(p.Title, r.Title) {
					title = p.Title
					start = p.Ft
				} else if r.HasKeyword() {
					if strings.Contains(p.Title, r.Keyword) ||
						strings.Contains(p.SubTitle, r.Keyword) ||
						strings.Contains(p.Desc, r.Keyword) ||
						strings.Contains(p.Pfm, r.Keyword) ||
						strings.Contains(p.Info, r.Keyword) {
						title = p.Title
						start = p.Ft
					} else { // not found
						continue
					}
				} else { // both title and keyword are empty
					continue
				}

				startTime, err := time.ParseInLocation(DatetimeLayout, start, location)
				if err != nil {
					log.Fatalf("invalid start time format '%s': %s", start, err)
				}

				if startTime.After(currentTime) { // if it is in the future, skip
					continue
				}

				output, err := radigo.NewOutputConfig(
					fmt.Sprintf(
						"%s-%s_%s",
						startTime.In(location).Format(DatetimeLayout),
						stationID,
						title,
					),
					radigo.AudioFormatAAC,
				)
				if err != nil {
					log.Fatalf("failed to configure output: %s", err)
				}

				if err := output.SetupDir(); err != nil {
					log.Fatalf("failed to setup the output dir: %s", err)
				}

				if output.IsExist() {
					log.Printf("skip [%s]%s at %s", stationID, title, start)
					log.Printf("the output file already exists: %s", output.AbsPath())
					continue
				}
				log.Printf("start downloading [%s]%s at %s", stationID, title, start)

				// record the timeshift recording
				uri, err := client.TimeshiftPlaylistM3U8(ctx, stationID, startTime)
				if err != nil {
					log.Fatalf("failed to get playlist.m3u8: %s", err)
				}

				// detach the download job
				wg.Add(1)
				go DownloadProgram(ctx, &wg, uri, output)
			} // programs
		} // rules
	} // stations

	// wait for all the downloading jobs
	wg.Wait()

	delta, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("fetching completed – sleeping until %v", currentTime.Add(delta))
}

func main() {
	// Set the config location
	cwd, _ := os.Getwd()
	conf := flag.String("c", "config.yml", "the config.yml to use.")
	version := flag.Bool("v", false, "print version.")
	flag.Parse()

	// use the version from build
	if *version {
		bi, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", bi.Main.Version)
		os.Exit(0)
	}

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting radiko-auto-downloader")

	// load config
	if *conf != "config.yml" {
		configPath, err := filepath.Abs(*conf)
		if err != nil {
			panic(err)
		}
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(cwd)
	}

	// read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config file: %s \n", err)
	}

	// set the default area_id
	viper.SetDefault("area-id", "JP13")
	// set the default interval as weekly
	viper.SetDefault("interval", "168h")

	// get the global config paramters
	var (
		areaID   = viper.GetString("area-id")
		interval = viper.GetString("interval")
	)
	log.Printf("[config] area-id: %s", areaID)
	log.Printf("[config] interval: %s", interval)

	// initialize radiko client with context
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	client, err := NewRadikoClient(ctx, areaID)
	if err != nil {
		log.Fatal(err)
	}

	// let us use this in Japan
	location, err = time.LoadLocation(TZTokyo)
	if err != nil {
		log.Fatal(err)
	}

	// fetch the available stations from all the regions
	region, err := GetRegion()
	if err != nil {
		log.Fatal(err)
	}
	for _, ss := range region.Region {
		for _, s := range ss.Stations {
			if s.AreaID == areaID {
				stations = append(stations, s.ID)
			}
		}
	}
	log.Printf("available stations in %s: %q", areaID, stations)

	// put the runner to a scheduler
	s := gocron.NewScheduler(location)
	job, err := s.Every(interval).Do(run, ctx, client, interval)
	if err != nil {
		log.Fatalf("job: %v, error: %v", job, err)
	}
	s.StartBlocking()
}
