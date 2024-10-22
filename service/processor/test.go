package processor

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ExpectedCall interface {
	CallCount() int
	PathHandler(t *testing.T) (string, http.HandlerFunc)
	Signature() string
}

type ExpectedAPICall[IN, OUT any] struct {
	Method              string
	APIPath             string
	ExpectedRequestBody *IN
	APIResponse         OUT
	callCount           int
}

func NewExpectedModelCreateCall(expectedDatasetID string, modelID string, expectedCreate clientmodels.ModelCreate) *ExpectedAPICall[clientmodels.ModelCreate, models.APIResponse] {
	return &ExpectedAPICall[clientmodels.ModelCreate, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts", expectedDatasetID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models.APIResponse{
			Name: expectedCreate.Name,
			ID:   modelID,
		},
	}
}

func NewExpectedPropertiesCreateCall(expectedDatasetID string, modelID string, expectedCreate clientmodels.PropertiesCreate) *ExpectedAPICall[clientmodels.PropertiesCreate, []models.APIResponse] {
	var apiResponse []models.APIResponse
	for _, prop := range expectedCreate {
		apiResponse = append(apiResponse, models.APIResponse{
			Name: prop.Name,
			ID:   uuid.NewString(),
		})
	}
	return &ExpectedAPICall[clientmodels.PropertiesCreate, []models.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/properties", expectedDatasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse:         apiResponse,
	}
}

func NewExpectedRecordCreateCall(datasetID string, modelID string, expectedCreate clientmodels.RecordValues) *ExpectedAPICall[clientmodels.RecordValues, models.APIResponse] {
	return &ExpectedAPICall[clientmodels.RecordValues, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}

func NewExpectedRecordUpdateCall(datasetID string, modelID string, recordID clientmodels.PennsieveInstanceID, expectedUpdate clientmodels.RecordValues) *ExpectedAPICall[clientmodels.RecordValues, models.APIResponse] {
	return &ExpectedAPICall[clientmodels.RecordValues, models.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances/%s", datasetID, modelID, recordID),
		ExpectedRequestBody: &expectedUpdate,
		APIResponse: models.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}

func (e *ExpectedAPICall[I, _]) HandlerFunction(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		e.callCount += 1
		require.Equal(t, e.Method, request.Method, "expected method %s for %s, got %s", e.Method, request.URL, request.Method)
		if e.ExpectedRequestBody == nil {
			var bytes []byte
			_, err := request.Body.Read(bytes)
			require.ErrorIs(t, err, io.EOF)
		} else {
			var actualRequestBody I
			require.NoError(t, json.NewDecoder(request.Body).Decode(&actualRequestBody))
			require.Equal(t, *e.ExpectedRequestBody, actualRequestBody)
		}
		responseBytes, err := json.Marshal(e.APIResponse)
		require.NoError(t, err)
		_, err = writer.Write(responseBytes)
		require.NoError(t, err)
	}
}

func (e *ExpectedAPICall[_, _]) CallCount() int {
	return e.callCount
}

func (e *ExpectedAPICall[_, _]) PathHandler(t *testing.T) (string, http.HandlerFunc) {
	return e.APIPath, e.HandlerFunction(t)
}

func (e *ExpectedAPICall[_, _]) Signature() string {
	return fmt.Sprintf("%s %s", e.Method, e.APIPath)
}

type MockModelServer struct {
	Server        *httptest.Server
	ExpectedCalls []ExpectedCall
}

func NewMockModelServer(t *testing.T, expectedCall ...ExpectedCall) *MockModelServer {
	mux := http.NewServeMux()
	for _, ph := range expectedCall {
		mux.HandleFunc(ph.PathHandler(t))
	}
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		require.Fail(t, "unexpected call to Pennsieve", "%s %s", request.Method, request.URL)
	})
	return &MockModelServer{
		Server:        httptest.NewServer(mux),
		ExpectedCalls: expectedCall,
	}
}

func (m *MockModelServer) AssertAllCalledExactlyOnce(t *testing.T) bool {
	for _, expectedCall := range m.ExpectedCalls {
		if !assert.Equal(t, 1, expectedCall.CallCount(), "%s was called %d times", expectedCall.Signature(), expectedCall.CallCount()) {
			return false
		}
	}
	return true
}

func (m *MockModelServer) Close() {
	m.Server.Close()
}

func (m *MockModelServer) URL() string {
	return m.Server.URL
}

type TestProcessorBuilder struct {
	integrationID   *string
	inputDirectory  *string
	outputDirectory *string
	sessionToken    *string
}

func NewTestProcessorBuilder() *TestProcessorBuilder {
	return &TestProcessorBuilder{}
}

func (b *TestProcessorBuilder) WithIntegrationID(integrationID string) *TestProcessorBuilder {
	b.integrationID = &integrationID
	return b
}

func (b *TestProcessorBuilder) WithInputDirectory(inputDirectory string) *TestProcessorBuilder {
	b.inputDirectory = &inputDirectory
	return b
}

func (b *TestProcessorBuilder) WithOutputDirectory(outputDirectory string) *TestProcessorBuilder {
	b.outputDirectory = &outputDirectory
	return b
}

func (b *TestProcessorBuilder) WithSessionToken(sessionToken string) *TestProcessorBuilder {
	b.sessionToken = &sessionToken
	return b
}

func (b *TestProcessorBuilder) Build(t *testing.T, mockServerURL string) *MetadataPostProcessor {
	var integrationID string
	if b.integrationID == nil {
		integrationID = uuid.NewString()
	} else {
		integrationID = *b.integrationID
	}

	var inputDirectory string
	if b.inputDirectory == nil {
		inputDirectory = t.TempDir()
	} else {
		inputDirectory = *b.inputDirectory
	}

	var outputDirectory string
	if b.outputDirectory == nil {
		outputDirectory = t.TempDir()
	} else {
		outputDirectory = *b.outputDirectory
	}

	var sessionToken string
	if b.sessionToken == nil {
		sessionToken = uuid.NewString()
	} else {
		sessionToken = *b.sessionToken
	}

	processor, err := NewMetadataPostProcessor(integrationID, inputDirectory, outputDirectory, sessionToken, mockServerURL, mockServerURL)
	require.NoError(t, err)
	return processor
}
