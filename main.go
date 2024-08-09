package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func compressVideo(inputPath, outputPath string, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	err := ffmpeg.Input(inputPath).
		Output(outputPath, ffmpeg.KwArgs{"vcodec": "libx264", "crf": "28"}).
		OverWriteOutput().
		Run()
	if err != nil {
		errChan <- err
	} else {
		log.Printf("Video %s compressed successfully", inputPath)
	}
}

func main() {
	inputDir := "./input"
	outputDir := "./output"

	var wg sync.WaitGroup
	errChan := make(chan error, 10)

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			outputPath := filepath.Join(outputDir, info.Name())
			wg.Add(1)
			go compressVideo(path, outputPath, &wg, errChan)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through input directory: %v", err)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		log.Printf("Error compressing video: %v", err)
	}
}
