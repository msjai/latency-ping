package logger

import (
	"go.uber.org/zap"
)

// New -.
func New() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() //nolint:errcheck    // очищение всех буферезированных записей журналаа
	sugar := logger.Sugar()

	return sugar
}
