package eseis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type fiscalYearResponse struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"display_name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type FiscalYear struct {
	ID          int
	DisplayName string
	StartDate   time.Time
	EndDate     time.Time
}

func (e *EseisClient) GetFiscalYears(placeID int) ([]FiscalYear, error) {
	path := fmt.Sprintf("/v1/places/%d/fiscal_years", placeID)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fiscal_years request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send fiscal_years request: %w", err)
	}
	defer resp.Body.Close()

	var fiscalYearsResponse []fiscalYearResponse
	err = json.NewDecoder(resp.Body).Decode(&fiscalYearsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode fiscal_years response: %w", err)
	}

	var fiscalYears = make([]FiscalYear, len(fiscalYearsResponse))
	for i, fiscalYearResponse := range fiscalYearsResponse {
		fiscalYears[i] = FiscalYear{
			ID:          fiscalYearResponse.ID,
			DisplayName: fiscalYearResponse.DisplayName,
			StartDate:   fiscalYearResponse.StartDate,
			EndDate:     fiscalYearResponse.EndDate,
		}
	}
	return fiscalYears, nil
}

type budgetResponse struct {
	ID              int    `json:"id"`
	DisplayName     string `json:"display_name"`
	AllocatedAmount int    `json:"allocated_amount"`
	SpentAmount     int    `json:"spent_amount"`
}

type Budget struct {
	ID              int
	DisplayName     string
	AllocatedAmount int
	SpentAmount     int
}

func (e *EseisClient) GetBudgets(placeID int, fiscalYearID int) ([]Budget, error) {
	path := fmt.Sprintf("/v1/places/%d/budgets?fiscal_year_id=%d", placeID, fiscalYearID)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create budgets request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send budgets request: %w", err)
	}
	defer resp.Body.Close()

	var budgetsResponse []budgetResponse
	err = json.NewDecoder(resp.Body).Decode(&budgetsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode budgets response: %w", err)
	}

	var budgets = make([]Budget, len(budgetsResponse))
	for i, budgetResponse := range budgetsResponse {
		budgets[i] = Budget{
			ID:              budgetResponse.ID,
			DisplayName:     budgetResponse.DisplayName,
			AllocatedAmount: budgetResponse.AllocatedAmount,
			SpentAmount:     budgetResponse.SpentAmount,
		}
	}
	return budgets, nil
}

type accountPlaceEntriesResponse struct {
	ID            int       `json:"id"`
	UUID          string    `json:"uuid"`
	DisplayName   string    `json:"display_name"`
	Amount        int       `json:"amount"`
	OperationDate time.Time `json:"operation_date"`
	UpdatedAt     time.Time `json:"updated_at"`
	Attachment    struct {
		DisplayName string `json:"display_name"`
		FileURL     string `json:"file_url"`
	} `json:"attachment"`
	ProviderID          int    `json:"provider_id"`
	ProviderDisplayName string `json:"provider_display_name"`
	DisplayState        string `json:"display_state"`
}

type AccountPlaceEntry struct {
	ID                  int
	UUID                string
	DisplayName         string
	Amount              int
	OperationDate       time.Time
	UpdatedAt           time.Time
	Attachment          *AccountPlaceEntryAttachment
	ProviderID          int
	ProviderDisplayName string
	DisplayState        string
}

type AccountPlaceEntryAttachment struct {
	DisplayName string
	FileURL     string
}

func (e *EseisClient) GetAccountPlaceEntries(budgetID int) ([]AccountPlaceEntry, error) {
	path := fmt.Sprintf("/v1/budgets/%d/account_place_entries", budgetID)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response []accountPlaceEntriesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var elements = make([]AccountPlaceEntry, len(response))
	for i, responseElement := range response {
		elements[i] = AccountPlaceEntry{
			ID:            responseElement.ID,
			UUID:          responseElement.UUID,
			DisplayName:   responseElement.DisplayName,
			Amount:        responseElement.Amount,
			OperationDate: responseElement.OperationDate,
			UpdatedAt:     responseElement.UpdatedAt,
			Attachment: &AccountPlaceEntryAttachment{
				DisplayName: responseElement.Attachment.DisplayName,
				FileURL:     responseElement.Attachment.FileURL,
			},
			ProviderID:          responseElement.ProviderID,
			ProviderDisplayName: responseElement.ProviderDisplayName,
			DisplayState:        responseElement.DisplayState,
		}
	}
	return elements, nil
}
