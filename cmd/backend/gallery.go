package main

import (
	"compress/gzip"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/go-json-experiment/json"
)


var (
	ttlTime = 1
	firstPageCache     *GalleryResponse
	firstPageCacheTime time.Time     
	firstPageCacheTTL  = time.Hour * time.Duration(ttlTime)
)

type Page struct {
	PageNumber    int `json:"page"`
	PostLimit     int `json:"postlimit"`
	NumberOfPages int `json:"total"`
}

type GalleryResponse struct {
	Data []Gallery `json:"data"`
	Page Page      `json:"page"`
}

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

	// Default pagination values.
	page := 1
	postLimit := 50

	var gallery []Gallery

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	offset := (page - 1) * postLimit

	if page == 1 
	{
		if firstPageCache != nil && time.Since(firstPageCacheTime) < firstPageCacheTTL
		{
			gw := gzip.NewWriter(w)
			defer gw.Close()
		
			err = json.MarshalWrite(gw, firstPageCache)
			if err != nil 
			{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	err := db.Table("user_post_media").
		Select(`users.username, user_posts.slug, user_post_media.filename,
			user_post_media.type, user_posts.likes, user_posts.post_created_at`).
		Joins("JOIN user_posts ON user_post_media.user_post_id = user_posts.id").
		Joins("JOIN users ON user_posts.user_id = users.id").
		Order("user_posts.post_created_at DESC").
		Limit(postLimit). 
		Offset(offset).  
		Scan(&gallery).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the total number of records for pagination.
	var totalCount int
	err = db.Table("user_post_media").Count(&totalCount).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(postLimit)))

	response := GalleryResponse{
		Data: gallery,
		Page: Page{
			PageNumber:    page,
			PostLimit:     postLimit,
			NumberOfPages: totalPages,
		},
	}

	if page == 1 {
		firstPageCache = &response
		firstPageCacheTime = time.Now()
	}


	gw := gzip.NewWriter(w)
	defer gw.Close()

	err = json.MarshalWrite(gw, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
