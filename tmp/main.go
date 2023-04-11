package main

import (
	"github.com/bamcop/kit/preset"
	"golang.org/x/exp/slog"
)

func main() {
	preset.SetDefaultSlog("log/stdout.log")

	slog.Info("A")
	slog.Warn("B")
	slog.Error("C")
}
