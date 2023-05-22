package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/yyoshiki41/go-radiko"
	"github.com/yyoshiki41/radigo"
)

const (
	DatetimeLayout = "20060102150405"
	MaxAttempts    = 4
	MaxConcurrents = 64
	TZTokyo        = "Asia/Tokyo"
)

// reimplement some internal functions from
// https://github.com/yyoshiki41/radigo/blob/main/internal/download.go

var sem = make(chan struct{}, MaxConcurrents)

func bulkDownload(list []string, output string) error {
	var errFlag bool
	var wg sync.WaitGroup

	for _, v := range list {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			var err error
			for i := 0; i < MaxAttempts; i++ {
				sem <- struct{}{}
				err = downloadLink(link, output)
				<-sem
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Printf("failed to download: %s", err)
				errFlag = true
			}
		}(v)
	}
	wg.Wait()

	if errFlag {
		return errors.New("lack of aac files")
	}
	return nil
}

func downloadLink(link, output string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, fileName := filepath.Split(link)
	file, err := os.Create(filepath.Join(output, fileName))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}

// DownloadProgram manages the download for the given program
// in a go routine and notify the wg when finished
func DownloadProgram(
	ctx context.Context, // the context for the request
	wg *sync.WaitGroup, // the wg to notify
	uri string, // the m3u8 URI for the program
	output *radigo.OutputConfig, // the file configuration
) {
	defer wg.Done()
	chunklist, err := radiko.GetChunklistFromM3U8(uri)
	if err != nil {
		log.Fatalf("failed to get chunklist: %s", err)
	}

	aacDir, err := output.TempAACDir()
	if err != nil {
		log.Fatalf("failed to create the aac dir: %s", err)
	}
	defer os.RemoveAll(aacDir) // clean up

	if err := bulkDownload(chunklist, aacDir); err != nil {
		log.Fatalf("failed to download aac files: %s", err)
	}

	concatedFile, err := radigo.ConcatAACFilesFromList(ctx, aacDir)
	if err != nil {
		log.Fatalf("failed to concat aac files: %s", err)
	}

	err = os.Rename(concatedFile, output.AbsPath())
	if err != nil {
		log.Fatalf(
			"Failed to output a result file: %s", err)
	}
	log.Printf("+file saved: %s", output.AbsPath())
}

// NewRadikoClient initializes the radiko client
// reimplemented from https://github.com/yyoshiki41/radigo/blob/main/client.go
func NewRadikoClient(ctx context.Context, areaID string) (*radiko.Client, error) {
	var client *radiko.Client
	var currentAreaID string
	var err error

	// check the current area
	currentAreaID, err = radiko.AreaID()
	if err != nil {
		return nil, err
	}

	client, err = radiko.New("")
	if err != nil {
		return nil, err
	}

	// When currentAreaID is not the same as the given areaID
	// we cannot download any programs
	// TODO: login with a premium account
	if areaID != "" && areaID != currentAreaID {
		return nil, fmt.Errorf(
			"the specified area-id (%s) differs from your location's area-id (%s)",
			areaID,
			currentAreaID,
		)
	}

	// authorize the token
	_, err = client.AuthorizeToken(ctx)
	if err != nil {
		log.Println("failed to get auth_token")
		return nil, err
	}

	return client, nil
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting radiko-auto-recorder")

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
	viper.SetDefault("areaID", "JP13")
	// set the default interval as weekly
	viper.SetDefault("check-interval-days", 7)

	// get the global config paramters
	var (
		areaID            = viper.GetString("area-id")
		checkIntervalDays = viper.GetInt("check-interval-days")
	)
	log.Printf("[config] area-id: %v", areaID)
	log.Printf("[config] check-interval-days: %v (days)", checkIntervalDays)

	// set the interval
	checkInterval := time.Hour * 24 * time.Duration(checkIntervalDays)

	// initialize radiko client with context
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	client, err := NewRadikoClient(ctx, areaID)
	if err != nil {
		log.Fatal(err)
	}

	location, err := time.LoadLocation(TZTokyo)
	if err != nil {
		panic(err)
	}

	for {
		// save the current time and the time after the check
		currentTime := time.Now().In(location)

		// fetch the programs to download
		programs := viper.GetStringMap("programs")

		// create the wait group for downloading
		var wg sync.WaitGroup
		for stationID := range programs {
			// TODO: check the stationID is valid
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
				keyword := viper.GetString(fmt.Sprintf("programs.%s.%s.keyword", stationID, ts))
				log.Printf("checking %s: [%s]%s", ts, stationID, keyword)

				for _, p := range weeklyPrograms[0].Progs.Progs {
					var title, start string
					if strings.Contains(p.Title, keyword) {
						title = p.Title
						start = p.Ft

						startTime, err := time.ParseInLocation(DatetimeLayout, start, location)
						if err != nil {
							log.Fatalf("invalid start time format '%s': %s", start, err)
						}

						if startTime.After(currentTime) { // if it is in the future, skip
							continue
						}
						log.Printf("found [%s]%s at %s", stationID, title, start)

						// check in the found ts
						if startTime.After(lastSavedTime) {
							viper.Set(
								fmt.Sprintf("programs.%s.%s.last-saved", stationID, ts),
								start,
							)
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
		err = viper.WriteConfig()
		if err != nil {
			log.Fatal(err)
		}

		// sleep until the next time
		diff := currentTime.Add(checkInterval).Sub(currentTime)
		log.Printf("sleeping %v seconds...", diff.Seconds())
		time.Sleep(diff)
	} // forever end
}
