package eseis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type coownershipFolderResponse struct {
	ID             int    `json:"id"`
	DisplayName    string `json:"display_name"`
	DocumentsCount int    `json:"documents_count"`
}

type CoownershipFolder struct {
	ID          int
	DisplayName string
}

func (e *EseisClient) GetCoownershipFolders(placeID int, page int) ([]CoownershipFolder, error) {
	if err := e.checkAuthenticated(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v2/places/%d/coownership_folders?page=%d&sort=display_name", placeID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create coownership_folders request: %w", err)
	}
	e.setAuthentication(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send coownership_folders request: %w", err)
	}
	defer resp.Body.Close()

	var foldersResponse []coownershipFolderResponse
	err = json.NewDecoder(resp.Body).Decode(&foldersResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode coownership_folders response: %w", err)
	}

	var folders = make([]CoownershipFolder, len(foldersResponse))
	for i, folderResponse := range foldersResponse {
		folders[i] = CoownershipFolder{ID: folderResponse.ID, DisplayName: folderResponse.DisplayName}
	}
	return folders, nil
}

type coownershipDocumentResponse struct {
	ID            int       `json:"id"`
	UUID          string    `json:"uuid"`
	DisplayName   string    `json:"display_name"`
	DocumentsName string    `json:"documents_name"`
	FileURL       string    `json:"file_url"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CoownershipDocument struct {
	ID            int
	UUID          string
	DisplayName   string
	DocumentsName string
	FileURL       string
	UpdatedAt     time.Time
}

func (e *EseisClient) GetCoownershipDocuments(placeID int, folderID int, page int) ([]CoownershipDocument, error) {
	if err := e.checkAuthenticated(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/places/%d/coownership_documents?by_folder=%d&page=%d", placeID, folderID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create coownership_documents request: %w", err)
	}
	e.setAuthentication(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send coownership_documents request: %w", err)
	}
	defer resp.Body.Close()

	var coownershipDocumentsResponse []coownershipDocumentResponse
	err = json.NewDecoder(resp.Body).Decode(&coownershipDocumentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode coownership_documents response: %w", err)
	}

	var coownershipDocuments = make([]CoownershipDocument, len(coownershipDocumentsResponse))
	for i, coownershipDocument := range coownershipDocumentsResponse {
		coownershipDocuments[i] = CoownershipDocument{
			ID:            coownershipDocument.ID,
			UUID:          coownershipDocument.UUID,
			DisplayName:   coownershipDocument.DisplayName,
			DocumentsName: coownershipDocument.DocumentsName,
			FileURL:       coownershipDocument.FileURL,
			UpdatedAt:     coownershipDocument.UpdatedAt,
		}
	}
	return coownershipDocuments, nil
}

type maintenanceContractCategoryResponse struct {
	ID                   int    `json:"id"`
	DisplayName          string `json:"display_name"`
	MaintenanceContracts []struct {
		ID          int    `json:"id"`
		Reference   string `json:"reference"`
		CompanyName string `json:"company_name"`
	} `json:"maintenance_contracts"`
	MaintenanceContractsCount int `json:"maintenance_contracts_count"`
}

type MaintenanceContractCategory struct {
	ID                        int
	DisplayName               string
	MaintenanceContracts      []MaintenanceContract
	MaintenanceContractsCount int
}

type MaintenanceContract struct {
	ID          int
	Reference   string
	CompanyName string
}

func (e *EseisClient) GetMaintenanceContractCategories(placeID int) ([]MaintenanceContractCategory, error) {
	if err := e.checkAuthenticated(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/places/%d/maintenance_contract_categories", placeID)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create maintenance_contract_categories request: %w", err)
	}
	e.setAuthentication(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send maintenance_contract_categories request: %w", err)
	}
	defer resp.Body.Close()

	var maintenanceContractCategoriesResponse []maintenanceContractCategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&maintenanceContractCategoriesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode maintenance_contract_categories response: %w", err)
	}

	var maintenanceContractCategories = make([]MaintenanceContractCategory, len(maintenanceContractCategoriesResponse))
	for i, maintenanceContractCategory := range maintenanceContractCategoriesResponse {
		maintenanceContracts := make([]MaintenanceContract, len(maintenanceContractCategory.MaintenanceContracts))
		for j, maintenanceContract := range maintenanceContractCategory.MaintenanceContracts {
			maintenanceContracts[j] = MaintenanceContract{
				ID:          maintenanceContract.ID,
				Reference:   maintenanceContract.Reference,
				CompanyName: maintenanceContract.CompanyName,
			}
		}
		maintenanceContractCategories[i] = MaintenanceContractCategory{
			ID:                   maintenanceContractCategory.ID,
			DisplayName:          maintenanceContractCategory.DisplayName,
			MaintenanceContracts: maintenanceContracts,
		}
	}
	return maintenanceContractCategories, nil
}

type maintenanceContractDetailsResponse struct {
	ID                           int         `json:"id"`
	Reference                    string      `json:"reference"`
	CompanyName                  string      `json:"company_name"`
	OriginDate                   time.Time   `json:"origin_date"`
	EffectiveDate                time.Time   `json:"effective_date"`
	DueDate                      time.Time   `json:"due_date"`
	ContractLengthInMonth        int         `json:"contract_length_in_month"`
	NoticeInMonth                int         `json:"notice_in_month"`
	ProviderPhoneNumber          interface{} `json:"provider_phone_number"`
	CreatedAt                    time.Time   `json:"created_at"`
	UpdatedAt                    time.Time   `json:"updated_at"`
	MaintenanceContractDocuments []struct {
		ID           int       `json:"id"`
		UUID         string    `json:"uuid"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		DisplayName  string    `json:"display_name"`
		DocumentName string    `json:"document_name"`
		FileURL      string    `json:"file_url"`
	} `json:"maintenance_contract_documents"`
}

type MaintenanceContractDetails struct {
	ID                           int
	Reference                    string
	CompanyName                  string
	UpdatedAt                    time.Time
	MaintenanceContractDocuments []MaintenanceContractDocument
}

type MaintenanceContractDocument struct {
	ID           int
	UUID         string
	DisplayName  string
	DocumentName string
	UpdatedAt    time.Time
	FileURL      string
}

func (e *EseisClient) GetMaintenanceContractDetails(maintenanceContractID int) (MaintenanceContractDetails, error) {
	if err := e.checkAuthenticated(); err != nil {
		return MaintenanceContractDetails{}, err
	}

	path := fmt.Sprintf("/v1/maintenance_contracts/%d", maintenanceContractID)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return MaintenanceContractDetails{}, fmt.Errorf("failed to create maintenance_contract request: %w", err)
	}
	e.setAuthentication(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MaintenanceContractDetails{}, fmt.Errorf("failed to send maintenance_contract request: %w", err)
	}
	defer resp.Body.Close()

	var response maintenanceContractDetailsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return MaintenanceContractDetails{}, fmt.Errorf("failed to decode maintenance_contract response: %w", err)
	}

	maintenanceContractDocuments := make([]MaintenanceContractDocument, len(response.MaintenanceContractDocuments))
	for j, maintenanceContractDocument := range response.MaintenanceContractDocuments {
		maintenanceContractDocuments[j] = MaintenanceContractDocument{
			ID:           maintenanceContractDocument.ID,
			UUID:         maintenanceContractDocument.UUID,
			DisplayName:  maintenanceContractDocument.DisplayName,
			DocumentName: maintenanceContractDocument.DocumentName,
			UpdatedAt:    maintenanceContractDocument.UpdatedAt,
			FileURL:      maintenanceContractDocument.FileURL,
		}
	}
	maintenanceContractDetails := MaintenanceContractDetails{
		ID:                           response.ID,
		Reference:                    response.Reference,
		CompanyName:                  response.CompanyName,
		UpdatedAt:                    response.UpdatedAt,
		MaintenanceContractDocuments: maintenanceContractDocuments,
	}
	return maintenanceContractDetails, nil
}
