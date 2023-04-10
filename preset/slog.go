package preset

import (
	"os"

	"golang.org/x/exp/slog"
)

func SetDefaultSlog() {
	h := slog.HandlerOptions{
		AddSource: true,
	}.NewTextHandler(os.Stdout)

	logger := slog.New(h)

	slog.SetDefault(logger)
}
