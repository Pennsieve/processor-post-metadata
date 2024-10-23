package processortest

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/service/processor"
	"github.com/stretchr/testify/require"
	"testing"
)

type Builder struct {
	integrationID   *string
	inputDirectory  *string
	outputDirectory *string
	sessionToken    *string
	idStore         *processor.IDStore
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithIntegrationID(integrationID string) *Builder {
	b.integrationID = &integrationID
	return b
}

func (b *Builder) WithInputDirectory(inputDirectory string) *Builder {
	b.inputDirectory = &inputDirectory
	return b
}

func (b *Builder) WithOutputDirectory(outputDirectory string) *Builder {
	b.outputDirectory = &outputDirectory
	return b
}

func (b *Builder) WithSessionToken(sessionToken string) *Builder {
	b.sessionToken = &sessionToken
	return b
}

func (b *Builder) WithIDStore(idStore *processor.IDStore) *Builder {
	b.idStore = idStore
	return b
}

func (b *Builder) Build(t *testing.T, mockServerURL string) *processor.MetadataPostProcessor {
	var integrationID string
	if b.integrationID == nil {
		integrationID = uuid.NewString()
	} else {
		integrationID = *b.integrationID
	}

	var inputDirectory string
	if b.inputDirectory == nil {
		inputDirectory = t.TempDir()
	} else {
		inputDirectory = *b.inputDirectory
	}

	var outputDirectory string
	if b.outputDirectory == nil {
		outputDirectory = t.TempDir()
	} else {
		outputDirectory = *b.outputDirectory
	}

	var sessionToken string
	if b.sessionToken == nil {
		sessionToken = uuid.NewString()
	} else {
		sessionToken = *b.sessionToken
	}

	testProcessor, err := processor.NewMetadataPostProcessor(integrationID, inputDirectory, outputDirectory, sessionToken, mockServerURL, mockServerURL, b.idStore)
	require.NoError(t, err)
	return testProcessor
}
