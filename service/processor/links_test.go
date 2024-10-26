package processor_test

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/clienttest"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"github.com/pennsieve/processor-post-metadata/service/processor"
	"github.com/pennsieve/processor-post-metadata/service/processor/internal/processortest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadataPostProcessor_ProcessLinkChanges(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"create link schema":                       createLinkSchema,
		"link schema exists; create link instance": createLinkInstance,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func createLinkSchema(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	fromModelName := uuid.NewString()
	fromModelID := clienttest.NewPennsieveSchemaID()
	toModelName := uuid.NewString()
	toModelID := clienttest.NewPennsieveSchemaID()

	instanceCreate := clienttest.NewInstanceLinkedPropertyCreate()
	fromRecordID := clienttest.NewPennsieveInstanceID()
	toRecordID := clienttest.NewPennsieveInstanceID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(fromModelName, fromModelID).
		WithModel(toModelName, toModelID).
		WithRecord(fromModelID, instanceCreate.FromExternalID, fromRecordID).
		WithRecord(toModelID, instanceCreate.ToExternalID, toRecordID).
		Build()

	schemaCreate := clienttest.NewSchemaLinkedPropertyCreate()
	expectedSchemaCreateBody := models.CreateLinkSchemaBody{
		Name:        schemaCreate.Name,
		DisplayName: schemaCreate.DisplayName,
		To:          toModelID,
		Position:    schemaCreate.Position,
	}

	expectedSchemaCreateCall := mock.NewExpectedCreateLinkSchemaCall(datasetID, fromModelID, expectedSchemaCreateBody)

	expectedInstanceCreateBody := models.CreateLinkInstanceBody{
		SchemaLinkedPropertyId: clientmodels.PennsieveSchemaID(expectedSchemaCreateCall.APIResponse.ID),
		To:                     toRecordID,
	}
	expectedInstanceCreateCall := mock.NewExpectedCreateLinkInstanceCall(datasetID, fromModelID, fromRecordID, expectedInstanceCreateBody)

	mockServer := mock.NewModelService(t, expectedSchemaCreateCall, expectedInstanceCreateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessLinkChanges(datasetID, clientmodels.LinkedPropertyChanges{
		FromModelName: fromModelName,
		ToModelName:   toModelName,
		Create:        &schemaCreate,
		Instances: clientmodels.InstanceChanges{
			Create: []clientmodels.InstanceLinkedPropertyCreate{instanceCreate},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}

func createLinkInstance(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	linkSchemaID := clienttest.NewPennsieveSchemaID()

	fromModelName := uuid.NewString()
	fromModelID := clienttest.NewPennsieveSchemaID()
	toModelName := uuid.NewString()
	toModelID := clienttest.NewPennsieveSchemaID()

	instanceCreate := clienttest.NewInstanceLinkedPropertyCreate()
	fromRecordID := clienttest.NewPennsieveInstanceID()
	toRecordID := clienttest.NewPennsieveInstanceID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(fromModelName, fromModelID).
		WithModel(toModelName, toModelID).
		WithRecord(fromModelID, instanceCreate.FromExternalID, fromRecordID).
		WithRecord(toModelID, instanceCreate.ToExternalID, toRecordID).
		Build()

	expectedInstanceCreateBody := models.CreateLinkInstanceBody{
		SchemaLinkedPropertyId: linkSchemaID,
		To:                     toRecordID,
	}
	expectedInstanceCreateCall := mock.NewExpectedCreateLinkInstanceCall(datasetID, fromModelID, fromRecordID, expectedInstanceCreateBody)

	mockServer := mock.NewModelService(t, expectedInstanceCreateCall)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessLinkChanges(datasetID, clientmodels.LinkedPropertyChanges{
		FromModelName: fromModelName,
		ToModelName:   toModelName,
		ID:            linkSchemaID,
		Instances: clientmodels.InstanceChanges{
			Create: []clientmodels.InstanceLinkedPropertyCreate{instanceCreate},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}

func TestMetadataPostProcessor_ProcessLinkChangesInstanceDeletes(t *testing.T) {
	for scenario, testFunc := range map[string]func(t *testing.T){
		"no deletes, schema does not exist": noDeletesLinkSchemaDoesNotExist,
		"no deletes, schema exists":         noDeletesLinkSchemaExists,
		"deletes":                           linkDeletes,
	} {
		t.Run(scenario, func(t *testing.T) {
			testFunc(t)
		})
	}
}

func noDeletesLinkSchemaDoesNotExist(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	fromModelName := uuid.NewString()
	fromModelID := clienttest.NewPennsieveSchemaID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(fromModelName, fromModelID).
		Build()

	schemaCreate := clienttest.NewSchemaLinkedPropertyCreate()

	mockServer := mock.NewModelService(t)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessLinkChangesInstanceDeletes(datasetID, clientmodels.LinkedPropertyChanges{
		FromModelName: fromModelName,
		ToModelName:   uuid.NewString(),
		Create:        &schemaCreate,
	}))

}

func noDeletesLinkSchemaExists(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	fromModelName := uuid.NewString()
	fromModelID := clienttest.NewPennsieveSchemaID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(fromModelName, fromModelID).
		Build()

	linkSchemaID := clienttest.NewPennsieveSchemaID()

	mockServer := mock.NewModelService(t)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessLinkChangesInstanceDeletes(datasetID, clientmodels.LinkedPropertyChanges{
		FromModelName: fromModelName,
		ToModelName:   uuid.NewString(),
		ID:            linkSchemaID,
	}))

}

func linkDeletes(t *testing.T) {
	datasetID := processortest.NewDatasetID()

	fromModelName := uuid.NewString()
	fromModelID := clienttest.NewPennsieveSchemaID()

	linkInstance1ID := clienttest.NewPennsieveInstanceID()
	fromExternal1ID := clienttest.NewExternalInstanceID()
	fromRecord1ID := clienttest.NewPennsieveInstanceID()

	linkInstance2ID := clienttest.NewPennsieveInstanceID()
	fromExternal2ID := clienttest.NewExternalInstanceID()
	fromRecord2ID := clienttest.NewPennsieveInstanceID()

	initialIDStore := processor.NewIDStoreBuilder().
		WithModel(fromModelName, fromModelID).
		WithRecord(fromModelID, fromExternal1ID, fromRecord1ID).
		WithRecord(fromModelID, fromExternal2ID, fromRecord2ID).
		Build()

	linkSchemaID := clienttest.NewPennsieveSchemaID()

	expectedCall1 := mock.NewExpectedDeleteLinkInstanceCall(datasetID, fromModelID, fromRecord1ID, linkInstance1ID)
	expectedCall2 := mock.NewExpectedDeleteLinkInstanceCall(datasetID, fromModelID, fromRecord2ID, linkInstance2ID)

	mockServer := mock.NewModelService(t, expectedCall1, expectedCall2)
	defer mockServer.Close()

	testProcessor := processortest.NewBuilder().WithIDStore(initialIDStore).Build(t, mockServer.URL())

	require.NoError(t, testProcessor.ProcessLinkChangesInstanceDeletes(datasetID, clientmodels.LinkedPropertyChanges{
		FromModelName: fromModelName,
		ToModelName:   uuid.NewString(),
		ID:            linkSchemaID,
		Instances: clientmodels.InstanceChanges{
			Delete: []clientmodels.InstanceLinkedPropertyDelete{
				{
					FromRecordID:             fromRecord1ID,
					InstanceLinkedPropertyID: linkInstance1ID,
				},
				{
					FromRecordID:             fromRecord2ID,
					InstanceLinkedPropertyID: linkInstance2ID,
				},
			},
		},
	}))

	mockServer.AssertAllCalledExactlyOnce(t)
}
