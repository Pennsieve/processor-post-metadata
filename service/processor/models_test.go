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
	modelID := clienttest.NewPennsieveSchemaID()

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
	assert.Equal(t, modelID, idStore.ModelByName[modelCreate.Name])

	recordKey := processor.RecordIDKey{
		ModelID:    modelID,
		ExternalID: createdRecordExternalID,
	}
	assert.Contains(t, idStore.RecordIDbyKey, recordKey)
	assert.Equal(t, clientmodels.PennsieveInstanceID(expectedRecordCreateCall.APIResponse.ID), idStore.RecordIDbyKey[recordKey])
}

func createRecordModelExists(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelName := uuid.NewString()
	modelID := clienttest.NewPennsieveSchemaID()
	externalRecordID := clienttest.NewExternalInstanceID()

	recordCreateValues := clienttest.NewRecordValues(clienttest.NewRecordValueSimple(t, datatypes.DoubleType))
	expectedRecordCreateCall := mock.NewExpectedRecordCreateCall(datasetID, modelID, recordCreateValues)

	mockServer := mock.NewModelService(t, expectedRecordCreateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIDStore(processor.NewIDStoreBuilder().WithModel(modelName, modelID).Build()).
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
		ModelID:    modelID,
		ExternalID: externalRecordID,
	}
	assert.Contains(t, idStore.RecordIDbyKey, recordKey)
	assert.Equal(t,
		clientmodels.PennsieveInstanceID(expectedRecordCreateCall.APIResponse.ID),
		idStore.RecordIDbyKey[recordKey])
}

func updateRecord(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelID := clienttest.NewPennsieveSchemaID()

	recordID := clienttest.NewPennsieveInstanceID()

	recordUpdateValues := clienttest.NewRecordValues(
		clienttest.NewRecordValueSimple(t, datatypes.DoubleType),
		clienttest.NewRecordValueSimple(t, datatypes.StringType),
	)
	expectedRecordUpdateCall := mock.NewExpectedRecordUpdateCall(datasetID, modelID, recordID, recordUpdateValues)

	mockServer := mock.NewModelService(t, expectedRecordUpdateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIDStore(processor.NewIDStoreBuilder().WithModel(uuid.NewString(), modelID).Build()).
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

func TestMetadataPostProcessor_ProcessModelChangeRecordDeletes(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"no deletes, model does not exist": noDeletesModelDoesNotExist,
		"no deletes, model exists":         noDeletesModelExists,
		"one delete":                       oneRecordDelete,
		"several deletes":                  severalRecordDeletes,
		"failed deletes":                   failedRecordDeletes,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func noDeletesModelDoesNotExist(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	mockServer := mock.NewModelService(t)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().Build(t, mockServer.URL())

	modelCreate := clienttest.NewModelCreate()

	require.NoError(t, testProcessor.ProcessModelChangeRecordDeletes(datasetID, clientmodels.ModelChanges{
		Create: &clientmodels.ModelPropsCreate{
			Model:      modelCreate,
			Properties: clientmodels.PropertiesCreate{clienttest.NewPropertyCreateSimple(t, datatypes.DoubleType)},
		},
	}))
}

func noDeletesModelExists(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelID := clienttest.NewPennsieveSchemaID()

	mockServer := mock.NewModelService(t)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessModelChangeRecordDeletes(datasetID, clientmodels.ModelChanges{ID: modelID}))
}

func oneRecordDelete(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelID := clienttest.NewPennsieveSchemaID()

	toDelete := []clientmodels.PennsieveInstanceID{clienttest.NewPennsieveInstanceID()}

	expectedDeleteCall := mock.NewExpectedRecordDeleteCall(datasetID, modelID, toDelete)

	mockServer := mock.NewModelService(t, expectedDeleteCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessModelChangeRecordDeletes(datasetID, clientmodels.ModelChanges{
		ID: modelID,
		Records: clientmodels.RecordChanges{
			Delete: toDelete,
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}

func severalRecordDeletes(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelID := clienttest.NewPennsieveSchemaID()

	toDelete := []clientmodels.PennsieveInstanceID{clienttest.NewPennsieveInstanceID(), clienttest.NewPennsieveInstanceID(), clienttest.NewPennsieveInstanceID()}

	expectedDeleteCall := mock.NewExpectedRecordDeleteCall(datasetID, modelID, toDelete)

	mockServer := mock.NewModelService(t, expectedDeleteCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessModelChangeRecordDeletes(datasetID, clientmodels.ModelChanges{
		ID: modelID,
		Records: clientmodels.RecordChanges{
			Delete: toDelete,
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}

func failedRecordDeletes(t *testing.T) {
	datasetID := processortest.NewDatasetID()
	modelID := clienttest.NewPennsieveSchemaID()

	toDelete := []clientmodels.PennsieveInstanceID{clienttest.NewPennsieveInstanceID(), clienttest.NewPennsieveInstanceID(), clienttest.NewPennsieveInstanceID()}

	expectedDeleteCall := mock.NewExpectedRecordDeleteCallFailure(datasetID, modelID, toDelete)

	mockServer := mock.NewModelService(t, expectedDeleteCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().Build(t, mockServer.URL())

	err := testProcessor.ProcessModelChangeRecordDeletes(datasetID, clientmodels.ModelChanges{
		ID: modelID,
		Records: clientmodels.RecordChanges{
			Delete: toDelete,
		},
	})

	require.Error(t, err)
	for _, recordID := range toDelete {
		assert.ErrorContains(t, err, string(recordID))
	}

	mockServer.AssertAllCalledExactlyOnce(t)
}
