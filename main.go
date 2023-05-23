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
	location *time.Location
)

func run(ctx context.Context, client *radiko.Client) {
	// save the current time
	currentTime := time.Now().In(location)

	// fetch the programs to download
	programs := viper.GetStringMap("programs")

	// create the wait group for downloading
	var wg sync.WaitGroup
	for stationID := range programs {
		// TODO: check the stationID is valid
		// TODO: use * for all the avaialable stationID
		stationID = strings.ToUpper(stationID)

		// Fetch the weekly program
		weeklyPrograms, err := client.GetWeeklyPrograms(ctx, stationID)
		if err != nil {
			log.Fatal(err)
		}

		// cycle through the programs to record
		for ts := range viper.GetStringMap(fmt.Sprintf("programs.%s", stationID)) {
			lastSaved := viper.GetString(
				fmt.Sprintf("programs.%s.%s.last-saved", stationID, ts),
			)
			lastSavedTime, err := time.ParseInLocation(DatetimeLayout, lastSaved, location)
			if err != nil { // set the lastSavedTime as RFC3389
				lastSavedTime, _ = time.Parse(time.RFC3339, time.RFC3339)
			}
			title := viper.GetString(fmt.Sprintf("programs.%s.%s.title", stationID, ts))
			log.Printf("checking %s: [%s]%s", ts, stationID, title)

			for _, p := range weeklyPrograms[0].Progs.Progs {
				var start string
				// TODO: other matching methods
				if strings.Contains(p.Title, title) {
					title = p.Title
					start = p.Ft

					startTime, err := time.ParseInLocation(DatetimeLayout, start, location)
					if err != nil {
						log.Fatalf("invalid start time format '%s': %s", start, err)
					}

					if startTime.After(currentTime) { // if it is in the future, skip
						continue
					}

					// check in the found ts
					if startTime.After(lastSavedTime) {
						viper.Set(
							fmt.Sprintf("programs.%s.%s.last-saved", stationID, ts),
							start,
						)
						log.Printf("start downloading [%s]%s at %s", stationID, title, start)
					} else {
						log.Printf("skip [%s]%s at %s", stationID, title, start)
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
						log.Printf("the output file already exists: %s", output.AbsPath())
						log.Printf("skipping [%s]%s at %s", stationID, title, start)
						continue
					}

					// record the timeshift recording
					uri, err := client.TimeshiftPlaylistM3U8(ctx, stationID, startTime)
					if err != nil {
						log.Fatalf("failed to get playlist.m3u8: %s", err)
					}

					// detach the download job
					wg.Add(1)
					go DownloadProgram(ctx, &wg, uri, output)
				} // if found end
			} // ts in weeklyPrograms end
		} // ts in programs.[station-id] end
	} // stationID in the programs end

	// wait for all the downloading jobs
	wg.Wait()

	// save the last-saved dates
	if err := viper.WriteConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Set the config location
	cwd, _ := os.Getwd()
	conf := flag.String("c", "config.toml", "the config.[toml|yml] to use.")
	version := flag.Bool("v", false, "print version.")
	flag.Parse()

	// use the version from build
	if *version {
		bi, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", bi.Main.Version)
		os.Exit(0)
	}

	// to change the flags on the default logger
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting radiko-auto-downloader")

	// load config
	if *conf != "config.toml" {
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
	log.Printf("[config] area-id: %v", areaID)
	log.Printf("[config] interval: %v", interval)

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

	// put the runner to a scheduler
	s := gocron.NewScheduler(location)
	job, err := s.Every(interval).Do(run, ctx, client)
	if err != nil {
		log.Fatalf("job: %v, error: %v", job, err)
	}
	s.StartBlocking()
}
