package main

import (
	"flag"
	"os"
	"fmt"
	"runtime"
	"github.com/zts1993/crp"
)

var (
	Cfg     = flag.String("config", "setting.ini", "proxy config file")
	Version = flag.Bool("version", false, "show current version")
)

const logo string = `
-------------------------------------------------
  Go Version: %s.
  Go OS/Arch: %s/%s.
-------------------------------------------------
`


func main() {
	fmt.Printf(logo, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	flag.Parse()

	if *Version {
		os.Exit(0)
	}


	config :=  crp.NewCRPConfig(*Cfg)
	s := crp.NewCRPServer(config)
	s.Init()
	s.Wg.Wrap(s.MonitorQuit)
	s.Wg.Wrap(s.Run)
	s.Wg.Wrap(s.LinkRun)

	s.WaitingForQuit()
}
