package main

import (
	"flag"
	"github.com/zts1993/crp"
)

var (
	Cfg     = flag.String("config", "client.conf", "proxy config file")
 )



func main() {
	flag.Parse()

	config :=  crp.NewCRPConfig(*Cfg)
	s := crp.NewCRPClient(config)
	s.Init()
	s.Wg.Wrap(s.MonitorQuit)
	s.Wg.Wrap(s.ClientLink)
	//s.Wg.Wrap(s.TargetLinkRun)

	s.WaitingForQuit()


}
