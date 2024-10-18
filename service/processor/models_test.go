package processor

import (
	"encoding/json"
	"github.com/google/uuid"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadataPostProcessor_ProcessModelChanges(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"create model": createModel,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func createModel(t *testing.T) {
	datasetID := newDatasetID()
	modelID := clientmodels.PennsieveInstanceID(uuid.NewString())

	modelCreate := clientmodels.ModelCreate{
		Name:        uuid.NewString(),
		DisplayName: uuid.NewString(),
		Description: uuid.NewString(),
		Locked:      false,
	}

	propertiesCreate := clientmodels.PropertiesCreate{
		clientmodels.PropertyCreate{
			DisplayName:  uuid.NewString(),
			Name:         uuid.NewString(),
			DataType:     json.RawMessage(`"String"`),
			ConceptTitle: true,
			Required:     true,
		},
	}

	expectedCreateCall := NewExpectedModelCreateCall(datasetID, modelID, modelCreate)
	expectedPropsCreateCall := NewExpectedPropertiesCreateCall(datasetID, modelID, propertiesCreate)

	mockServer := NewMockModelServer(t, expectedCreateCall, expectedPropsCreateCall)
	defer mockServer.Close()

	processor := NewTestProcessorBuilder().WithInputDirectory("testdata/input").Build(t, mockServer.URL())

	require.NoError(t, processor.ProcessModelChanges(datasetID,
		clientmodels.ModelChanges{
			Create: &clientmodels.ModelPropsCreate{
				Model:      modelCreate,
				Properties: propertiesCreate,
			},
		}))
	mockServer.AssertAllCalledExactlyOnce(t)

	assert.Contains(t, processor.IDStore.ModelByName, modelCreate.Name)
	assert.Equal(t, modelID, processor.IDStore.ModelByName[modelCreate.Name])
}
