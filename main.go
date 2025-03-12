package main

import (
	"context"
	"errors"
	"log"
	"main/logger"
)

var (
	DEBUG_MODE = true
)

func main() {
	log.SetFlags(0)
	logger.ConfigureLogger(DEBUG_MODE)

	// LOG WITHOUT CONTEXT
	logger.Info("Info log without context")
	logger.Debug("Debug log without context")
	logger.Error("Error log without context")

	// LOG WITH correlationId context
	ctx := logger.GetCorrelationIDLoggerCtx(context.Background(), "cid-123897123")
	appLogger := logger.Ctx(ctx)
	appLogger.Info("Info Log : %s", "some info")
	appLogger.Debug("Debug Log : %s", "some debug")
	appLogger.Info("Error Log : %v", errors.New("some error"))

	// LOG WITH CONTEXT AND additionalInfo
	additionalInfo := map[string]interface{}{
		"some_key": "some_value",
	}
	appLogger.AdditionalInfo(additionalInfo).Info("Info Log with additionalInfo")
	appLogger.AdditionalInfo(additionalInfo).StackTrace(errors.New("stack trace")).Error("Error Log with additionalInfo")

}
