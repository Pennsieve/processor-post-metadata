package processor

import (
	"fmt"
	"github.com/pennsieve/processor-post-metadata/service/logging"
	"github.com/pennsieve/processor-post-metadata/service/pennsieve"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"log/slog"
)

var processorLogger = logging.PackageLogger("processor")

type MetadataPostProcessor struct {
	IntegrationID   string
	InputDirectory  string
	OutputDirectory string
	MetadataReader  *metadataclient.Reader
	Pennsieve       *pennsieve.Session
}

func NewMetadataPostProcessor(
	integrationID string,
	inputDirectory string,
	outputDirectory string,
	sessionToken string,
	apiHost string,
	api2Host string) (*MetadataPostProcessor, error) {
	reader, err := metadataclient.NewReader(inputDirectory)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata reader for %s: %w", inputDirectory, err)
	}
	session := pennsieve.NewSession(sessionToken, apiHost, api2Host)
	return &MetadataPostProcessor{
		IntegrationID:   integrationID,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		Pennsieve:       session,
		MetadataReader:  reader,
	}, nil
}

func (p *MetadataPostProcessor) Run() error {
	integration, err := p.Pennsieve.GetIntegration(p.IntegrationID)
	if err != nil {
		return fmt.Errorf("error getting integration %s from Pennsieve: %w", p.IntegrationID, err)
	}
	logger := processorLogger.With(slog.String("datasetID", integration.DatasetNodeID))
	logger.Info("starting metadata processing")
	logger.Info("finished metadata processing")
	return nil
}
