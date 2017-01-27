package crp

import (
	"fmt"
	"github.com/zts1993/crp/log"
	"github.com/zts1993/crp/util"
	"net"
	"os"
	"runtime"
	"strings"
)

type CRPServer struct {
	LinkListen  net.Listener
	Listen      net.Listener
	LinkSession *Session

	Conf *CRPConfig
	Quit chan bool //quit
	Wg   util.WaitGroupWrapper
}

func NewCRPServer(config *CRPConfig) *CRPServer {

	crp := &CRPServer{
		Conf: config,
	}

	return crp
}

func (ps *CRPServer) Init() {

	log.Info("Proxy Server Init ....")

	listenOn := fmt.Sprintf("%s:%d", ps.Conf.Addr, ps.Conf.Port)
	l, err := net.Listen("tcp4", listenOn)
	if err != nil {
		log.Fatalf("Proxy Server Listen on %s failed ", listenOn)
	}
	log.Info("Proxy Server Listen on", listenOn)
	ps.Listen = l

	linkListenOn := fmt.Sprintf("%s:%d", ps.Conf.LinkAddr, ps.Conf.LinkPort)
	ll, err := net.Listen("tcp4", linkListenOn)
	if err != nil {
		log.Fatalf("Proxy Server LinkListen on %s failed ", listenOn)
	}
	log.Info("Proxy Server Listen on", listenOn)
	ps.LinkListen = ll

}

func (ps *CRPServer) MonitorQuit() {

	util.RegisterSignalAndWait()
	ps.Quit <- true
}

func (ps *CRPServer) LinkRun() {
	log.Info("Proxy Server Run ....")

	ch := make(chan net.Conn, 4096)
	defer close(ch)

	go func() {
		for c := range ch {

			if ps.LinkSession != nil {
				ps.LinkSession.Close()
			}
			ps.LinkSession = NewSession(ps, c)

		}
	}()

	for {
		conn, err := ps.LinkListen.Accept()
		if err != nil {
			log.Info("Got error when Accept network connect ", err)
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Warningf("NOTICE: temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}

			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Warningf("ERROR: listener.Accept() - %s", err)
			}
			break
		}

		ch <- conn
	}
}

func (ps *CRPServer) Run() {

	log.Info("Proxy Server Run ....")

	ch := make(chan net.Conn, 4096)
	defer close(ch)

	go func() {
		for c := range ch {

			//cop := Options{
			//	Network: "tcp",
			//	Addr:    "127.0.0.1:6379",
			//}
			//
			//upstream, err := NewConnection(&cop)
			//if err != nil {
			//	log.Error(err)
			//	return nil
			//}

			if ps.LinkSession != nil {
				s := NewSession(ps, c)
				go s.readLoop(ps.LinkSession.Wt)
				go s.writeLoop(ps.LinkSession.Rd)

			} else {
				log.Error("LinkSession is nil")
				c.Close()
			}

		}
	}()

	for {
		conn, err := ps.Listen.Accept()
		if err != nil {
			log.Info("Got error when Accept network connect ", err)
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Warningf("NOTICE: temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}

			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Warningf("ERROR: listener.Accept() - %s", err)
			}
			break
		}

		ch <- conn
	}
}
func (ps *CRPServer) Close() {

	var err error

	log.Info("Proxy Server Close Socket LinkListener ")
	err = ps.LinkListen.Close()
	if err != nil {
		log.Error("Close Listener err ", err)
	}

	log.Info("Proxy Server Close Socket Listener ")
	err = ps.Listen.Close()
	if err != nil {
		log.Error("Close Listener err ", err)
	}

	log.Warning("Proxy Server Closing ....")

	ps.Wg.Wait()

}

func (ps *CRPServer) WaitingForQuit() {

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
