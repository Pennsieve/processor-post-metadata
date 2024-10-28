package clienttest

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/stretchr/testify/require"
	"math/rand"
)

func NewModelCreate() models.ModelCreate {
	return models.ModelCreate{
		Name:        uuid.NewString(),
		DisplayName: uuid.NewString(),
		Description: uuid.NewString(),
		Locked:      false,
	}
}

func NewPropertyCreateSimple(t require.TestingT, dataType datatypes.SimpleType) models.PropertyCreate {
	bytes, err := json.Marshal(dataType)
	require.NoError(t, err)
	return models.PropertyCreate{
		DisplayName:  uuid.NewString(),
		Name:         uuid.NewString(),
		DataType:     bytes,
		ConceptTitle: true,
		Required:     true,
	}
}

func NewArrayDataType(itemType datatypes.SimpleType) datatypes.ArrayDataType {
	return datatypes.ArrayDataType{
		Type:  datatypes.ArrayType,
		Items: datatypes.ItemsType{Type: itemType},
	}
}
func NewPropertyCreateArray(t require.TestingT, itemType datatypes.SimpleType) models.PropertyCreate {
	dataTypeBytes, err := json.Marshal(NewArrayDataType(itemType))
	require.NoError(t, err)
	return models.PropertyCreate{
		DisplayName:  uuid.NewString(),
		Name:         uuid.NewString(),
		DataType:     dataTypeBytes,
		IsMultiValue: true,
	}
}

func NewRecordValueSimple(t require.TestingT, dataType datatypes.SimpleType) models.RecordValue {
	var value any
	switch dataType {
	case datatypes.StringType:
		value = uuid.NewString()
	case datatypes.DoubleType:
		value = rand.ExpFloat64()
	case datatypes.BooleanType:
		value = rand.Intn(2) == 0
	case datatypes.LongType:
		require.FailNow(t, "case not implemented", "cannot use %s", dataType)
	case datatypes.DateType:
		require.FailNow(t, "case not implemented", "cannot use %s", dataType)
	default:
		require.FailNow(t, "unknown datatype", dataType)

	}
	return models.RecordValue{
		Value: value,
		Name:  uuid.NewString(),
	}
}
func NewRecordValues(values ...models.RecordValue) models.RecordValues {
	return models.RecordValues{Values: values}
}
