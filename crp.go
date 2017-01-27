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
	Listen net.Listener //server
	Conf   *CRPConfig
	Quit   chan bool //quit
	Wg     util.WaitGroupWrapper


}

func NewCRPServer(config *CRPConfig) *CRPServer {

	crp := &CRPServer{
		Conf: config,
	}

	return crp
}

func (ps *CRPServer) Init() {

	log.Info("Proxy Server Init ....")

	listenOn := fmt.Sprintf("0.0.0.0:%d", ps.Conf.Port)

	l, err := net.Listen("tcp4", listenOn)

	if err != nil {
		log.Fatalf("Proxy Server Listen on %s failed ", listenOn)
	}
	log.Info("Proxy Server Listen on", listenOn)

	ps.Listen = l

}


func (ps *CRPServer) MonitorQuit() {

	util.RegisterSignalAndWait()
	ps.Quit <- true
}

func (ps *CRPServer) Run() {


	log.Info("Proxy Server Run ....")

	ch := make(chan net.Conn, 4096)
	defer close(ch)

	go func() {
		for c := range ch {
			s := NewSession(ps, c)
 			go s.readLoop()
			go s.writeLoop()
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
