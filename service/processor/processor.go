package processor

import (
	"encoding/json"
	"fmt"
	"github.com/pennsieve/processor-post-metadata/client"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/logging"
	"github.com/pennsieve/processor-post-metadata/service/pennsieve"
	"log/slog"
	"os"
	"path/filepath"
)

var logger = logging.PackageLogger("processor")

type MetadataPostProcessor struct {
	IntegrationID   string
	InputDirectory  string
	OutputDirectory string
	Pennsieve       *pennsieve.Session
	IDStore         *IDStore
}

func NewMetadataPostProcessor(
	integrationID string,
	inputDirectory string,
	outputDirectory string,
	sessionToken string,
	apiHost string,
	api2Host string,
	idStore *IDStore) (*MetadataPostProcessor, error) {
	session := pennsieve.NewSession(sessionToken, apiHost, api2Host)
	return &MetadataPostProcessor{
		IntegrationID:   integrationID,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		Pennsieve:       session,
		IDStore:         idStore,
	}, nil
}

func (p *MetadataPostProcessor) Run() error {
	integration, err := p.Pennsieve.GetIntegration(p.IntegrationID)
	if err != nil {
		return fmt.Errorf("error getting integration %s from Pennsieve: %w", p.IntegrationID, err)
	}
	datasetID := integration.DatasetNodeID
	logger.Info("starting metadata processing", slog.String("datasetID", datasetID))
	datasetChanges, err := readChangesetFile(p.changesetFilePath())
	if err != nil {
		return err
	}
	logger.Info("read dataset changeset file", slog.String("path", p.changesetFilePath()))
	if err := p.ProcessModels(datasetID, datasetChanges.Models); err != nil {
		return err
	}
	if err := p.ProcessLinks(datasetID, datasetChanges.LinkedProperties); err != nil {
		return err
	}
	logger.Info("finished metadata processing")
	return nil
}

func (p *MetadataPostProcessor) changesetFilePath() string {
	return filepath.Join(p.OutputDirectory, client.Filename)
}

func readChangesetFile(filePath string) (clientmodels.Dataset, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return clientmodels.Dataset{}, fmt.Errorf("error opening changeset file %s: %w", filePath, err)
	}
	var datasetChangeset clientmodels.Dataset
	if err := json.NewDecoder(file).Decode(&datasetChangeset); err != nil {
		return clientmodels.Dataset{}, fmt.Errorf("error decoding changeset file %s: %w", filePath, err)
	}
	return datasetChangeset, nil
}
