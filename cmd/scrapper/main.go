package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v7"
	"github.com/idkw/eseisscrapper/pkg/infrastructure/eseis"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	OutDir string `env:"ESEIS_SCRAPPER_OUT_DIR,required"`
}

const (
	sergicOffer    = "ESE"
	individualDir  = "individual"
	coownershipDir = "coownership"
	maintenanceDir = "maintenance"
	infoFile       = "info.txt"
	fileExtension  = ".pdf"
)

func main() {
	config, err := newConfig()
	mustBeNilErr(err, "failed to create config")
	client := eseis.NewEseisClientFatal()
	exportContracts(client, config.OutDir)
	logrus.Info("Done scrapping Eseis documents")
}

func exportContracts(client *eseis.EseisClient, outDir string) {
	mkDirFatal(outDir)

	contracts, err := client.GetContracts(sergicOffer)
	mustBeNilErr(err, "failed to get contracts for sergicOffer %s", sergicOffer)

	for _, contract := range contracts {
		contractOutDir := joinFilePath(outDir, sanitizePath(contract.DisplayName))
		logrus.Infof("processing contract %d - %s", contract.ID, contract.DisplayName)
		exportIndividualDocuments(client, contract, contractOutDir)
		exportCoownershipDocuments(client, contract, contractOutDir)
		exportMaintenanceContractDocuments(client, contract, contractOutDir)
	}
}

func exportIndividualDocuments(client *eseis.EseisClient, contract eseis.Contract, outDir string) {
	foldersPage := 1
	for {
		folders, err := client.GetContractFolders(contract.ID, foldersPage)
		mustBeNilErr(err, "failed to get contract folders for id=%d page=%d", contract.ID, foldersPage)
		if len(folders) == 0 {
			break
		}

		for _, folder := range folders {
			logrus.Infof("----------\nFolder %d:%s", folder.ID, folder.DisplayName)

			folderPath := joinFilePath(outDir, individualDir, sanitizePath(folder.DisplayName))
			mkDirFatal(folderPath)

			documentsPage := 1
			for {
				documents, err := client.GetContractDocuments(contract.ID, folder.ID, documentsPage)
				mustBeNilErr(err, "failed to get contract documents for id=%d folder=%d, page=%d", contract.ID, folder.ID, foldersPage)
				if len(documents) == 0 {
					break
				}
				for _, document := range documents {
					exportDocument(client, document.UUID, document.DisplayName, document.UpdatedAt, folderPath)
				}
				documentsPage++
			}
		}

		foldersPage++
	}
}

func exportCoownershipDocuments(client *eseis.EseisClient, contract eseis.Contract, outDir string) {
	foldersPage := 1
	for {
		coownershipFolders, err := client.GetCoownershipFolders(contract.PlaceID, foldersPage)
		mustBeNilErr(err, "failed to get coownership folders for placeId=%d page=%d", contract.PlaceID, foldersPage)
		if len(coownershipFolders) == 0 {
			break
		}

		for _, coownershipFolder := range coownershipFolders {
			logrus.Infof("----------\nCoownershipFolder %d:%s", coownershipFolder.ID, coownershipFolder.DisplayName)

			folderPath := joinFilePath(outDir, coownershipDir, sanitizePath(coownershipFolder.DisplayName))
			mkDirFatal(folderPath)

			documentsPage := 1
			for {
				documents, err := client.GetCoownershipDocuments(contract.PlaceID, coownershipFolder.ID, documentsPage)
				mustBeNilErr(err, "failed to get coownership documents for placeId=%d coownershipFolder=%d, page=%d", contract.PlaceID, coownershipFolder.ID, foldersPage)
				if len(documents) == 0 {
					break
				}
				for _, document := range documents {
					exportDocument(client, document.UUID, document.DisplayName, document.UpdatedAt, folderPath)
				}
				documentsPage++
			}
		}

		foldersPage++
	}
}

func exportMaintenanceContractDocuments(client *eseis.EseisClient, contract eseis.Contract, outDir string) {
	categories, err := client.GetMaintenanceContractCategories(contract.PlaceID)
	mustBeNilErr(err, "failed to get maintenance contract categories for placeID %d", contract.PlaceID)
	for _, category := range categories {
		categoryFolderPath := joinFilePath(outDir, maintenanceDir, sanitizePath(category.DisplayName))
		for _, maintenanceContract := range category.MaintenanceContracts {
			maintenanceContractFolderPath := joinFilePath(categoryFolderPath, sanitizePath(maintenanceContract.CompanyName+"_"+maintenanceContract.Reference))
			mkDirFatal(maintenanceContractFolderPath)

			maintenanceContractDetails, err := client.GetMaintenanceContractDetails(maintenanceContract.ID)
			mustBeNilErr(err, "failed to get maintenance contract details for id %d", maintenanceContract.ID)
			for _, document := range maintenanceContractDetails.MaintenanceContractDocuments {
				exportDocument(client, document.UUID, document.DisplayName, document.UpdatedAt, maintenanceContractFolderPath)
			}

			// add additional info file for metadata
			exportInfoFile(maintenanceContractDetails, maintenanceContractFolderPath)
		}
	}
}

func exportDocument(client *eseis.EseisClient, documentUUID string, documentName string, updatedAt time.Time, folderPath string) {
	logrus.Infof("Exporting document %s:%s to folder %s", documentUUID, documentName, folderPath)

	documentFilePath := joinFilePath(folderPath, sanitizePath(documentName+fileExtension))

	fileInfo, err := os.Stat(documentFilePath)
	if errors.Is(err, os.ErrNotExist) {
		// nominal case, new unseen file
	} else {
		// already seen file
		mustBeNilErr(err, "failed to stat file at path %s", documentFilePath)
		if updatedAt.Before(fileInfo.ModTime()) {
			logrus.Infof("document %s:%s already downloaded", documentUUID, documentName)
			return
		}
	}

	documentFile, err := os.Create(documentFilePath)
	mustBeNilErr(err, "failed to create output document at path %s", documentFilePath)
	defer documentFile.Close()

	documentBytes, err := client.GetDocument(documentUUID)
	mustBeNilErr(err, "failed to get document for uuid %s", documentUUID)
	_, err = documentFile.Write(documentBytes)
	mustBeNilErr(err, "failed to write output document at path %s", documentFilePath)
}

func exportInfoFile(content any, folderPath string) {
	// add aditional info file for metadata
	contentJson, err := json.MarshalIndent(content, "", "  ")
	mustBeNilErr(err, "failed to serialize info file for %+v", content)
	err = os.WriteFile(joinFilePath(folderPath, infoFile), contentJson, 0660)
	mustBeNilErr(err, "failed to write info file for id=%+v", content)
}

func mustBeNilErr(err error, message string, args ...interface{}) {
	if err != nil {
		logrus.Fatalf(message+": %s", err, args)
	}
}

func mkDirFatal(path string) {
	err := os.MkdirAll(path, 0770)
	mustBeNilErr(err, "failed to create dir %s", path)
}

func joinFilePath(elements ...string) string {
	return filepath.Join(elements...)
}

func sanitizePath(filePath string) string {
	return strings.ReplaceAll(filePath, "/", "_")
}
func newConfig() (*config, error) {
	config := &config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment: %w", err)
	}
	return config, nil
}
