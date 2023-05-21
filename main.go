package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	//"regexp"
	"runtime/debug"
	"strings"
	//"time"

	"github.com/spf13/viper"
	"github.com/yyoshiki41/go-radiko"
)

func main() {
	// Set the config location
	cwd, _ := os.Getwd()
	conf := flag.String("c", "config.toml", "The config.[toml|yml] to use.")
	version := flag.Bool("v", false, "Print version.")
	flag.Parse()

	// use the version from build
	if *version {
		bi, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", bi.Main.Version)
		os.Exit(0)
	}

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("--- Starting radiko-crawler")

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
	if err := viper.ReadInConfig(); err != nil { // handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	// initialize radiko
	client, err := radiko.New("")
	if err != nil {
		log.Fatal(err)
	}
	// enter with timer
	// get programs data
	stations, err := client.GetNowPrograms(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, st := range stations {
		for _, p := range st.Scd.Progs.Progs {
			log.Printf("[%v] %v %v", st.ID, p.Ft, p.Title)
			// check the programs
			programs := viper.GetStringMap("programs")
			for program := range programs {
				lastSaved := viper.GetString(fmt.Sprintf("programs.%s.last-saved", program))
				station := viper.GetString(fmt.Sprintf("programs.%s.station", program))
				keyword := viper.GetString(fmt.Sprintf("programs.%s.keyword", program))
				log.Printf("Checking %s: %s", program, keyword)

				if st.ID != station {
					continue
				}
				var title, start string
				if strings.Contains(p.Title, keyword) {
					title = p.Title
					start = p.Ft
				}
				log.Printf("Found %s at [%s] %s", title, station, start)
				// TODO: dispatch recording

				if start != lastSaved {
					viper.Set(fmt.Sprintf("programs.%s.last-saved", program), start)
				}
				// write out the config
				viper.WriteConfig()
			} // programs end
		}
	}
}
