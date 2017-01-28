package crp

import (
	"bufio"
	"github.com/zts1993/crp/log"
	"github.com/zts1993/crp/util"
	"net"
	"sync/atomic"
	"time"
)

type Connection struct {
	netcn        net.Conn
	Rd           *bufio.Reader
	Wt           *bufio.Writer

	ctime        time.Time //created time
	closed       int32     //1 for closed
}

func NewConnection(opt *Options) (cn *Connection, err error) {

	log.Debug("new connection to ", opt.Addr)

	dialer := opt.getDialer()

	netcn, err := dialer()
	if err != nil {
		return nil, err
	}
	cn = &Connection{
		netcn: netcn,
	}

	cn.Rd = bufio.NewReaderSize(cn, opt.ReaderBufferSize)
	cn.Wt = bufio.NewWriterSize(cn, opt.WriterBufferSize)

	cn.ctime = time.Now()
	go cn.writeLoop()
	go cn.readLoop()

	return cn, nil
}

func (cn *Connection) writeLoop() {
	defer util.PanicProcess(cn.writeLoop, nil)
}

func (cn *Connection) readLoop() {
	defer util.PanicProcess(cn.writeLoop, nil)
}

func (cn *Connection) Read(b []byte) (int, error) {
	return cn.netcn.Read(b)
}

func (cn *Connection) Write(b []byte) (int, error) {
	return cn.netcn.Write(b)
}

func (cn *Connection) RemoteAddr() net.Addr {
	return cn.netcn.RemoteAddr()
}

func (cn *Connection) Close() error {

	need_closed := atomic.CompareAndSwapInt32(&cn.closed, 0, 1)

	if need_closed {
		cn.netcn.Close()
	}

	return nil
}
