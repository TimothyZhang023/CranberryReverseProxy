package util

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/zts1993/crp/log"
)

func RegisterSignalAndWait() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	quit := <-sc
	log.Warning("Receive signal", quit.String())
}
