package processor

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/clienttest"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadataPostProcessor_ProcessModelChanges(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"create model and record":             createModel,
		"create record; model already exists": createRecordModelExists,
		"update record":                       updateRecord,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func createModel(t *testing.T) {
	datasetID := newDatasetID()
	modelID := uuid.NewString()

	modelCreate := clienttest.NewModelCreate()

	propertiesCreate := clientmodels.PropertiesCreate{
		clienttest.NewPropertyCreateSimple(t, datatypes.StringType),
		clienttest.NewPropertyCreateArray(t, datatypes.DoubleType),
	}

	recordCreateValues := clienttest.NewRecordValues(clienttest.NewRecordValueSimple(t, datatypes.StringType))

	expectedCreateCall := NewExpectedModelCreateCall(datasetID, modelID, modelCreate)
	expectedPropsCreateCall := NewExpectedPropertiesCreateCall(datasetID, modelID, propertiesCreate)
	expectedRecordCreateCall := NewExpectedRecordCreateCall(datasetID, modelID, recordCreateValues)

	mockServer := NewMockModelServer(t,
		expectedCreateCall,
		expectedPropsCreateCall,
		expectedRecordCreateCall)
	defer mockServer.Close()

	processor := NewTestProcessorBuilder().WithInputDirectory("testdata/input_no_model").Build(t, mockServer.URL())

	createdRecordExternalID := clienttest.NewExternalInstanceID()
	require.NoError(t, processor.ProcessModelChanges(datasetID,
		clientmodels.ModelChanges{
			Create: &clientmodels.ModelPropsCreate{
				Model:      modelCreate,
				Properties: propertiesCreate,
			},
			Records: clientmodels.RecordChanges{
				Create: []clientmodels.RecordCreate{{
					ExternalID:   createdRecordExternalID,
					RecordValues: recordCreateValues,
				}},
			},
		}))
	mockServer.AssertAllCalledExactlyOnce(t)

	idStore := processor.IDStore
	assert.Contains(t, idStore.ModelByName, modelCreate.Name)
	pennsieveModelID := clientmodels.PennsieveSchemaID(modelID)
	assert.Equal(t, pennsieveModelID, idStore.ModelByName[modelCreate.Name])

	recordKey := RecordIDKey{
		ModelID:    pennsieveModelID,
		ExternalID: createdRecordExternalID,
	}
	assert.Contains(t, idStore.RecordIDbyKey, recordKey)
	assert.Equal(t, clientmodels.PennsieveInstanceID(expectedRecordCreateCall.APIResponse.ID), idStore.RecordIDbyKey[recordKey])
}

func createRecordModelExists(t *testing.T) {
	datasetID := newDatasetID()
	modelID := uuid.NewString()
	externalRecordID := clienttest.NewExternalInstanceID()

	recordCreateValues := clienttest.NewRecordValues(clienttest.NewRecordValueSimple(t, datatypes.DoubleType))
	expectedRecordCreateCall := NewExpectedRecordCreateCall(datasetID, modelID, recordCreateValues)

	mockServer := NewMockModelServer(t, expectedRecordCreateCall)
	defer mockServer.Close()

	processor := NewTestProcessorBuilder().WithInputDirectory("testdata/input_no_records").Build(t, mockServer.URL())

	require.NoError(t, processor.ProcessModelChanges(datasetID, clientmodels.ModelChanges{
		ID: modelID,
		Records: clientmodels.RecordChanges{
			Create: []clientmodels.RecordCreate{
				{
					ExternalID:   externalRecordID,
					RecordValues: recordCreateValues,
				},
			},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)

	idStore := processor.IDStore

	recordKey := RecordIDKey{
		ModelID:    clientmodels.PennsieveSchemaID(modelID),
		ExternalID: externalRecordID,
	}
	assert.Contains(t, idStore.RecordIDbyKey, recordKey)
	assert.Equal(t,
		clientmodels.PennsieveInstanceID(expectedRecordCreateCall.APIResponse.ID),
		idStore.RecordIDbyKey[recordKey])
}

func updateRecord(t *testing.T) {
	datasetID := newDatasetID()
	modelID := uuid.NewString()

	recordID := clienttest.NewPennsieveInstanceID()

	recordUpdateValues := clienttest.NewRecordValues(
		clienttest.NewRecordValueSimple(t, datatypes.DoubleType),
		clienttest.NewRecordValueSimple(t, datatypes.StringType),
	)
	expectedRecordUpdateCall := NewExpectedRecordUpdateCall(datasetID, modelID, recordID, recordUpdateValues)

	mockServer := NewMockModelServer(t, expectedRecordUpdateCall)
	defer mockServer.Close()

	processor := NewTestProcessorBuilder().WithInputDirectory("testdata/input").Build(t, mockServer.URL())

	require.NoError(t, processor.ProcessModelChanges(datasetID, clientmodels.ModelChanges{
		ID: modelID,
		Records: clientmodels.RecordChanges{
			Update: []clientmodels.RecordUpdate{
				{
					PennsieveID:  recordID,
					RecordValues: recordUpdateValues,
				},
			},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)

}
