package processor_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/clienttest"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock/expectedcalls"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"github.com/pennsieve/processor-post-metadata/service/processor"
	"github.com/pennsieve/processor-post-metadata/service/processor/internal/processortest"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCurationExportSyncProcessor_Run(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"handle empty changeset without panic":     testEmptyChangeset,
		"create model and record":                  testCreateModelAndRecord,
		"create link between two existing records": testCreateLinkBetweenTwoExistingRecords,
		"link package to existing record":          testLinkPackageToExistingRecord,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func testEmptyChangeset(t *testing.T) {
	integrationID := uuid.NewString()
	datasetID := processortest.NewDatasetID()
	outputDirectory := t.TempDir()

	changeset := clientmodels.Dataset{}
	changesetFilePath := processor.ChangesetFilePath(outputDirectory)
	writeChangeset(t, changeset, changesetFilePath)

	mockServer := mock.NewModelService(t, expectedcalls.GetIntegration(integrationID, datasetID))
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIntegrationID(integrationID).
		WithOutputDirectory(outputDirectory).
		Build(t, mockServer.URL())

	require.NoError(t, testProcessor.Run())

	mockServer.AssertAllCalledExactlyOnce(t)
}

func testCreateModelAndRecord(t *testing.T) {
	integrationID := uuid.NewString()
	datasetID := processortest.NewDatasetID()
	outputDirectory := t.TempDir()

	modelID := clienttest.NewPennsieveSchemaID()

	modelCreate := clienttest.NewModelCreate()

	propertiesCreate := clientmodels.PropertiesCreateParams{
		clienttest.NewPropertyCreateSimple(t, datatypes.StringType),
		clienttest.NewPropertyCreateArray(t, datatypes.DoubleType),
	}

	recordCreateValues := clienttest.NewRecordValues(clienttest.NewRecordValueSimple(t, datatypes.StringType))

	expectedCreateCall := expectedcalls.ModelCreate(datasetID, modelID, modelCreate)
	expectedPropsCreateCall := expectedcalls.PropertiesCreate(datasetID, modelID, propertiesCreate)
	expectedRecordCreateCall := expectedcalls.RecordCreate(datasetID, modelID, recordCreateValues)

	createdRecordExternalID := clienttest.NewExternalInstanceID()

	changeset := clientmodels.Dataset{
		Models: clientmodels.ModelChanges{
			Creates: []clientmodels.ModelCreate{{
				Create: clientmodels.ModelPropsCreate{
					Model:      modelCreate,
					Properties: propertiesCreate,
				},
				Records: []clientmodels.RecordCreate{
					{
						ExternalID:   createdRecordExternalID,
						RecordValues: recordCreateValues,
					},
				},
			}},
		},
	}
	changesetFilePath := processor.ChangesetFilePath(outputDirectory)
	writeChangeset(t, changeset, changesetFilePath)

	mockServer := mock.NewModelService(t,
		expectedcalls.GetIntegration(integrationID, datasetID),
		expectedCreateCall,
		expectedPropsCreateCall,
		expectedRecordCreateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIntegrationID(integrationID).
		WithOutputDirectory(outputDirectory).
		Build(t, mockServer.URL())

	require.NoError(t, testProcessor.Run())

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

func testCreateLinkBetweenTwoExistingRecords(t *testing.T) {
	integrationID := uuid.NewString()
	datasetID := processortest.NewDatasetID()
	outputDirectory := t.TempDir()

	linkSchemaID := clienttest.NewPennsieveSchemaID()
	fromModelName := uuid.NewString()
	fromModelID := clienttest.NewPennsieveSchemaID()
	toModelName := uuid.NewString()

	fromExternalID := clienttest.NewExternalInstanceID()
	fromRecordID := clienttest.NewPennsieveInstanceID()

	toExternalID := clienttest.NewExternalInstanceID()
	toRecordID := clienttest.NewPennsieveInstanceID()

	changeset := clientmodels.Dataset{
		ExistingModelIDMap: map[string]clientmodels.PennsieveSchemaID{
			fromModelName: fromModelID,
			toModelName:   clienttest.NewPennsieveSchemaID(),
		},
		RecordIDMaps: []clientmodels.RecordIDMap{
			{
				ModelName: fromModelName,
				ExternalToPennsieve: map[clientmodels.ExternalInstanceID]clientmodels.PennsieveInstanceID{
					fromExternalID: fromRecordID,
				},
			},
			{
				ModelName: toModelName,
				ExternalToPennsieve: map[clientmodels.ExternalInstanceID]clientmodels.PennsieveInstanceID{
					toExternalID: toRecordID,
				},
			},
		},
		LinkedProperties: []clientmodels.LinkedPropertyChanges{
			{
				FromModelName: fromModelName,
				ToModelName:   toModelName,
				ID:            linkSchemaID,
				Instances: clientmodels.InstanceChanges{
					Create: []clientmodels.InstanceLinkedPropertyCreate{
						{
							FromExternalID: fromExternalID,
							ToExternalID:   toExternalID,
						},
					},
				},
			},
		},
	}
	changesetFilePath := processor.ChangesetFilePath(outputDirectory)
	writeChangeset(t, changeset, changesetFilePath)

	mockServer := mock.NewModelService(t,
		expectedcalls.GetIntegration(integrationID, datasetID),
		expectedcalls.CreateLinkInstance(datasetID, fromModelID, fromRecordID, models.CreateLinkInstanceBody{
			SchemaLinkedPropertyId: linkSchemaID,
			To:                     toRecordID,
		}))
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIntegrationID(integrationID).
		WithOutputDirectory(outputDirectory).
		Build(t, mockServer.URL())

	require.NoError(t, testProcessor.Run())

	mockServer.AssertAllCalledExactlyOnce(t)
}

func testLinkPackageToExistingRecord(t *testing.T) {
	integrationID := uuid.NewString()
	datasetID := processortest.NewDatasetID()
	outputDirectory := t.TempDir()

	modelName := uuid.NewString()
	modelID := clienttest.NewPennsieveSchemaID()

	targetExternalID := clienttest.NewExternalInstanceID()
	targetRecordID := clienttest.NewPennsieveInstanceID()

	packageNodeID := NewPackageNodeID()

	changeset := clientmodels.Dataset{
		ExistingModelIDMap: map[string]clientmodels.PennsieveSchemaID{
			modelName: modelID,
		},
		RecordIDMaps: []clientmodels.RecordIDMap{
			{
				ModelName: modelName,
				ExternalToPennsieve: map[clientmodels.ExternalInstanceID]clientmodels.PennsieveInstanceID{
					targetExternalID: targetRecordID,
				},
			},
		},
		Proxies: &clientmodels.ProxyChanges{
			CreateProxyRelationshipSchema: false,
			RecordChanges: []clientmodels.ProxyRecordChanges{
				{
					ModelName:        modelName,
					RecordExternalID: targetExternalID,
					NodeIDCreates:    []string{packageNodeID},
				},
			},
		},
	}
	changesetFilePath := processor.ChangesetFilePath(outputDirectory)
	writeChangeset(t, changeset, changesetFilePath)

	mockServer := mock.NewModelService(t,
		expectedcalls.GetIntegration(integrationID, datasetID),
		expectedcalls.CreateProxyInstance(datasetID, models.NewCreateProxyInstanceBody(targetRecordID, packageNodeID)),
	)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().
		WithIntegrationID(integrationID).
		WithOutputDirectory(outputDirectory).
		Build(t, mockServer.URL())

	require.NoError(t, testProcessor.Run())

	mockServer.AssertAllCalledExactlyOnce(t)
}

func writeChangeset(t *testing.T, changeset clientmodels.Dataset, filePath string) {
	file, err := os.Create(filePath)
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(file).Encode(changeset))
}
