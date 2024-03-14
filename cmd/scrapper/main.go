package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v7"
	"github.com/idkw/eseisscrapper/pkg/infrastructure/eseis"
	"github.com/idkw/eseisscrapper/pkg/infrastructure/utils"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type config struct {
	OutDir string `env:"ESEIS_SCRAPPER_OUT_DIR,required"`
}

const (
	sergicOffer            = "ESE"
	individualDir          = "individual"
	coownershipDir         = "coownership"
	maintenanceDir         = "maintenance"
	reportsDir             = "reports"
	reportsOpenedDir       = "opened"
	reportsAcknowledgedDir = "acknowledged"
	reportsResolvedDir     = "resolved"
	forumTopicsDir         = "forum"
	budgetsDir             = "budgets"
	pdfFileExtension       = ".pdf"
	jpgFileExtension       = ".jpg"
)

func main() {
	config, err := newConfig()
	utils.MustBeNilErr(err, "failed to create config")
	client := eseis.NewEseisClientFatal()
	exportContracts(client, config.OutDir)
	logrus.Info("Done scrapping Eseis documents")
}

func exportContracts(client *eseis.EseisClient, outDir string) {
	utils.MkDirFatal(outDir)

	contracts, err := client.GetContracts(sergicOffer)
	utils.MustBeNilErr(err, "failed to get contracts for sergicOffer %s", sergicOffer)

	for _, contract := range contracts {
		contractOutDir := utils.JoinFilePath(outDir, utils.SanitizePath(contract.DisplayName))
		logrus.Infof("processing contract %d - %s", contract.ID, contract.DisplayName)
		exportIndividualDocuments(client, contract, contractOutDir)
		exportCoownershipDocuments(client, contract, contractOutDir)
		exportMaintenanceContractDocuments(client, contract, contractOutDir)
		exportReports(client, contract, contractOutDir)
		exportForumTopics(client, contract, contractOutDir)
		exportBudgets(client, contract, contractOutDir)
	}
}

func exportIndividualDocuments(client *eseis.EseisClient, contract eseis.Contract, outDir string) {
	foldersPage := 1
	for {
		folders, err := client.GetContractFolders(contract.ID, foldersPage)
		utils.MustBeNilErr(err, "failed to get contract folders for id=%d page=%d", contract.ID, foldersPage)
		if len(folders) == 0 {
			break
		}

		for _, folder := range folders {
			logrus.Infof("----------\nFolder %d:%s", folder.ID, folder.DisplayName)

			folderPath := utils.JoinFilePath(outDir, individualDir, utils.SanitizePath(folder.DisplayName))
			utils.MkDirFatal(folderPath)

			documentsPage := 1
			for {
				documents, err := client.GetContractDocuments(contract.ID, folder.ID, documentsPage)
				utils.MustBeNilErr(err, "failed to get contract documents for id=%d folder=%d, page=%d", contract.ID, folder.ID, foldersPage)
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
		utils.MustBeNilErr(err, "failed to get coownership folders for placeId=%d page=%d", contract.PlaceID, foldersPage)
		if len(coownershipFolders) == 0 {
			break
		}

		for _, coownershipFolder := range coownershipFolders {
			logrus.Infof("----------\nCoownershipFolder %d:%s", coownershipFolder.ID, coownershipFolder.DisplayName)

			folderPath := utils.JoinFilePath(outDir, coownershipDir, utils.SanitizePath(coownershipFolder.DisplayName))
			utils.MkDirFatal(folderPath)

			documentsPage := 1
			for {
				documents, err := client.GetCoownershipDocuments(contract.PlaceID, coownershipFolder.ID, documentsPage)
				utils.MustBeNilErr(err, "failed to get coownership documents for placeId=%d coownershipFolder=%d, page=%d", contract.PlaceID, coownershipFolder.ID, foldersPage)
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
	utils.MustBeNilErr(err, "failed to get maintenance contract categories for placeID %d", contract.PlaceID)
	for _, category := range categories {
		categoryFolderPath := utils.JoinFilePath(outDir, maintenanceDir, utils.SanitizePath(category.DisplayName))
		for _, maintenanceContract := range category.MaintenanceContracts {
			maintenanceContractFolderPath := utils.JoinFilePath(categoryFolderPath, utils.SanitizePath(maintenanceContract.CompanyName+"_"+maintenanceContract.Reference))
			utils.MkDirFatal(maintenanceContractFolderPath)

			maintenanceContractDetails, err := client.GetMaintenanceContractDetails(maintenanceContract.ID)
			utils.MustBeNilErr(err, "failed to get maintenance contract details for id %d", maintenanceContract.ID)
			for _, document := range maintenanceContractDetails.MaintenanceContractDocuments {
				exportDocument(client, document.UUID, document.DisplayName, document.UpdatedAt, maintenanceContractFolderPath)
			}

			// add additional info file for metadata
			exportInfoFile(maintenanceContractDetails, utils.JoinFilePath(maintenanceContractFolderPath, "info.json"))
		}
	}
}

func exportReports(client *eseis.EseisClient, contract eseis.Contract, outDir string) {
	utils.MkDirFatal(utils.JoinFilePath(outDir, reportsDir, reportsOpenedDir))
	utils.MkDirFatal(utils.JoinFilePath(outDir, reportsDir, reportsAcknowledgedDir))
	utils.MkDirFatal(utils.JoinFilePath(outDir, reportsDir, reportsResolvedDir))

	reportsPage := 1
	for {
		reports, err := client.GetReports(contract.PlaceID, reportsPage)
		utils.MustBeNilErr(err, "failed to get contract folders for placeId=%d page=%d", contract.PlaceID, reportsPage)
		if len(reports) == 0 {
			break
		}

		for _, report := range reports {
			logrus.Infof("----------\nReport %d:%s", report.ID, report.DisplayName)
			err := client.CreateReportScreenshot(report, utils.JoinFilePath(outDir, reportsDir, report.State))
			utils.MustBeNilErr(err, "failed screenshot for report %d", report.ID)
		}

		reportsPage++
	}
}

func exportForumTopics(client *eseis.EseisClient, contract eseis.Contract, outDir string) {
	utils.MkDirFatal(utils.JoinFilePath(outDir, forumTopicsDir))

	page := 1
	for {
		forumTopics, err := client.GetForumTopics(contract.PlaceID, page)
		utils.MustBeNilErr(err, "failed to get forum topics for placeId=%d page=%d", contract.PlaceID, page)
		if len(forumTopics) == 0 {
			break
		}

		for _, forumTopic := range forumTopics {
			logrus.Infof("----------\nReport %d:%s", forumTopic.ID, forumTopic.DisplayName)

			year, month, day := forumTopic.CreatedAt.Date()
			forumTopicDir := utils.JoinFilePath(
				outDir,
				forumTopicsDir,
				fmt.Sprintf(
					"%d_%d_%d__%d__%s", year, month, day, forumTopic.ID, forumTopic.CleanDisplayName(),
				))
			utils.MkDirFatal(forumTopicDir)
			exportInfoFile(forumTopic.Raw, utils.JoinFilePath(forumTopicDir, "topic.json"))
			for _, attachment := range forumTopic.Raw.Attachments {
				exportAttachment(client, attachment.FileURL, attachment.SourceFileName, attachment.SourceContentType, attachment.SourceUpdatedAt, forumTopicDir)
			}

			topicPosts, err := client.GetAllTopicPosts(contract.PlaceID, forumTopic.ID)
			utils.MustBeNilErr(err, "failed to get topic posts for placeID=%d forumTopic=%d", contract.PlaceID, forumTopic.ID)
			exportInfoFile(topicPosts, utils.JoinFilePath(forumTopicDir, "posts.json"))
			for _, post := range topicPosts {
				for _, attachment := range post.Attachments {
					exportAttachment(client, attachment.FileURL, attachment.SourceFileName, attachment.SourceContentType, attachment.SourceUpdatedAt, forumTopicDir)
				}
			}

		}

		page++
	}
}

func exportBudgets(client *eseis.EseisClient, contract eseis.Contract, outDir string) {
	utils.MkDirFatal(utils.JoinFilePath(outDir, budgetsDir))

	fiscalYears, err := client.GetFiscalYears(contract.PlaceID)
	utils.MustBeNilErr(err, "failed to get fiscal years for placeId=%d", contract.PlaceID)

	for _, fiscalYear := range fiscalYears {
		budgets, err := client.GetBudgets(contract.PlaceID, fiscalYear.ID)
		utils.MustBeNilErr(err, "failed to get budgets for placeId=%d and fiscalYear=%d", contract.PlaceID, fiscalYear.ID)
		fiscalYearDirName := utils.JoinFilePath(
			outDir,
			budgetsDir,
			strings.ReplaceAll(fiscalYear.DisplayName, "/", "_"),
		)
		utils.MkDirFatal(fiscalYearDirName)
		exportInfoFile(fiscalYear, utils.JoinFilePath(fiscalYearDirName, "info.json"))

		for _, budget := range budgets {
			budgetDirName := utils.JoinFilePath(
				fiscalYearDirName,
				utils.SanitizePath(budget.DisplayName),
			)
			utils.MkDirFatal(budgetDirName)
			exportInfoFile(budget, utils.JoinFilePath(budgetDirName, "info.json"))

			accountPlaceEntries, err := client.GetAccountPlaceEntries(budget.ID)
			utils.MustBeNilErr(err, "failed to get account place entries for budgetID=%d", budget.ID)
			for _, accountPlaceEntry := range accountPlaceEntries {
				exportDocumentName := fmt.Sprintf(
					"%s_%d_%s",
					accountPlaceEntry.OperationDate.Format(time.RFC3339),
					accountPlaceEntry.Amount,
					utils.SanitizePath(accountPlaceEntry.DisplayName),
				)
				exportInfoFile(accountPlaceEntry, utils.JoinFilePath(budgetDirName, fmt.Sprintf("%s.json", exportDocumentName)))
				exportDocument(client, accountPlaceEntry.UUID, exportDocumentName, accountPlaceEntry.UpdatedAt, budgetDirName)
			}
		}
	}
}

func exportDocument(client *eseis.EseisClient, documentUUID string, documentName string, updatedAt time.Time, folderPath string) {
	logrus.Infof("Exporting document %s:%s to folder %s", documentUUID, documentName, folderPath)

	documentFilePath := utils.JoinFilePath(folderPath, utils.SanitizePath(documentName+pdfFileExtension))

	fileInfo, err := os.Stat(documentFilePath)
	if errors.Is(err, os.ErrNotExist) {
		// nominal case, new unseen file
	} else {
		// already seen file
		utils.MustBeNilErr(err, "failed to stat file at path %s", documentFilePath)
		if updatedAt.Before(fileInfo.ModTime()) {
			logrus.Infof("document %s:%s already downloaded", documentUUID, documentName)
			return
		}
	}

	documentFile, err := os.Create(documentFilePath)
	utils.MustBeNilErr(err, "failed to create output document at path %s", documentFilePath)
	defer documentFile.Close()

	documentBytes, err := client.GetDocument(documentUUID)
	utils.MustBeNilErr(err, "failed to get document for uuid %s", documentUUID)
	_, err = documentFile.Write(documentBytes)
	utils.MustBeNilErr(err, "failed to write output document at path %s", documentFilePath)
}

func exportAttachment(client *eseis.EseisClient, url string, attachmentName string, attachmentFileType string, updatedAt time.Time, folderPath string) {
	logrus.Infof("Exporting attachment %s:%s to folder %s", url, attachmentName, folderPath)

	fileExtension := ""
	switch attachmentFileType {
	case "application/pdf":
		fileExtension = pdfFileExtension
		break
	case "image/jpeg":
		fileExtension = jpgFileExtension
	}
	attachmentFilePath := utils.JoinFilePath(folderPath, utils.SanitizePath(attachmentName+fileExtension))

	fileInfo, err := os.Stat(attachmentFilePath)
	if errors.Is(err, os.ErrNotExist) {
		// nominal case, new unseen file
	} else {
		// already seen file
		utils.MustBeNilErr(err, "failed to stat file at path %s", attachmentFilePath)
		if updatedAt.Before(fileInfo.ModTime()) {
			logrus.Infof("attachment %s:%s already downloaded", url, attachmentName)
			return
		}
	}

	attachmentFile, err := os.Create(attachmentFilePath)
	utils.MustBeNilErr(err, "failed to create output attachment at path %s", attachmentFilePath)
	defer attachmentFile.Close()

	attachmentBytes, err := client.GetAttachment(url)
	utils.MustBeNilErr(err, "failed to get attachment for url %s", url)
	_, err = attachmentFile.Write(attachmentBytes)
	utils.MustBeNilErr(err, "failed to write output attachment at path %s", attachmentFilePath)
}

func exportInfoFile(content any, infoFilePath string) {
	// add aditional info file for metadata
	contentJson, err := json.MarshalIndent(content, "", "  ")
	utils.MustBeNilErr(err, "failed to serialize info file for %+v", content)
	err = os.WriteFile(infoFilePath, contentJson, 0660)
	utils.MustBeNilErr(err, "failed to write info file for id=%+v", content)
}

func newConfig() (*config, error) {
	config := &config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment: %w", err)
	}
	return config, nil
}
