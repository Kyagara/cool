package main

import (
	"context"
	"cool/internal/utils"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func main() {
	h2 := flag.Bool("h2", false, "Use HTTP/2 instead of HTTP/1.1")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("PRIVATE_DATABASE")
	host := os.Getenv("PUBLIC_HOST")
	port := os.Getenv("PUBLIC_PORT")
	domain := fmt.Sprintf("%s:%s", host, port)

	closeLog, err := utils.SetupLogging("backend.log")
	if err != nil {
		log.Fatalf("Error setting up logging: %v\n", err)
	}

	defer closeLog()

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/gallery", handleGallery)

	// Built website
	mux.Handle("/", handleWebsite())
	// Provider content
	mux.Handle("/provider/", http.StripPrefix("/provider/", handleCDN()))

	server := http.Server{
		Addr:    domain,
		Handler: mux,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		log.Println("Shutting down")
		err = server.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("Error shutting down server: %v\n", err)
		}
	}()

	var url string
	if *h2 {
		url = fmt.Sprintf("https://%s", domain)
	} else {
		url = fmt.Sprintf("http://%s", domain)
	}

	log.Printf("Running on %s\n", url)

	if *h2 {
		err = server.ListenAndServeTLS("cert.pem", "key.pem")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}

		return
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

func handleWebsite() http.Handler {
	buildDir := "static\\dist"
	fs := http.FileServer(http.Dir(buildDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")

		path := filepath.Join(buildDir, r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(buildDir, "index.html"))
			return
		}

		fs.ServeHTTP(w, r)
	})
}

func handleCDN() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")

		if strings.HasSuffix(r.URL.Path, "/") {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}

		http.FileServer(http.Dir("static\\provider")).ServeHTTP(w, r)
	})
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Encoding", "gzip")
}
