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
	// initialize the IDStore with model name -> id map for existing models
	// If we create models in this changeset, those name -> id entries will be added as well
	p.IDStore.AddModels(datasetChanges.ExistingModelIDMap)
	if err := p.ProcessDeletes(datasetID, datasetChanges); err != nil {
		return err
	}
	if err := p.ProcessModelCreatesUpdates(datasetID, datasetChanges.Models.Creates, datasetChanges.Models.Updates); err != nil {
		return err
	}
	// Wait til after ProcessModels to add these so that the IDStore now should have the complete mapping
	// of model names to model IDs for any models that were created.
	// ProcessLinks and ProcessProxies will need these record ID maps.
	if err := p.IDStore.AddRecordIDMaps(datasetChanges.RecordIDMaps); err != nil {
		return err
	}
	if err := p.ProcessLinks(datasetID, datasetChanges.LinkedProperties); err != nil {
		return err
	}
	if err := p.ProcessProxyChanges(datasetID, datasetChanges.Proxies); err != nil {
		return err
	}
	logger.Info("finished metadata processing")
	return nil
}

func (p *MetadataPostProcessor) ProcessDeletes(datasetID string, datasetChanges clientmodels.Dataset) error {
	// Delete dependent objects, links and proxies before deleting records
	logger.Info("starting deletes")
	if err := p.ProcessLinkInstanceDeletes(datasetID, datasetChanges.LinkedProperties); err != nil {
		return err
	}
	if proxyChanges := datasetChanges.Proxies; proxyChanges != nil {
		if err := p.ProcessProxyInstanceDeletes(datasetID, *proxyChanges); err != nil {
			return err
		}
	}
	if err := p.ProcessRecordModelDeletes(datasetID, datasetChanges.Models.Updates, datasetChanges.Models.Deletes); err != nil {
		return err
	}
	logger.Info("finished deletes")
	return nil
}

// ChangesetFilePath joins the given output directory with the
// changeset file name.
// Visible for testing.
func ChangesetFilePath(outputDirectory string) string {
	return filepath.Join(outputDirectory, client.Filename)
}

func (p *MetadataPostProcessor) changesetFilePath() string {
	return ChangesetFilePath(p.OutputDirectory)
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
