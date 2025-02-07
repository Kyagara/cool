package main

import (
	"cool/internal/utils"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	videos = make(map[string]string, 20000)
	images = make(map[string]bool, 20000)
)

// Download and convert to webp
func handleImage(client http.Client, url string, filename string, postFolder string) (string, error) {
	webpOutputFilename := strings.Replace(filename, filepath.Ext(filename), ".webp", 1)

	if images[webpOutputFilename] {
		log.Printf("Found %s\n", webpOutputFilename)
		return webpOutputFilename, nil
	}

	// Test if a webp file exists
	webpPath := fmt.Sprintf("%s\\%s", postFolder, webpOutputFilename)
	size := utils.GetFileSize(webpPath)
	if size == 0 {
		originalPath := fmt.Sprintf("%s\\%s", postFolder, filename)
		err := os.MkdirAll(postFolder, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error creating directory %s", postFolder)
		}

		log.Printf("Downloading %s\n", webpOutputFilename)
		err = downloadImage(client, url, originalPath)
		if err != nil {
			return "", fmt.Errorf("error downloading %s %s", filename, err)
		}

		err = utils.ToWebP(originalPath, webpPath)
		if err != nil {
			return "", fmt.Errorf("error converting %s %s", filename, err)
		}

		images[webpOutputFilename] = true
		return webpOutputFilename, nil
	}

	images[webpOutputFilename] = true
	return webpOutputFilename, nil
}

// Use yt-dlp to download the video and make sure the container is mp4
func handleVideo(client http.Client, url string, postFolder string) (string, error) {
	// Early video check, check in the map if the video is already downloaded (avoids a request)
	v := videos[url]
	if v != "" {
		log.Printf("Found %s\n", v)
		return v, nil
	}

	// Get video.m3u8
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating video request")
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending video request")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("error status code when retrieving video")
	}

	oldFilename := res.Header.Get("stream-media-id")
	if oldFilename == "" {
		return "", fmt.Errorf("no filename found in header")
	}

	outputFilename := fmt.Sprintf("%s.mp4", oldFilename)

	// Test if the video file exists
	videoPath := fmt.Sprintf("%s\\%s", postFolder, outputFilename)
	size := utils.GetFileSize(videoPath)
	if size == 0 {
		err := os.MkdirAll(postFolder, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error creating directory %s", postFolder)
		}

		log.Printf("Downloading %s\n", outputFilename)
		err = downloadVideo(url, postFolder)
		if err != nil {
			return "", fmt.Errorf("error downloading %s %s", oldFilename, err)
		}

		// Test if a preview gif file exists
		previewFilename := fmt.Sprintf("%s.gif", oldFilename)
		previewPath := fmt.Sprintf("%s\\%s", postFolder, previewFilename)
		size := utils.GetFileSize(previewPath)
		if size == 0 {
			log.Printf("Creating preview GIF for %s\n", oldFilename)
			err := utils.MakePreviewGif(videoPath, previewPath)
			if err != nil {
				log.Printf("error creating preview GIF %s: %s", oldFilename, err)
			}
		}

		videos[url] = outputFilename
		return outputFilename, nil
	}

	videos[url] = outputFilename
	return outputFilename, nil
}

func loadImages(dataPath string) error {
	err := filepath.WalkDir(dataPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if strings.HasSuffix(name, ".webp") {
			images[name] = true
			return nil
		}

		return nil
	})

	return err
}
