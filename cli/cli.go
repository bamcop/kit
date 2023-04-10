package cli

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
