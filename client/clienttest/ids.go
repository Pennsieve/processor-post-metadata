package clienttest

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/models"
)

func NewExternalInstanceID() models.ExternalInstanceID {
	return models.ExternalInstanceID(uuid.NewString())
}

func NewPennsieveInstanceID() models.PennsieveInstanceID {
	return models.PennsieveInstanceID(uuid.NewString())
}
