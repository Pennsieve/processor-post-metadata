package processor_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/clienttest"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock/expectedcalls"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"github.com/pennsieve/processor-post-metadata/service/processor"
	"github.com/pennsieve/processor-post-metadata/service/processor/internal/processortest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadataPostProcessor_ProcessProxyInstanceDeletes(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"no deletes, schema does not exist": noDeletesProxySchemaDoesNotExist,
		"no deletes, schema exists":         noDeletesProxySchemaExists,
		"deletes":                           proxyDeletes,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func noDeletesProxySchemaDoesNotExist(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	modelName := uuid.NewString()
	modelID := clienttest.NewPennsieveSchemaID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(modelName, modelID).
		Build()

	targetExternalID := clienttest.NewExternalInstanceID()

	mockServer := mock.NewModelService(t)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessProxyInstanceDeletes(datasetID, clientmodels.ProxyChanges{
		CreateProxyRelationshipSchema: true,
		RecordChanges: []clientmodels.ProxyRecordChanges{
			{
				ModelName:        modelName,
				RecordExternalID: targetExternalID,
				NodeIDCreates:    []string{NewPackageNodeID(), NewPackageNodeID()},
			},
		},
	}))
}

func noDeletesProxySchemaExists(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	modelName := uuid.NewString()
	modelID := clienttest.NewPennsieveSchemaID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(modelName, modelID).
		Build()

	targetExternalID := clienttest.NewExternalInstanceID()

	mockServer := mock.NewModelService(t)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessProxyInstanceDeletes(datasetID, clientmodels.ProxyChanges{
		CreateProxyRelationshipSchema: false,
		RecordChanges: []clientmodels.ProxyRecordChanges{
			{
				ModelName:        modelName,
				RecordExternalID: targetExternalID,
				NodeIDCreates:    []string{NewPackageNodeID(), NewPackageNodeID()},
			},
		},
	}))
}

func proxyDeletes(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	modelName := uuid.NewString()
	modelID := clienttest.NewPennsieveSchemaID()

	model2Name := uuid.NewString()
	model2ID := clienttest.NewPennsieveSchemaID()

	targetExternalID := clienttest.NewExternalInstanceID()
	targetRecordID := clienttest.NewPennsieveInstanceID()

	target2ExternalID := clienttest.NewExternalInstanceID()
	target2RecordID := clienttest.NewPennsieveInstanceID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(modelName, modelID).
		WithModel(model2Name, model2ID).
		WithRecord(modelID, targetExternalID, targetRecordID).
		WithRecord(model2ID, target2ExternalID, target2RecordID).
		Build()

	targetProxyInstanceID := clienttest.NewPennsieveInstanceID()
	targetProxyInstance2IO := clienttest.NewPennsieveInstanceID()

	target2ProxyInstanceID := clienttest.NewPennsieveInstanceID()

	expectedCall := expectedcalls.DeleteProxyInstances(datasetID, models.DeleteProxyInstancesBody{
		SourceRecordID:   targetRecordID,
		ProxyInstanceIDs: []clientmodels.PennsieveInstanceID{targetProxyInstanceID, targetProxyInstance2IO},
	}, models.DeleteProxyInstancesBody{
		SourceRecordID:   target2RecordID,
		ProxyInstanceIDs: []clientmodels.PennsieveInstanceID{target2ProxyInstanceID},
	})

	mockServer := mock.NewModelService(t, expectedCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessProxyInstanceDeletes(datasetID, clientmodels.ProxyChanges{
		CreateProxyRelationshipSchema: false,
		RecordChanges: []clientmodels.ProxyRecordChanges{
			{
				ModelName:         modelName,
				RecordExternalID:  targetExternalID,
				NodeIDCreates:     []string{NewPackageNodeID(), NewPackageNodeID()},
				InstanceIDDeletes: []clientmodels.PennsieveInstanceID{targetProxyInstanceID, targetProxyInstance2IO},
			},
			{
				ModelName:         model2Name,
				RecordExternalID:  target2ExternalID,
				InstanceIDDeletes: []clientmodels.PennsieveInstanceID{target2ProxyInstanceID},
			},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}

func TestMetadataPostProcessor_ProcessProxyChanges(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"create schema and instances":     createProxySchemaAndInstances,
		"schema exists; create instances": proxySchemaExistsCreateInstances,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func createProxySchemaAndInstances(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	modelName := uuid.NewString()
	modelID := clienttest.NewPennsieveSchemaID()
	targetExternalID := clienttest.NewExternalInstanceID()
	targetPennsieveID := clienttest.NewPennsieveInstanceID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(modelName, modelID).
		WithRecord(modelID, targetExternalID, targetPennsieveID).
		Build()

	nodeID1 := NewPackageNodeID()
	nodeID2 := NewPackageNodeID()

	expectedCreateCals := expectedcalls.CreateProxyInstance(datasetID,
		models.NewCreateProxyInstanceBody(targetPennsieveID, nodeID1),
		models.NewCreateProxyInstanceBody(targetPennsieveID, nodeID2))
	mockServer := mock.NewModelService(t,
		expectedcalls.CreateProxyRelationshipSchema(datasetID),
		expectedCreateCals)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessProxyChanges(datasetID, &clientmodels.ProxyChanges{
		CreateProxyRelationshipSchema: true,
		RecordChanges: []clientmodels.ProxyRecordChanges{
			{
				ModelName:        modelName,
				RecordExternalID: targetExternalID,
				NodeIDCreates:    []string{nodeID1, nodeID2},
			},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}

func proxySchemaExistsCreateInstances(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	modelName := uuid.NewString()
	modelID := clienttest.NewPennsieveSchemaID()

	model2Name := uuid.NewString()
	model2ID := clienttest.NewPennsieveSchemaID()

	targetExternalID := clienttest.NewExternalInstanceID()
	targetRecordID := clienttest.NewPennsieveInstanceID()

	targetNodeID1 := NewPackageNodeID()
	targetNodeID2 := NewPackageNodeID()

	target2ExternalID := clienttest.NewExternalInstanceID()
	target2RecordID := clienttest.NewPennsieveInstanceID()

	target2NodeID := NewPackageNodeID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(modelName, modelID).
		WithModel(model2Name, model2ID).
		WithRecord(modelID, targetExternalID, targetRecordID).
		WithRecord(model2ID, target2ExternalID, target2RecordID).
		Build()

	mockServer := mock.NewModelService(t,
		expectedcalls.CreateProxyInstance(datasetID,
			models.NewCreateProxyInstanceBody(targetRecordID, targetNodeID1),
			models.NewCreateProxyInstanceBody(targetRecordID, targetNodeID2),
			models.NewCreateProxyInstanceBody(target2RecordID, target2NodeID)),
	)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessProxyChanges(datasetID, &clientmodels.ProxyChanges{
		CreateProxyRelationshipSchema: false,
		RecordChanges: []clientmodels.ProxyRecordChanges{
			{
				ModelName:         modelName,
				RecordExternalID:  targetExternalID,
				NodeIDCreates:     []string{targetNodeID1, targetNodeID2},
				InstanceIDDeletes: []clientmodels.PennsieveInstanceID{clienttest.NewPennsieveInstanceID(), clienttest.NewPennsieveInstanceID()},
			},
			{
				ModelName:         model2Name,
				RecordExternalID:  target2ExternalID,
				NodeIDCreates:     []string{target2NodeID},
				InstanceIDDeletes: []clientmodels.PennsieveInstanceID{clienttest.NewPennsieveInstanceID()},
			},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}

func NewPackageNodeID() string {
	return fmt.Sprintf("N:collection:%s", uuid.NewString())
}
