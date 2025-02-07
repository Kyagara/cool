package main

import (
	"cool/internal/utils"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/go-json-experiment/json"
)

var (
	PROVIDER_NAME = flag.String("provider", "umate", "The provider to scan")

	// The folder where the scanner will save the data
	SCANNER_OUTPUT_PATH = ""
)

type Scanner interface {
	OpenDB() error
	Start(page int, pageSize int) error
	GetNextPage() string
	GetTotalPages() int
	Save(client http.Client, data any) error
}

func main() {
	workers := flag.Int("workers", runtime.NumCPU(), "Number of concurrent workers")
	startPage := flag.Int("page", 1, "Start page")
	perPage := flag.Int("per", 20, "Posts per page")
	flag.Parse()

	closeLog, err := utils.SetupLogging("main.log")
	if err != nil {
		log.Printf("Error setting up logging: %v\n", err)
		return
	}
	defer closeLog()

	var scanner Scanner

	switch *PROVIDER_NAME {
	case "umate":
		scanner = &Umate{}
	default:
		log.Printf("Unknown scanner: %s\n", *PROVIDER_NAME)
		return
	}

	SCANNER_OUTPUT_PATH = fmt.Sprintf("static\\provider\\%s", *PROVIDER_NAME)
	err = os.MkdirAll(SCANNER_OUTPUT_PATH, os.ModePerm)
	if err != nil {
		log.Printf("Error creating data directory: %v\n", err)
		return
	}

	// Loading all the images from the provider, this is used to check if the image is already downloaded, skipping the download if it is
	err = loadImages(SCANNER_OUTPUT_PATH)
	if err != nil {
		return
	}

	// Opening the database, also should create the database if it doesn't exist and do db.AutoMigrate()
	err = scanner.OpenDB()
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		return
	}

	numWorkers := *workers
	urls := make(chan string)

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go worker(scanner, urls, &wg)
	}

	err = scanner.Start(*startPage, *perPage)
	if err != nil {
		log.Printf("Error starting scanner: %s\n", err)
		return
	}

	for i := range scanner.GetTotalPages() {
		if i%10 == 0 {
			log.Printf("Processing page %d\n", i)
		}

		url := scanner.GetNextPage()
		if url == "" {
			break
		}

		urls <- url
	}

	close(urls)
	wg.Wait()

	log.Println("All done.")
}

func worker(scanner Scanner, urls <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{}

	for url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Error creating request: %v\n", err)
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request: %v\n", err)
			continue
		}

		if res.StatusCode != 200 {
			log.Printf("Error status code: %d\n", res.StatusCode)
			res.Body.Close()
			continue
		}

		var data any

		switch *PROVIDER_NAME {
		case "umate":
			data = UmateModel{}
			err = json.UnmarshalRead(res.Body, &data)
			if err != nil {
				res.Body.Close()
				log.Printf("Error unmarshalling response: %v, %s\n", err, url)
				continue
			}

		default:
			log.Printf("Unknown scanner: %s\n", *PROVIDER_NAME)
			res.Body.Close()
			continue
		}

		res.Body.Close()

		err = scanner.Save(client, data)
		if err != nil {
			log.Printf("Error saving data: %v\n", err)
			continue
		}
	}
}
