package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/bogem/id3v2"
	"github.com/yyoshiki41/go-radiko"
	"github.com/yyoshiki41/radigo"
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
	prog radiko.Prog, // the program metadata
	uri string, // the m3u8 URI for the program
	output *radigo.OutputConfig, // the file configuration
) {
	defer wg.Done()

	chunklist, err := radiko.GetChunklistFromM3U8(uri)
	if err != nil {
		log.Printf("failed to get chunklist: %s", err)
		return
	}

	aacDir, err := output.TempAACDir()
	if err != nil {
		log.Printf("failed to create the aac dir: %s", err)
		return
	}
	defer os.RemoveAll(aacDir) // clean up

	if err := bulkDownload(chunklist, aacDir); err != nil {
		log.Printf("failed to download aac files: %s", err)
		return
	}

	concatedFile, err := radigo.ConcatAACFilesFromList(ctx, aacDir)
	if err != nil {
		log.Printf("failed to concat aac files: %s", err)
		return
	}

	switch output.AudioFormat() {
	case radigo.AudioFormatAAC:
		err = os.Rename(concatedFile, output.AbsPath())
	case radigo.AudioFormatMP3:
		err = radigo.ConvertAACtoMP3(ctx, concatedFile, output.AbsPath())
	default:
		log.Fatal("invalid file format")
	}

	if err != nil {
		log.Printf("failed to output a result file: %s", err)
		return
	}
	if err != nil {
		log.Printf("failed to open the output file: %s", err)
		return
	}
	tag, err := id3v2.Open(output.AbsPath(), id3v2.Options{Parse: true})
	if err != nil {
		log.Printf("error while opening the output file: %s", err)
	}
	defer tag.Close()

	// Set tags
	tag.SetTitle(output.FileBaseName)
	tag.SetArtist(prog.Pfm)
	tag.SetAlbum(prog.Title)
	tag.SetYear(prog.Ft[:4])
	tag.AddCommentFrame(id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "jpn",
		Description: prog.Info,
	})

	// write tag to the aac
	if err = tag.Save(); err != nil {
		log.Printf("error while saving a tag: %s", err)
	}

	// finish downloading the file
	log.Printf("+file saved: %s", output.AbsPath())
}

func Download(
	ctx context.Context,
	wg *sync.WaitGroup,
	client *radiko.Client,
	prog radiko.Prog,
	stationID string,
) error {
	title := prog.Title
	start := prog.Ft

	startTime, err := time.ParseInLocation(DatetimeLayout, start, location)
	if err != nil {
		return fmt.Errorf("invalid start time format '%s': %s", start, err)
	}

	if startTime.After(currentTime) { // if it is in the future, skip
		log.Printf("the program is in the future [%s]%s (%s)", stationID, title, start)
		return nil
	}

	output, err := radigo.NewOutputConfig(
		fmt.Sprintf(
			"%s_%s_%s",
			startTime.In(location).Format(OutputDatetimeLayout),
			stationID,
			title,
		),
		fileFormat,
	)
	if err != nil {
		return fmt.Errorf("failed to configure output: %s", err)
	}

	if err := output.SetupDir(); err != nil {
		return fmt.Errorf("failed to setup the output dir: %s", err)
	}

	if output.IsExist() {
		log.Printf("skip [%s]%s at %s", stationID, title, start)
		log.Printf("the output file already exists: %s", output.AbsPath())
		return nil
	}

	// detach the download job
	wg.Add(1)
	go func() {
		// try fetching the recording m3u8 uri
		var uri string
		err = retry.Do(
			func() error {
				uri, err = client.TimeshiftPlaylistM3U8(ctx, stationID, startTime)
				return err
			},
			retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
				retry.DefaultDelay = 60 * time.Second
				delay := retry.BackOffDelay(n, err, config)
				log.Printf(
					"failed to get playlist.m3u8 for [%s]%s (%s): %s (retrying in %s)",
					stationID,
					title,
					start,
					err,
					delay,
				)
				// apply a default exponential back off strategy
				return delay
			}),
			retry.Attempts(MaxRetryAttempts),          // maximum retry = 8 (~6hrs)
			retry.Delay(RetryDelaySecond*time.Second), // initial delay for BackOffDelay
		)
		if len(uri) == 0 {
			wg.Done()
			return
		}
		log.Printf("start downloading [%s]%s (%s)", stationID, title, start)
		go DownloadProgram(ctx, wg, prog, uri, output)
	}()
	return nil
}
