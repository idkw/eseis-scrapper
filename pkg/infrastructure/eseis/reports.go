package eseis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type reportResponse struct {
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

type Report struct {
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

type Author struct {
	DisplayPlaceRole string
	DisplayName      string
}

func (e *EseisClient) GetReports(placeID int, page int) ([]Report, error) {
	path := fmt.Sprintf("/v1/places/%d/reports?page=%d&per_page=10&sort=created_at", placeID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create contracts request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send contracts request: %w", err)
	}
	defer resp.Body.Close()

	var reportsResponse []reportResponse
	err = json.NewDecoder(resp.Body).Decode(&reportsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contracts response: %w", err)
	}

	var reports = make([]Report, len(reportsResponse))
	for i, report := range reportsResponse {
		reports[i] = Report{
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

func (e *EseisClient) CreateReportScreenshot(report Report, outDir string) error {
	reportName := strings.Trim(strings.ReplaceAll(report.DisplayName, "/", "_"), "")
	year, month, day := report.CreatedAt.Date()
	reportFileName := fmt.Sprintf("%d_%d_%d__%d__%s.pdf", year, month, day, report.ID, reportName)
	reportPath := filepath.Join(outDir, reportFileName)
	return e.SavePDF(report.URL, reportPath, WaitForReportPageActions()...)
}
