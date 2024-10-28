package processor_test

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/processor"
	"github.com/pennsieve/processor-post-metadata/service/processor/internal/processortest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCurationExportSyncProcessor_Run(t *testing.T) {
	t.Skip("need some changes to the curation file in testdata")
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := "testdata/output"
	sessionToken := uuid.NewString()

	datasetID := processortest.NewDatasetID()
	expectedIntegrationsCall := mock.NewExpectedGetIntegrationCall(integrationID, datasetID)
	mockServer := mock.NewModelService(t, expectedIntegrationsCall)
	defer mockServer.Close()

	idStore := processor.NewIDStoreBuilder().Build()

	testProcessor, err := processor.NewMetadataPostProcessor(integrationID, inputDirectory, outputDirectory, sessionToken, mockServer.URL(), mockServer.URL(), idStore)
	require.NoError(t, err)

	require.NoError(t, testProcessor.Run())
}
