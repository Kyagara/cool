package utils

import (
	"io"
	"log"
	"os"
)

// https://gist.github.com/jerblack/4b98ba48ed3fb1d9f7544d2b1a1be287
func SetupLogging(stdoutPath string) (func(), error) {
	file, err := os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	writer := io.MultiWriter(os.Stdout, file)
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	os.Stdout = writePipe
	os.Stderr = writePipe

	log.SetOutput(writer)

	exit := make(chan bool)

	go func() {
		_, err = io.Copy(writer, readPipe)
		if err != nil {
			log.Printf("Error copying from pipe: %v\n", err)
		}

		exit <- true
	}()

	return func() {
		err = writePipe.Close()
		if err != nil {
			log.Printf("Error closing write pipe: %v\n", err)
		}

		<-exit

		err = file.Close()
		if err != nil {
			log.Printf("Error closing file: %v\n", err)
		}
	}, nil
}
