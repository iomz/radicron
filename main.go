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
	AreaID            string         // radiko's area-id
	AvailableStations []string       // the available stations
	CurrentTime       time.Time      // time in the location
	FileFormat        string         // the file format to save
	InitialDelay      time.Duration  // the initial delay to back off from
	Interval          string         // the checking interval
	Location          *time.Location // the current location
	MaxConcurrency    int            // for downloading
	MaxRetryAttempts  uint           // for downloading
)

func configure(filename string) {
	cwd, _ := os.Getwd()

	// check ${RADIGO_HOME}
	if len(os.Getenv("RADIGO_HOME")) == 0 {
		os.Setenv("RADIGO_HOME", filepath.Join(cwd, "./downloads"))
	}

	// load params from a config file
	if filename != "config.yml" && filename != "config.toml" {
		configPath, err := filepath.Abs(filename)
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
		log.Fatalf("error reading config: %s \n", err)
	}

	// set the default area_id
	viper.SetDefault("area-id", DefaultArea)
	// set the default file-format as aac
	viper.SetDefault("file-format", radigo.AudioFormatAAC)
	// set the default initial retry delay as 60 seconds
	viper.SetDefault("initial-delay", DefaultInitialDelaySeconds)
	// set the default interval as weekly
	viper.SetDefault("interval", DefaultInterval)
	// set the default max concurrency as 64
	viper.SetDefault("max-concurrency", DefaultMaxConcurrency)
	// set the default max retries as 8
	viper.SetDefault("max-retry", DefaultMaxRetryAttempts)

	// get the global config parameters
	AreaID = viper.GetString("area-id")
	FileFormat = viper.GetString("file-format")
	InitialDelay = time.Second * time.Duration(viper.GetInt("initial-delay"))
	Interval = viper.GetString("interval")
	MaxConcurrency = viper.GetInt("max-concurrency")
	MaxRetryAttempts = viper.GetUint("max-retry")

	// check if the interval is invalid or is too short
	intervalDuration, err := time.ParseDuration(Interval)
	if err != nil || intervalDuration < time.Hour {
		log.Fatalf("invalid interval: %s, setting to %v", Interval, DefaultInterval)
		Interval = DefaultInterval
	}
	// check the output file format
	if FileFormat != radigo.AudioFormatAAC &&
		FileFormat != radigo.AudioFormatMP3 {
		log.Fatalf("unsupported audio format: %s", FileFormat)
	}

	log.Printf("[config] area-id: %s", AreaID)
	log.Printf("[config] file-format: %s", FileFormat)
	log.Printf("[config] initial-delay: %v", InitialDelay)
	log.Printf("[config] interval: %s", Interval)
	log.Printf("[config] max-concurrency: %v", MaxConcurrency)
	log.Printf("[config] max-retry: %v", MaxRetryAttempts)

	// radiko is in Japan
	Location, err = time.LoadLocation(TZTokyo)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, client *radiko.Client, interval string) {
	// log the current time
	CurrentTime = time.Now().In(Location)

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
	for _, stationID := range AvailableStations {
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
	log.Printf("fetching completed – sleeping until %v", CurrentTime.Add(delta))
}

func main() {
	// Set the config location
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

	// load config params
	configure(*conf)

	// initialize radiko client with context
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	client, err := NewRadikoClient(ctx, AreaID)
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
			if s.AreaID == AreaID {
				AvailableStations = append(AvailableStations, s.ID)
			}
		}
	}
	log.Printf("available stations in %s: %q", AreaID, AvailableStations)

	// put the runner to a scheduler
	s := gocron.NewScheduler(Location)
	job, err := s.Every(Interval).Do(run, ctx, client, Interval)
	if err != nil {
		log.Fatalf("job: %v, error: %v", job, err)
	}
	s.StartBlocking()
}
