package processor

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurationExportSyncProcessor_Run(t *testing.T) {
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := t.TempDir()
	sessionToken := uuid.NewString()

	datasetID := newDatasetID()
	mockServer := newMockServer(t, integrationID, datasetID)
	defer mockServer.Close()

	processor, err := NewMetadataPostProcessor(integrationID, inputDirectory, outputDirectory, sessionToken, mockServer.URL, mockServer.URL)
	require.NoError(t, err)

	require.NoError(t, processor.Run())
}

func newDatasetID() string {
	return newDatasetIDWithUUID(uuid.NewString())
}

func newDatasetIDWithUUID(datasetUUID string) string {
	return fmt.Sprintf("N:dataset:%s", datasetUUID)
}

func newMockServer(t *testing.T, integrationID string, datasetID string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/integrations/%s", integrationID), func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, http.MethodGet, request.Method, "expected method %s for %s, got %s", http.MethodGet, request.URL, request.Method)
		integration := models.Integration{
			Uuid:          uuid.NewString(),
			ApplicationID: 0,
			DatasetNodeID: datasetID,
		}
		integrationResponse, err := json.Marshal(integration)
		require.NoError(t, err)
		_, err = writer.Write(integrationResponse)
		require.NoError(t, err)
	})
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		require.Fail(t, "unexpected call to Pennsieve", "%s %s", request.Method, request.URL)
	})
	return httptest.NewServer(mux)
}
