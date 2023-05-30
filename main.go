package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"github.com/yyoshiki41/go-radiko"
	"github.com/yyoshiki41/radigo"
)

var (
	areaID      string         // store the area-id
	currentTime time.Time      // store the current location
	fileFormat  string         // store the fileFormat to save
	interval    string         // store the checking interval
	location    *time.Location // store the current location
	stations    []string       // store the available stations
)

func run(ctx context.Context, client *radiko.Client, interval string) {
	// log the current time
	currentTime = time.Now().In(location)

	// refresh the rules from the file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config file: %s \n", err)
	}
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

		// fetch the weekly program
		log.Printf("fetching the %s program", stationID)
		weeklyPrograms, err := client.GetWeeklyPrograms(ctx, stationID)
		if err != nil {
			log.Printf("failed to fetch the %s program: %v", stationID, err)
			continue
		}

		// iterate through the rules to download
		for _, r := range rules {
			for _, p := range weeklyPrograms[0].Progs.Progs {
				if r.Match(stationID, p) {
					err = Download(ctx, &wg, client, p, stationID)
					if err != nil {
						log.Printf("downlod faild: %s", err)
					}
				}
			} // weeklyPrograms for stationID
		} // rules
	} // stations

	// wait for all the downloading jobs
	log.Println("waiting for all the downloads to complete")
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
	enableDebug := flag.Bool("d", false, "enable debug mode.")
	version := flag.Bool("v", false, "print version.")
	flag.Parse()

	// use the version from build
	if *version {
		bi, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", bi.Main.Version)
		os.Exit(0)
	}

	// to change the flags on the default logger
	if *enableDebug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	log.Println("starting radicron")

	// check ${RADIGO_HOME}
	if len(os.Getenv("RADIGO_HOME")) == 0 {
		os.Setenv("RADIGO_HOME", filepath.Join(cwd, "./downloads"))
	}

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
	viper.SetDefault("area-id", DefaultArea)
	// set the default interval as weekly
	viper.SetDefault("interval", DefaultInterval)
	// set the default file-format as aac
	viper.SetDefault("file-format", radigo.AudioFormatAAC)

	// get the global config parameters
	areaID = viper.GetString("area-id")
	fileFormat = viper.GetString("file-format")
	interval = viper.GetString("interval")

	// TODO: verify the area-id
	// check if the interval is invalid or is too short
	intervalDuration, err := time.ParseDuration(interval)
	if err != nil || intervalDuration < time.Hour {
		log.Fatalf("invalid interval: %s, setting to %v", interval, DefaultInterval)
		interval = DefaultInterval
	}
	// check the output file format
	if fileFormat != radigo.AudioFormatAAC &&
		fileFormat != radigo.AudioFormatMP3 {
		log.Fatalf("unsupported audio format: %s", fileFormat)
	}

	log.Printf("[config] area-id: %s", areaID)
	log.Printf("[config] file-format: %s", fileFormat)
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
