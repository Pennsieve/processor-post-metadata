package processor_test

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/clienttest"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/processor"
	"github.com/pennsieve/processor-post-metadata/service/processor/internal/processortest"
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
	datasetID := processortest.NewDatasetID()
	modelID := uuid.NewString()

	modelCreate := clienttest.NewModelCreate()

	propertiesCreate := clientmodels.PropertiesCreate{
		clienttest.NewPropertyCreateSimple(t, datatypes.StringType),
		clienttest.NewPropertyCreateArray(t, datatypes.DoubleType),
	}

	recordCreateValues := clienttest.NewRecordValues(clienttest.NewRecordValueSimple(t, datatypes.StringType))

	expectedCreateCall := mock.NewExpectedModelCreateCall(datasetID, modelID, modelCreate)
	expectedPropsCreateCall := mock.NewExpectedPropertiesCreateCall(datasetID, modelID, propertiesCreate)
	expectedRecordCreateCall := mock.NewExpectedRecordCreateCall(datasetID, modelID, recordCreateValues)

	mockServer := mock.NewModelService(t,
		expectedCreateCall,
		expectedPropsCreateCall,
		expectedRecordCreateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(processor.NewIDStoreBuilder().Build()).Build(t, mockServer.URL())

	createdRecordExternalID := clienttest.NewExternalInstanceID()
	require.NoError(t, testProcessor.ProcessModelChanges(datasetID,
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

	idStore := testProcessor.IDStore
	assert.Contains(t, idStore.ModelByName, modelCreate.Name)
	pennsieveModelID := clientmodels.PennsieveSchemaID(modelID)
	assert.Equal(t, pennsieveModelID, idStore.ModelByName[modelCreate.Name])

	recordKey := processor.RecordIDKey{
		ModelID:    pennsieveModelID,
		ExternalID: createdRecordExternalID,
	}
	assert.Contains(t, idStore.RecordIDbyKey, recordKey)
	assert.Equal(t, clientmodels.PennsieveInstanceID(expectedRecordCreateCall.APIResponse.ID), idStore.RecordIDbyKey[recordKey])
}

func createRecordModelExists(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelName := uuid.NewString()
	modelID := uuid.NewString()
	externalRecordID := clienttest.NewExternalInstanceID()

	recordCreateValues := clienttest.NewRecordValues(clienttest.NewRecordValueSimple(t, datatypes.DoubleType))
	expectedRecordCreateCall := mock.NewExpectedRecordCreateCall(datasetID, modelID, recordCreateValues)

	mockServer := mock.NewModelService(t, expectedRecordCreateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIDStore(processor.NewIDStoreBuilder().WithModel(modelName, clientmodels.PennsieveSchemaID(modelID)).Build()).
		Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessModelChanges(datasetID, clientmodels.ModelChanges{
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

	idStore := testProcessor.IDStore

	recordKey := processor.RecordIDKey{
		ModelID:    clientmodels.PennsieveSchemaID(modelID),
		ExternalID: externalRecordID,
	}
	assert.Contains(t, idStore.RecordIDbyKey, recordKey)
	assert.Equal(t,
		clientmodels.PennsieveInstanceID(expectedRecordCreateCall.APIResponse.ID),
		idStore.RecordIDbyKey[recordKey])
}

func updateRecord(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelID := uuid.NewString()

	recordID := clienttest.NewPennsieveInstanceID()

	recordUpdateValues := clienttest.NewRecordValues(
		clienttest.NewRecordValueSimple(t, datatypes.DoubleType),
		clienttest.NewRecordValueSimple(t, datatypes.StringType),
	)
	expectedRecordUpdateCall := mock.NewExpectedRecordUpdateCall(datasetID, modelID, recordID, recordUpdateValues)

	mockServer := mock.NewModelService(t, expectedRecordUpdateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIDStore(processor.NewIDStoreBuilder().WithModel(uuid.NewString(), clientmodels.PennsieveSchemaID(modelID)).Build()).
		Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessModelChanges(datasetID, clientmodels.ModelChanges{
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
