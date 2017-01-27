package crp

import (
	"bufio"
	//"github.com/zts1993/crp/util"
	"github.com/zts1993/crp/log"
	"net"
	"sync/atomic"
	"time"
	"io"
)

type Session struct {
	c      net.Conn // tcp for per client
	Rd     *bufio.Reader
	Wt     *bufio.Writer
	ps     *CRPServer

	ctime  int64
	mtime  int64

	closed int32 //1 for closed
}

func NewSession(ps *CRPServer, conn net.Conn) *Session {

	s := &Session{
		c:        conn,
		Rd:        bufio.NewReaderSize(conn, ps.Conf.ReaderBufSize),
		Wt:        bufio.NewWriterSize(conn, ps.Conf.WriterBufSize),
		ps:       ps,
		ctime:    time.Now().Unix(),
		mtime:    time.Now().Unix(),
	}

	return s
}

func (s *Session) readLoop(upstream *Connection) error {
	//defer util.PanicProcess(s.readLoop, nil)
	//defer s.Cleanup()
	for {
		var buf = make([]byte, 1024)

		n, err := s.Rd.Read(buf)

		if err != nil {
			log.Error(err)
			break
		}

		if n <= 0 {
			continue
		}

		upstream.Wt.Write(buf[0:n])
		upstream.Wt.Flush()

	}

	return nil
}

func (s *Session) writeLoop(rd io.Reader) error {
	//defer util.PanicProcess(s.readLoop, nil)
	//defer s.Cleanup()

	for {
		var buf = make([]byte, 1024)

		n, err := rd.Read(buf)

		if err != nil {
			log.Error(err)
			break
		}

		if n <= 0 {
			continue
		}

		s.Wt.Write(buf[0:n])
		s.Wt.Flush()

	}

	return nil
}

func (s *Session) Cleanup() {
	//defer util.PanicProcess(s.Cleanup, nil)

	s.Close()
	//log.Debug(s.String() ,"quit")
}

func (s *Session) Close() {
	need_closed := atomic.CompareAndSwapInt32(&s.closed, 0, 1)

	if need_closed {
		s.c.Close()
		//close
	}
}
