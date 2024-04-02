package eseis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type reportSummaryResponse struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	State       string    `json:"state"`
	Author      struct {
		DisplayPlaceRole string `json:"display_place_role"`
		DisplayName      string `json:"display_name"`
	} `json:"author,omitempty"`
}

type ReportSummary struct {
	ID          int
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	State       string
	Author      struct {
		DisplayPlaceRole string
		DisplayName      string
	}
	URL string
}

func (r ReportSummary) CleanDisplayName() any {
	return strings.Trim(strings.ReplaceAll(r.DisplayName, "/", "_"), "")
}

type Author struct {
	DisplayPlaceRole string
	DisplayName      string
}

type reportResponse struct {
	ID               int         `json:"id"`
	DisplayName      string      `json:"display_name"`
	Description      string      `json:"description"`
	WitnessesCount   int         `json:"witnesses_count"`
	ResolutionTime   interface{} `json:"resolution_time"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	UserWitness      bool        `json:"user_witness"`
	EquipmentID      int         `json:"equipment_id"`
	EquipmentName    string      `json:"equipment_name"`
	EquipmentIcon    string      `json:"equipment_icon"`
	PlaceID          int         `json:"place_id"`
	State            string      `json:"state"`
	PlaceDisplayName string      `json:"place_display_name"`
	PlacePictureURL  interface{} `json:"place_picture_url"`
	Category         struct {
		ID           int       `json:"id"`
		DisplayName  string    `json:"display_name"`
		DisplayColor string    `json:"display_color"`
		DisplayIcon  string    `json:"display_icon"`
		Kind         string    `json:"kind"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Sort         int       `json:"sort"`
	} `json:"category"`
	CanResolve       bool        `json:"can_resolve"`
	ReportResolution interface{} `json:"report_resolution"`
	Attachments      []struct {
		ID                int         `json:"id"`
		UUID              string      `json:"uuid"`
		SourceFileName    string      `json:"source_file_name"`
		SourceContentType string      `json:"source_content_type"`
		SourceFileSize    int         `json:"source_file_size"`
		FileURL           string      `json:"file_url"`
		Dimensions        []int       `json:"dimensions"`
		SourceUpdatedAt   time.Time   `json:"source_updated_at"`
		DisplayName       interface{} `json:"display_name"`
		Width             int         `json:"width"`
		Height            int         `json:"height"`
	} `json:"attachments"`
	ReportEvents []struct {
		ID          int       `json:"id"`
		Description string    `json:"description"`
		Kind        string    `json:"kind"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DisplayName string    `json:"display_name"`
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
		Attachments []struct {
			ID                int         `json:"id"`
			UUID              string      `json:"uuid"`
			SourceFileName    string      `json:"source_file_name"`
			SourceContentType string      `json:"source_content_type"`
			SourceFileSize    int         `json:"source_file_size"`
			FileURL           string      `json:"file_url"`
			Dimensions        []int       `json:"dimensions"`
			SourceUpdatedAt   time.Time   `json:"source_updated_at"`
			DisplayName       interface{} `json:"display_name"`
			Width             int         `json:"width"`
			Height            int         `json:"height"`
		} `json:"attachments"`
	} `json:"report_events"`
}

type Report struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	State       string    `json:"state"`
	Attachments []struct {
		ID                int         `json:"id"`
		UUID              string      `json:"uuid"`
		SourceFileName    string      `json:"source_file_name"`
		SourceContentType string      `json:"source_content_type"`
		SourceFileSize    int         `json:"source_file_size"`
		FileURL           string      `json:"file_url"`
		Dimensions        []int       `json:"dimensions"`
		SourceUpdatedAt   time.Time   `json:"source_updated_at"`
		DisplayName       interface{} `json:"display_name"`
		Width             int         `json:"width"`
		Height            int         `json:"height"`
	} `json:"attachments"`
	ReportEvents []struct {
		ID          int       `json:"id"`
		Description string    `json:"description"`
		Kind        string    `json:"kind"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DisplayName string    `json:"display_name"`
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
		Attachments []struct {
			ID                int         `json:"id"`
			UUID              string      `json:"uuid"`
			SourceFileName    string      `json:"source_file_name"`
			SourceContentType string      `json:"source_content_type"`
			SourceFileSize    int         `json:"source_file_size"`
			FileURL           string      `json:"file_url"`
			Dimensions        []int       `json:"dimensions"`
			SourceUpdatedAt   time.Time   `json:"source_updated_at"`
			DisplayName       interface{} `json:"display_name"`
			Width             int         `json:"width"`
			Height            int         `json:"height"`
		} `json:"attachments"`
	} `json:"report_events"`
}

func (e *EseisClient) GetReportSummaries(placeID int, page int) ([]ReportSummary, error) {
	path := fmt.Sprintf("/v1/places/%d/reports?page=%d&per_page=10&sort=created_at", placeID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create report summaries request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send report summaries request: %w", err)
	}
	defer resp.Body.Close()

	var reportsResponse []reportSummaryResponse
	err = json.NewDecoder(resp.Body).Decode(&reportsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode report summaries response: %w", err)
	}

	var reports = make([]ReportSummary, len(reportsResponse))
	for i, report := range reportsResponse {
		reports[i] = ReportSummary{
			ID:          report.ID,
			DisplayName: report.DisplayName,
			CreatedAt:   report.CreatedAt,
			UpdatedAt:   report.UpdatedAt,
			State:       report.State,
			Author: Author{
				DisplayName:      report.Author.DisplayName,
				DisplayPlaceRole: report.Author.DisplayPlaceRole,
			},
			URL: e.buildWebURL(fmt.Sprintf("/mes-echanges/signalements/%d", report.ID)),
		}
	}
	return reports, nil
}

func (e *EseisClient) GetReport(reportID int) (Report, error) {
	path := fmt.Sprintf("/v1/reports/%d", reportID)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return Report{}, fmt.Errorf("failed to create report request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return Report{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Report{}, fmt.Errorf("failed to send report request: %w", err)
	}
	defer resp.Body.Close()

	var rawReport reportResponse
	err = json.NewDecoder(resp.Body).Decode(&rawReport)
	if err != nil {
		return Report{}, fmt.Errorf("failed to decode report response: %w", err)
	}

	return Report{
		ID:           rawReport.ID,
		DisplayName:  rawReport.DisplayName,
		Description:  rawReport.Description,
		CreatedAt:    rawReport.CreatedAt,
		UpdatedAt:    rawReport.UpdatedAt,
		State:        rawReport.State,
		Attachments:  rawReport.Attachments,
		ReportEvents: rawReport.ReportEvents,
	}, nil
}

func (e *EseisClient) CreateReportScreenshot(report ReportSummary, outDir string) error {
	reportName := strings.Trim(strings.ReplaceAll(report.DisplayName, "/", "_"), "")
	year, month, day := report.CreatedAt.Date()
	reportFileName := fmt.Sprintf("%d_%d_%d__%d__%s.pdf", year, month, day, report.ID, reportName)
	reportPath := filepath.Join(outDir, reportFileName)
	return e.SavePDF(report.URL, reportPath, WaitForReportPageActions()...)
}
