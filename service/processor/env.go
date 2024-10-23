package processor

import (
	"fmt"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"os"
)

const IntegrationIDKey = "INTEGRATION_ID"
const InputDirectoryKey = "INPUT_DIR"
const OutputDirectoryKey = "OUTPUT_DIR"
const SessionTokenKey = "SESSION_TOKEN"
const PennsieveAPIHostKey = "PENNSIEVE_API_HOST"
const PennsieveAPI2HostKey = "PENNSIEVE_API_HOST2"

func FromEnv() (*MetadataPostProcessor, error) {
	integrationID, err := LookupRequiredEnvVar(IntegrationIDKey)
	if err != nil {
		return nil, err
	}
	inputDirectory, err := LookupRequiredEnvVar(InputDirectoryKey)
	if err != nil {
		return nil, err
	}
	outputDirectory, err := LookupRequiredEnvVar(OutputDirectoryKey)
	if err != nil {
		return nil, err
	}
	sessionToken, err := LookupRequiredEnvVar(SessionTokenKey)
	if err != nil {
		return nil, err
	}
	apiHost, err := LookupRequiredEnvVar(PennsieveAPIHostKey)
	if err != nil {
		return nil, err
	}
	api2Host, err := LookupRequiredEnvVar(PennsieveAPI2HostKey)
	if err != nil {
		return nil, err
	}
	reader, err := metadataclient.NewReader(inputDirectory)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata reader for %s: %w", inputDirectory, err)
	}
	idStore := NewIDStoreBuilder().WithSchema(reader.Schema).Build()
	return NewMetadataPostProcessor(integrationID,
		inputDirectory,
		outputDirectory,
		sessionToken,
		apiHost,
		api2Host,
		idStore,
	)
}

func LookupRequiredEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("no %s set", key)
	}
	return value, nil
}
