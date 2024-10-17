package main

import (
	"github.com/pennsieve/processor-post-metadata/service/logging"
	"github.com/pennsieve/processor-post-metadata/service/processor"
	"log/slog"
	"os"
)

var logger = logging.PackageLogger("main")

func main() {
	m, err := processor.FromEnv()
	if err != nil {
		logger.Error("error creating processor", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("created MetadataPostProcessor",
		slog.String("integrationID", m.IntegrationID),
		slog.String("inputDirectory", m.InputDirectory),
		slog.String("outputDirectory", m.OutputDirectory),
		slog.String("apiHost", m.Pennsieve.APIHost),
		slog.String("api2Host", m.Pennsieve.API2Host),
	)

	if err := m.Run(); err != nil {
		logger.Error("error running processor", slog.Any("error", err))
		os.Exit(1)
	}
}
