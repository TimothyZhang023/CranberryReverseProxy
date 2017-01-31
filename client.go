package crp

import (
	"fmt"
	"github.com/zts1993/crp/log"
	"github.com/zts1993/crp/util"
	"os"
)

type CRPClient struct {
	LinkSession   *Session
	TargetSession *Session

	Conf *CRPConfig
	Quit chan bool //quit
	Wg   util.WaitGroupWrapper
}

func NewCRPClient(config *CRPConfig) *CRPClient {

	crp := &CRPClient{
		Conf: config,
		Quit: make(chan bool, 2),

	}

	return crp
}

func (ps *CRPClient) Init() {
}

func (ps *CRPClient) MonitorQuit() {

	util.RegisterSignalAndWait()
	ps.Quit <- true
}

func (ps *CRPClient) Close() {

	ps.Wg.Wait()

}

func (ps *CRPClient) ClientLink() {

	opt := Options{
		Addr: fmt.Sprintf("%s:%d", ps.Conf.LinkAddr, ps.Conf.LinkPort),
	}

	cn, err := NewConnection(&opt)
	if err != nil {
		log.Fatal(err)
	}

	ps.LinkSession = NewSession(ps.Conf, cn.netcn)

	opt2 := Options{
		Addr: fmt.Sprintf("%s:%d", ps.Conf.TargetAddr, ps.Conf.TargetPort),
	}

	cn2, err := NewConnection(&opt2)
	if err != nil {
		log.Fatal(err)
	}

	ps.TargetSession = NewSession(ps.Conf, cn2.netcn)

	go ps.LinkSession.readLoop(ps.TargetSession.c)
	go ps.LinkSession.writeLoop(ps.TargetSession.c)

	go ps.TargetSession.readLoop(ps.LinkSession.c)
	go ps.TargetSession.writeLoop(ps.LinkSession.c)

}

func (ps *CRPClient) WaitingForQuit() {

	for {
		select {
		case quit := <-ps.Quit:
			log.Info("Receive quit signal", quit)
			if quit {
				goto quit
			}
		}
	}

quit:
	ps.Quit <- true

	ps.Close()

	log.Warning("Proxy server closed")

	os.Exit(0)
}
