package clienttest

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/models"
	"math/rand"
)

func NewSchemaLinkedPropertyCreate() models.SchemaLinkedPropertyCreate {
	return models.SchemaLinkedPropertyCreate{
		Name:        uuid.NewString(),
		DisplayName: uuid.NewString(),
		Position:    rand.Intn(9) + 1,
	}
}

func NewInstanceLinkedPropertyCreate() models.InstanceLinkedPropertyCreate {
	return models.InstanceLinkedPropertyCreate{
		FromExternalID: NewExternalInstanceID(),
		ToExternalID:   NewExternalInstanceID(),
	}
}
