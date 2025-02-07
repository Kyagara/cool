package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func downloadImage(client http.Client, url string, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadVideo(url string, dest string) error {
	format := "%(title)s.%(ext)s"

	cmd := fmt.Sprintf("%s -f \"bv+ba/bestvideo+bestaudio\" --merge-output-format mp4", url)
	command := exec.Command("yt-dlp", cmd, "-o", format)
	command.Dir = dest

	err := command.Run()
	if err != nil {
		return err
	}

	return nil
}
