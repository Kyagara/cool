package main

import (
	"compress/gzip"
	"log"
	"net/http"
	"sort"

	"github.com/go-json-experiment/json"
)

var (
	cachedTotalPosts                          = 0
	cachedUsersList    []UserListUser         = make([]UserListUser, 0, 1000)
	cachedUsersProfile map[string]UserProfile = make(map[string]UserProfile, 1000)
)

// Home page json
type UserList struct {
	Users      []UserListUser `json:"users"`
	TotalPosts int            `json:"totalPosts"`
}

type UserListUser struct {
	ID          uint        `json:"-"`
	Username    string      `json:"u"`
	DisplayName string      `json:"d"`
	Avatar      string      `json:"a"`
	Banner      string      `json:"b"`
	Posts       int         `json:"p"`
	Links       []UserLinks `json:"l" gorm:"-"`
}

// Home page and user page json
type UserLinks struct {
	UserID  uint   `json:"-"`
	Website string `json:"w"`
	URL     string `json:"u"`
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	username := r.URL.Query().Get("username")
	// provider := r.URL.Query().Get("provider")
	if username != "" {
		userProfile(w, r, username)
		return
	}

	if len(cachedUsersList) != 0 {
		log.Printf("Request: %s %s (cached)\n", r.Method, r.URL.Path)

		res := UserList{
			Users:      cachedUsersList,
			TotalPosts: cachedTotalPosts,
		}

		gw := gzip.NewWriter(w)
		defer gw.Close()

		err := json.MarshalWrite(gw, &res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	log.Printf("Request: %s %s (uncached)\n", r.Method, r.URL.Path)

	var usersList []UserListUser
	err := db.Table("users").
		Select(`users.id, users.username, users.display_name, users.avatar_filename as avatar, users.banner_filename as banner,
			COUNT(DISTINCT user_posts.id) as posts`).
		Joins("LEFT JOIN user_posts ON users.id = user_posts.user_id").
		Group("users.id").
		Scan(&usersList).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userLinks []UserLinks

	err = db.Table("user_links").
		Select("website,url,user_id").
		Scan(&userLinks).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userLinksMap := make(map[uint][]UserLinks)
	for _, link := range userLinks {
		userLinksMap[link.UserID] = append(userLinksMap[link.UserID], link)
	}

	totalPosts := 0
	for i := range usersList {
		totalPosts += usersList[i].Posts
		usersList[i].Links = userLinksMap[usersList[i].ID]
	}

	sort.Slice(usersList, func(i, j int) bool {
		return usersList[i].Posts > usersList[j].Posts
	})

	res := UserList{
		Users:      usersList,
		TotalPosts: totalPosts,
	}

	gw := gzip.NewWriter(w)
	defer gw.Close()

	err = json.MarshalWrite(gw, &res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cachedUsersList = usersList
	cachedTotalPosts = totalPosts
}

// User profile json
type UserProfile struct {
	ID          uint                       `json:"-"`
	Username    string                     `json:"username"`
	DisplayName string                     `json:"displayName"`
	Bio         string                     `json:"bio"`
	Avatar      string                     `json:"avatar"`
	Banner      string                     `json:"banner"`
	TotalImages int                        `json:"totalImages"`
	TotalVideos int                        `json:"totalVideos"`
	TotalPosts  int                        `json:"totalPosts"`
	Links       []UserLinks                `json:"links" gorm:"-"`
	Posts       map[string]UserProfilePost `json:"posts" gorm:"-"`
}

type UserProfilePost struct {
	ID            uint                   `json:"-"`
	Slug          string                 `json:"s"`
	Content       string                 `json:"c"`
	PostCreatedAt string                 `json:"ca"`
	Likes         int                    `json:"l"`
	IsPreview     bool                   `json:"p"`
	Media         []UserProfilePostMedia `json:"m" gorm:"-"`
}

type UserProfilePostMedia struct {
	PostID   uint   `json:"-"`
	Filename string `json:"f"`
	Type     int    `json:"t"`
}

func userProfile(w http.ResponseWriter, r *http.Request, username string) {
	setHeaders(w)

	v := cachedUsersProfile[username]
	if v.Username != "" {
		log.Printf("Request: %s %s?%s (cached)\n", r.Method, r.URL.Path, r.URL.RawQuery)

		profile := cachedUsersProfile[username]

		gw := gzip.NewWriter(w)
		defer gw.Close()

		err := json.MarshalWrite(gw, &profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	log.Printf("Request: %s %s?%s (uncached)\n", r.Method, r.URL.Path, r.URL.RawQuery)

	var user UserProfile
	err := db.Table("users").
		Select(`users.id, users.username, users.display_name, users.bio,
			users.avatar_filename as avatar, users.banner_filename as banner`).
		Where("users.username = ?", username).
		Group("users.id").
		Scan(&user).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.Table("user_links").
		Select("user_links.website, user_links.url, user_links.user_id").
		Where("user_id = ?", user.ID).
		Scan(&user.Links).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var posts []UserProfilePost

	err = db.Table("user_posts").
		Select("user_posts.id, user_posts.slug, user_posts.content, user_posts.post_created_at, user_posts.likes, user_posts.is_preview").
		Where("user_posts.user_id = ?", user.ID).
		Scan(&posts).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var postIDs []uint
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}

	var media []UserProfilePostMedia
	err = db.Table("user_post_media").
		Select("user_post_media.user_post_id AS post_id, user_post_media.filename, user_post_media.type").
		Where("user_post_media.user_post_id IN (?)", postIDs).
		Scan(&media).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mediaMap := make(map[uint][]UserProfilePostMedia)
	for _, m := range media {
		mediaMap[m.PostID] = append(mediaMap[m.PostID], UserProfilePostMedia{
			Filename: m.Filename,
			Type:     m.Type,
		})

		if m.Type == 0 {
			user.TotalImages++
		} else {
			user.TotalVideos++
		}
	}

	user.Posts = make(map[string]UserProfilePost, len(posts))
	for i, post := range posts {
		posts[i].Media = mediaMap[post.ID]
		user.Posts[post.Slug] = posts[i]
	}

	profile := UserProfile{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		Avatar:      user.Avatar,
		Banner:      user.Banner,
		TotalPosts:  len(user.Posts),
		TotalImages: user.TotalImages,
		TotalVideos: user.TotalVideos,
		Links:       user.Links,
		Posts:       user.Posts,
	}

	gw := gzip.NewWriter(w)
	defer gw.Close()

	err = json.MarshalWrite(gw, &profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cachedUsersProfile[username] = profile
}
