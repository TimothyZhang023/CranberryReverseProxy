package util

import (
	"syscall"

	"github.com/zts1993/crp/log"
)

func UpdateRLimit() {

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatal("Error Getting Rlimit", err)
	}
	log.Infof("Current Rlimit %+v", rLimit)
	rLimit.Max = 999999
	rLimit.Cur = 999999

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Error("Error Setting Rlimit", err)
	}

	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatal("Error Getting Rlimit", err)
	}
	log.Infof("Now Rlimit %+v", rLimit)

}
