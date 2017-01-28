package crp

import (
	"bufio"
	//"github.com/zts1993/crp/util"
	"bytes"
	"github.com/zts1993/crp/log"
 	"net"
	"sync/atomic"
	"time"
)

type Session struct {
	c  net.Conn // tcp for per client
	Rd *bufio.Reader
	Wt *bufio.Writer

	ctime int64
	mtime int64

	closed int32 //1 for closed
}

func NewSession(conf *CRPConfig, conn net.Conn) *Session {

	s := &Session{
		c:     conn,
		Rd:    bufio.NewReaderSize(conn, conf.ReaderBufSize),
		Wt:    bufio.NewWriterSize(conn, conf.WriterBufSize),
		ctime: time.Now().Unix(),
		mtime: time.Now().Unix(),
	}

	return s
}

func (s *Session) readLoop(wt net.Conn) {
	//defer util.PanicProcess(s.readLoop, nil)
	defer s.Cleanup()
	for {
		var buf = make([]byte, 1024)

		n, err := s.c.Read(buf)

		if err != nil {
			log.Error(err)
			return
		}

		if n <= 0 {
			continue
		}

		log.Info("read", s.c.LocalAddr(), s.c.RemoteAddr(), string(buf))

		wt.Write(bytes.Trim(buf, "\x00"))
		//wt.Write(buf[0:n])
		//wt.Flush()

	}

	return
}

func (s *Session) writeLoop(rd net.Conn) {
	//defer util.PanicProcess(s.readLoop, nil)
	defer s.Cleanup()

	for {
		var buf = make([]byte, 1024)

		n, err := rd.Read(buf)

		if err != nil {
			log.Error(err)
			return
		}

		if n <= 0 {
			continue
		}
		log.Info("write", s.c.LocalAddr(), s.c.RemoteAddr(), string(buf))

		s.c.Write(bytes.Trim(buf, "\x00"))
		//s.Wt.Write(buf[0:n])
		//s.Wt.Flush()

	}

	return
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
