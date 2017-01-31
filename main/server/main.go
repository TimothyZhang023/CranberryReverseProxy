package main

import (
	"flag"
	"github.com/zts1993/crp"
)

var (
	Cfg = flag.String("config", "server.conf", "proxy config file")
)

func main() {

	flag.Parse()

	config := crp.NewCRPConfig(*Cfg)
	s := crp.NewCRPServer(config)
	s.Init()
	s.Wg.Wrap(s.MonitorQuit)
	s.Wg.Wrap(s.Run)
	s.Wg.Wrap(s.ServerLink)

	s.WaitingForQuit()
}
