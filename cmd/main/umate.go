package main

import (
	"cool/internal/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Umate struct {
	CurrentPage int
	TotalPages  int
	pageSize    int
	users       map[string]*sync.Mutex
	db          *gorm.DB
	mu          sync.Mutex
}

func (u *Umate) OpenDB() error {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 1 * time.Second,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}

	dsn := os.Getenv("PRIVATE_DATABASE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.UserLink{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.UserPost{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.UserPostMedia{})
	if err != nil {
		return err
	}

	u.db = db
	return nil
}

func (u *Umate) Start(page int, pageSize int) error {
	client := http.Client{}
	res, err := client.Get(fmt.Sprintf("https://api.umate.me/api/post/home_post_list?page=%d&page_size=%d", page, pageSize))
	if err != nil {
		return fmt.Errorf("failed to get first page")
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to get first page")
	}

	var data UmateModel
	err = json.UnmarshalRead(res.Body, &data)
	if err != nil {
		return fmt.Errorf("failed to decode first page")
	}

	u.CurrentPage++
	u.TotalPages = data.Data.LastPage
	u.pageSize = pageSize
	u.users = make(map[string]*sync.Mutex, 1000)
	return nil
}

func (u *Umate) GetNextPage() string {
	urlF := "https://api.umate.me/api/post/home_post_list?page=%d&page_size=%d"

	u.mu.Lock()
	defer u.mu.Unlock()

	u.CurrentPage++
	if u.CurrentPage > u.TotalPages {
		return ""
	}

	url := fmt.Sprintf(urlF, u.CurrentPage, u.pageSize)
	return url
}

func (u *Umate) Save(client http.Client, data any) error {
	model := data.(UmateModel).Data.Data
	for _, post := range model {
		username := post.User.Name

		invalidChars := `<>:"/\|?[]{}*`
		if strings.ContainsAny(username, invalidChars) {
			log.Printf("Invalid username: %s\n", username)
			continue
		}

		slug := post.Slug

		userPath := fmt.Sprintf("%s\\%s", SCANNER_OUTPUT_PATH, username)
		postPath := fmt.Sprintf("%s\\%s", userPath, slug)

		err := os.MkdirAll(postPath, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory: %v\n", err)
			continue
		}

		mu := u.users[username]
		if mu == nil {
			mu = &sync.Mutex{}
			u.users[username] = mu
		}

		mu.Lock()

		var newUser models.User
		err = u.db.Where("username = ?", username).Find(&newUser).Error
		if err != nil {
			log.Printf("Error getting user: %v\n", err)
			mu.Unlock()
			continue
		}

		if newUser.Username == "" {
			newUser.Username = username
			newUser.DisplayName = post.User.Nickname
			newUser.Bio = post.User.Description

			avatarURL := u.getMediaPath(post.User.Avatar)
			avatar, err := handleImage(client, avatarURL, u.getImageFilename(avatarURL), userPath)
			if err != nil {
				log.Printf("Error checking %s avatar %s\n", username, err)
			}

			newUser.Avatar.Filename = avatar
			newUser.Avatar.URL = avatarURL

			bannerURL := u.getMediaPath(post.User.Banner)
			banner, err := handleImage(client, bannerURL, u.getImageFilename(bannerURL), userPath)
			if err != nil {
				log.Printf("Error checking %s banner %s\n", username, err)
			}

			newUser.Banner.Filename = banner
			newUser.Banner.URL = bannerURL

			err = u.db.Create(&newUser).Error
			if err != nil {
				log.Printf("Error creating user: %v\n", err)
				mu.Unlock()
				continue
			}

			log.Printf("Created user %s\n", username)

			social := make([]models.UserLink, 0, len(post.User.Social))

			if post.User.PersonalWebsite != "" {
				social = append(social, models.UserLink{
					UserID:           newUser.ID,
					Username:         username,
					Website:          "Website",
					URL:              post.User.PersonalWebsite,
					UniqueConstraint: fmt.Sprintf("%s_website", username),
				})

			}

			for _, socialMedia := range post.User.Social {
				social = append(social, models.UserLink{
					UserID:           newUser.ID,
					Username:         socialMedia.SocialUsername,
					Website:          socialMedia.SocialType,
					URL:              socialMedia.URL,
					UniqueConstraint: fmt.Sprintf("%s_%s", username, socialMedia.SocialType),
				})
			}

			if len(social) != 0 {
				err := u.db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&social).Error
				if err != nil {
					log.Printf("Error updating %s social: %v\n", username, err)
				}
			}
		}

		mu.Unlock()

		var newPost models.UserPost
		err = u.db.Where("slug = ?", slug).Find(&newPost).Error
		if err != nil {
			log.Printf("Error getting post: %v\n", err)
			continue
		}

		if newPost.Slug != "" {
			continue
		}

		layout := "2006-01-02 15:04:05"
		postCreation, _ := time.Parse(layout, post.CreatedAt)
		postUpdate, _ := time.Parse(layout, post.UpdatedAt)
		likes, _ := strconv.Atoi(post.LikeNum)
		var isPreview bool
		if post.IsPreview == 1 {
			isPreview = true
		}

		newPost = models.UserPost{
			Slug:          slug,
			UserID:        newUser.ID,
			Content:       post.Content,
			PostCreatedAt: postCreation,
			PostUpdatedAt: postUpdate,
			Likes:         likes,
			IsPreview:     isPreview,
		}

		err = u.db.Create(&newPost).Error
		if err != nil {
			log.Printf("Error creating %s %s post: %v\n", username, slug, err)
			continue
		}

		log.Printf("Created %s post %s\n", username, slug)

		newPostMedia := make([]models.UserPostMedia, 0, len(post.Media))
		for _, media := range post.Media {
			mediaURL := u.getMediaPath(media.Paths)
			if mediaURL == "" {
				continue
			}

			var name string

			if media.ResourceType == 0 {
				name, err = handleImage(client, mediaURL, u.getImageFilename(mediaURL), postPath)
				if err != nil {
					log.Printf("Error checking %s %s image %s: %s\n", username, slug, mediaURL, err)
					continue
				}
			} else {
				name, err = handleVideo(client, mediaURL, postPath)
				if err != nil {
					log.Printf("Error checking %s %s video %s: %s\n", username, slug, mediaURL, err)
					continue
				}
			}

			if name == "" {
				continue
			}

			newMedia := models.UserPostMedia{
				UserPostID: newPost.ID,
				UserMedia: models.UserMedia{
					Filename: name,
					URL:      mediaURL,
				},
				Type:   media.ResourceType,
				Width:  media.Width,
				Height: media.Height,
			}

			if media.ResourceType != 0 {
				newMedia.Duration = &media.Duration
			}

			newPostMedia = append(newPostMedia, newMedia)
		}

		if len(newPostMedia) == 0 {
			continue
		}

		err = u.db.Create(&newPostMedia).Error
		if err != nil {
			log.Printf("Error creating %s %s post: %v\n", username, slug, err)
			continue
		}
	}

	return nil
}

func (u *Umate) GetTotalPages() int {
	return u.TotalPages
}

func (u *Umate) getMediaPath(media any) string {
	sizes := []string{"large", "medium", "small"}

	switch object := media.(type) {
	case string:
		return object

	case map[string]string:
		for _, size := range sizes {
			objectValue := object[size]
			if objectValue != "" {
				return objectValue
			}
		}

	case []map[string]string:
		for _, objectMap := range object {
			for _, size := range sizes {
				objectValue := objectMap[size]
				if objectValue != "" {
					return objectValue
				}
			}
		}

	case map[string]any:
		for _, size := range sizes {
			objectValue := object[size]
			if objectValue != nil {
				return objectValue.(string)
			}
		}

	case []any:
		return ""

	default:
		log.Printf("Unknown type %T\n", object)
	}

	return ""
}

func (u *Umate) getImageFilename(url string) string {
	filenameRegex := regexp.MustCompile(`([^\/\?]+)(?:\?|$)`)
	matches := filenameRegex.FindStringSubmatch(url)
	if len(matches) == 0 {
		return ""
	}

	return matches[1]
}

type UmateModel struct {
	Status     string         `json:"status"`
	Data       UmateModelData `json:"data"`
	StatusCode int            `json:"status_code"`
}

type UmateModelData struct {
	Data       []UmateModelPost `json:"data"`
	LastPage   int              `json:"last_page"`
	TotalPosts int              `json:"total"`
}

type UmateModelPostDetail struct {
	Status     string         `json:"status"`
	Data       UmateModelPost `json:"data"`
	StatusCode int            `json:"status_code"`
}

type UmateModelPost struct {
	Slug               string            `json:"post_slug"`
	Content            string            `json:"content"`
	Price              string            `json:"price"`
	CreatedAt          string            `json:"created_at"`
	UpdatedAt          string            `json:"updated_at"`
	SingleUnlockPrice  string            `json:"single_unlock_price"`
	LikeNum            string            `json:"like_num"`
	CommentNum         string            `json:"comment_num"`
	User               UmateModelUser    `json:"author,omitzero"`
	Media              []UmateModelMedia `json:"media"`
	MediaType          int               `json:"media_type"`
	IsPreview          int               `json:"is_preview"`
	VideoConvertStatus int               `json:"video_convert_status"`
	Readable           int               `json:"readable"`
	IsBuy              int               `json:"is_buy"`
}

type UmateModelUser struct {
	Name            string             `json:"name"`
	Nickname        string             `json:"nickname"`
	Description     string             `json:"description"`
	PersonalWebsite string             `json:"personal_website"`
	Avatar          map[string]string  `json:"avatar"`
	Banner          map[string]string  `json:"banner"`
	Social          []UmateModelSocial `json:"social_media"`
}

type UmateModelMedia struct {
	Paths        any `json:"paths"`
	ResourceType int `json:"resource_type"`
	IsCf         int `json:"is_cf"`
	Scene        int `json:"scene"`
	Width        int `json:"width"`
	Height       int `json:"height"`
	Duration     int `json:"duration"`
}

type UmateModelSocial struct {
	SocialUsername string `json:"social_username"`
	SocialType     string `json:"social_type"`
	URL            string `json:"url"`
}
