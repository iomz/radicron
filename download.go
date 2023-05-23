package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

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

	err = os.Rename(concatedFile, output.AbsPath())
	if err != nil {
		log.Printf(
			"Failed to output a result file: %s", err)
		return
	}
	log.Printf("+file saved: %s", output.AbsPath())
}
