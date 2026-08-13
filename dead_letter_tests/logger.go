package dead_letter_tests

import (
	"log/slog"
	"os"
)

var (
	logger *slog.Logger
)

func init() {
	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}
