package eseis

import (
	"encoding/json"
	"fmt"
	"github.com/idkw/eseisscrapper/pkg/infrastructure/utils"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type forumTopicsResponse struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DisplayName string    `json:"display_name"`
	State       string    `json:"state"`
	Read        bool      `json:"read"`
	PostCount   int       `json:"post_count"`
	Description string    `json:"description"`
	Author      struct {
		ID               int         `json:"id"`
		TopPlaceRole     string      `json:"top_place_role"`
		InvitedBy        interface{} `json:"invited_by"`
		RealName         string      `json:"real_name"`
		IdentityID       int         `json:"identity_id"`
		DisplayPlaceRole string      `json:"display_place_role"`
		DisplayName      string      `json:"display_name"`
		AvatarURL        string      `json:"avatar_url"`
		SergicPartner    bool        `json:"sergic_partner"`
		CanTalk          bool        `json:"can_talk"`
		Contracts        []struct {
			ID                      int    `json:"id"`
			DisplayName             string `json:"display_name"`
			PendingAmount           int    `json:"pending_amount"`
			CustomerReferenceNumber string `json:"customer_reference_number"`
		} `json:"contracts"`
	} `json:"author"`
	EditedAt            interface{} `json:"edited_at"`
	DisplayCategoryKind string      `json:"display_category_kind"`
	PlaceDisplayName    string      `json:"place_display_name"`
	Category            struct {
		ID              int       `json:"id"`
		Kind            string    `json:"kind"`
		DisplayName     string    `json:"display_name"`
		DisplayIcon     string    `json:"display_icon"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
		DisplayIconURL  string    `json:"display_icon_url"`
		OpenTopicsCount int       `json:"open_topics_count"`
		TopicsCount     int       `json:"topics_count"`
	} `json:"category"`
	Attachments []struct {
		ID                int         `json:"id"`
		UUID              string      `json:"uuid"`
		SourceFileName    string      `json:"source_file_name"`
		SourceContentType string      `json:"source_content_type"`
		SourceFileSize    int         `json:"source_file_size"`
		FileURL           string      `json:"file_url"`
		Dimensions        interface{} `json:"dimensions"`
		SourceUpdatedAt   time.Time   `json:"source_updated_at"`
		DisplayName       interface{} `json:"display_name"`
	} `json:"attachments"`
}

func (f ForumTopic) CleanDisplayName() string {
	return strings.Trim(strings.ReplaceAll(f.DisplayName, "/", "_"), "")
}

type ForumTopic struct {
	ID          int
	UUID        string
	DisplayName string
	CreatedAt   time.Time `json:"created_at"`
	Raw         *forumTopicsResponse
}

func (e *EseisClient) GetForumTopics(placeID int, page int) ([]ForumTopic, error) {
	path := fmt.Sprintf("/v1/places/%d/forum/topics?page=%d&per_page=20&sort=-updated_at", placeID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create forum topics request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send forum topics request: %w", err)
	}
	defer resp.Body.Close()

	var response []forumTopicsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode forum topics response: %w", err)
	}

	var forumTopics = make([]ForumTopic, len(response))
	for i, r := range response {
		forumTopics[i] = ForumTopic{
			ID:          r.ID,
			UUID:        r.UUID,
			DisplayName: r.DisplayName,
			CreatedAt:   r.CreatedAt,
			Raw:         &r,
		}
	}
	return forumTopics, nil
}

func (e *EseisClient) CreateForumTopicScreenshot(forumTopic ForumTopic, outDir string) error {
	forumTopicName := forumTopic.CleanDisplayName()
	year, month, day := forumTopic.CreatedAt.Date()
	forumTopicFileName := fmt.Sprintf("%d_%d_%d__%d__%s.pdf", year, month, day, forumTopic.ID, forumTopicName)
	forumTopicPath := filepath.Join(outDir, forumTopicFileName)
	url := fmt.Sprintf("https://client.eseis-syndic.com/mes-echanges/forum/%d", forumTopic.ID)
	return e.SavePDF(url, forumTopicPath, WaitForForumPageActions()...)
}

type postResponse struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    struct {
		ID               int           `json:"id"`
		TopPlaceRole     string        `json:"top_place_role"`
		InvitedBy        interface{}   `json:"invited_by"`
		RealName         interface{}   `json:"real_name"`
		IdentityID       int           `json:"identity_id"`
		DisplayPlaceRole string        `json:"display_place_role"`
		DisplayName      string        `json:"display_name"`
		AvatarURL        string        `json:"avatar_url"`
		SergicPartner    bool          `json:"sergic_partner"`
		CanTalk          bool          `json:"can_talk"`
		Contracts        []interface{} `json:"contracts"`
	} `json:"author"`
	Attachments []struct {
		ID                int         `json:"id"`
		UUID              string      `json:"uuid"`
		SourceFileName    string      `json:"source_file_name"`
		SourceContentType string      `json:"source_content_type"`
		SourceFileSize    int         `json:"source_file_size"`
		FileURL           string      `json:"file_url"`
		Dimensions        interface{} `json:"dimensions"`
		SourceUpdatedAt   time.Time   `json:"source_updated_at"`
		DisplayName       interface{} `json:"display_name"`
	} `json:"attachments"`
}

func (e *EseisClient) GetAllTopicPosts(placeID int, topicID int) ([]postResponse, error) {
	allPosts := make([]postResponse, 0)
	page := 1
	for {
		posts, err := e.GetTopicPosts(placeID, topicID, page)
		utils.MustBeNilErr(err, "failed to get topic posts for placeID=%d topic=%d page=%d", placeID, topicID, page)
		if len(posts) == 0 {
			break
		}
		allPosts = append(allPosts, posts...)

		page++
	}
	return allPosts, nil
}

func (e *EseisClient) GetTopicPosts(placeID int, topicID int, page int) ([]postResponse, error) {
	path := fmt.Sprintf("/v1/forum/topics/%d/posts?page=%d&per_page=15&sort=-updated_at", topicID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	req.Header.Set("x-current-place-id", fmt.Sprintf("%d", placeID))
	if err != nil {
		return nil, fmt.Errorf("failed to create topic posts request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send topic posts request: %w", err)
	}
	defer resp.Body.Close()

	var response []postResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode topic posts response: %w", err)
	}
	return response, nil
}
