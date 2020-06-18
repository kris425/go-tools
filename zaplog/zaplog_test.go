package zaplog

import (
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(SetAppName("hello"),
		SetDevelopment(true),
		SetLevel("info"),
		SetLogFileDir("E:"),
	)
	logger.Debug("Info msg", zap.Any("user", "missi"))
	logger.Info("Info msg", zap.Any("user", "missi"))
	logger.Error("Error msg", zap.Any("e", "a"))
}
