package main

import (
	"compress/gzip"
	"log"
	"net/http"

	"github.com/go-json-experiment/json"
)

var (
	cachedGallery []Gallery = nil
)

// Gallery page json
type Gallery struct {
	Username      string `json:"u"`
	Slug          string `json:"s"`
	Filename      string `json:"f"`
	Likes         int    `json:"l"`
	PostCreatedAt string `json:"ca"`
	Type          int    `json:"t"`
}

func handleGallery(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	if cachedGallery != nil {
		log.Printf("Request: %s %s (cached)\n", r.Method, r.URL.Path)

		gw := gzip.NewWriter(w)
		defer gw.Close()

		err := json.MarshalWrite(gw, &cachedGallery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	log.Printf("Request: %s %s (uncached)\n", r.Method, r.URL.Path)

	var gallery []Gallery

	err := db.Table("user_post_media").
		Select(`users.username, user_posts.slug, user_post_media.filename,
				user_post_media.type, user_posts.likes, user_posts.post_created_at`).
		Joins("JOIN user_posts ON user_post_media.user_post_id = user_posts.id").
		Joins("JOIN users ON user_posts.user_id = users.id").
		Scan(&gallery).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gw := gzip.NewWriter(w)
	defer gw.Close()

	err = json.MarshalWrite(gw, &gallery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cachedGallery = gallery
}
