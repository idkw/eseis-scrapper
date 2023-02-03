package eseis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type contractResponse struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
	PlaceID     int    `json:"place_id"`
}

type Contract struct {
	ID          int
	DisplayName string
	PlaceID     int
}

func (e *EseisClient) GetContracts(sergicOffer string) ([]Contract, error) {
	if err := e.checkAuthenticated(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/users/me/contracts?by_sergic_offer=%s", sergicOffer)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create contracts request: %w", err)
	}
	e.setAuthentication(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send contracts request: %w", err)
	}
	defer resp.Body.Close()

	var contractsResponse []contractResponse
	err = json.NewDecoder(resp.Body).Decode(&contractsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contracts response: %w", err)
	}

	var contracts = make([]Contract, len(contractsResponse))
	for i, contractResponse := range contractsResponse {
		contracts[i] = Contract{
			ID:          contractResponse.ID,
			DisplayName: contractResponse.DisplayName,
			PlaceID:     contractResponse.PlaceID,
		}
	}
	return contracts, nil
}

//func (e *EseisClient) GetContracts(sergicOffer string, placeID int) ([]Contract, error) {
//	if err := e.checkAuthenticated(); err != nil {
//		return nil, err
//	}
//
//	path := fmt.Sprintf("/v1/users/me/contracts?by_sergic_offer=%s&by_place=%d", sergicOffer, placeID)
//	req, err := http.NewRequest("GET", e.buildURL(path), nil)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create contracts request: %w", err)
//	}
//	e.setAuthentication(req)
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("failed to send contracts request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	var contractsResponse []contractResponse
//	err = json.NewDecoder(resp.Body).Decode(&contractsResponse)
//	if err != nil {
//		return nil, fmt.Errorf("failed to decode contracts response: %w", err)
//	}
//
//	var contracts = make([]Contract, len(contractsResponse))
//	for i, contractResponse := range contractsResponse {
//		contracts[i] = Contract{
//			ID:          contractResponse.ID,
//			DisplayName: contractResponse.DisplayName,
//			PlaceID:     contractResponse.PlaceID,
//		}
//	}
//	return contracts, nil
//}

type contractFolderResponse struct {
	ID             int    `json:"id"`
	DisplayName    string `json:"display_name"`
	DocumentsCount int    `json:"documents_count"`
}

type ContractFolder struct {
	ID          int
	DisplayName string
}

func (e *EseisClient) GetContractFolders(contractID int, page int) ([]ContractFolder, error) {
	if err := e.checkAuthenticated(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v2/contract_folders?by_contract=%d&page=%d&sort=display_name", contractID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract_folders request: %w", err)
	}
	e.setAuthentication(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send contract_folders request: %w", err)
	}
	defer resp.Body.Close()

	var foldersResponse []contractFolderResponse
	err = json.NewDecoder(resp.Body).Decode(&foldersResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contract_folders response: %w", err)
	}

	var folders = make([]ContractFolder, len(foldersResponse))
	for i, folderResponse := range foldersResponse {
		folders[i] = ContractFolder{ID: folderResponse.ID, DisplayName: folderResponse.DisplayName}
	}
	return folders, nil
}

type contractDocumentResponse struct {
	ID            int       `json:"id"`
	UUID          string    `json:"uuid"`
	DisplayName   string    `json:"display_name"`
	DocumentsName string    `json:"documents_name"`
	FileURL       string    `json:"file_url"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ContractDocument struct {
	ID            int
	UUID          string
	DisplayName   string
	DocumentsName string
	FileURL       string
	UpdatedAt     time.Time
}

func (e *EseisClient) GetContractDocuments(contractID int, folderID, page int) ([]ContractDocument, error) {
	if err := e.checkAuthenticated(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/contracts/%d/contract_documents?by_folder=%d&page=%d", contractID, folderID, page)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract_documents request: %w", err)
	}
	e.setAuthentication(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send contract_documents request: %w", err)
	}
	defer resp.Body.Close()

	var contractDocumentsResponse []contractDocumentResponse
	err = json.NewDecoder(resp.Body).Decode(&contractDocumentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contract_documents response: %w", err)
	}

	var contractDocuments = make([]ContractDocument, len(contractDocumentsResponse))
	for i, contractDocument := range contractDocumentsResponse {
		contractDocuments[i] = ContractDocument{
			ID:            contractDocument.ID,
			UUID:          contractDocument.UUID,
			DisplayName:   contractDocument.DisplayName,
			DocumentsName: contractDocument.DocumentsName,
			FileURL:       contractDocument.FileURL,
			UpdatedAt:     contractDocument.UpdatedAt,
		}
	}
	return contractDocuments, nil
}
