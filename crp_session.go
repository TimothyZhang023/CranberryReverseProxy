package crp

import (
	"bufio"
	//"github.com/zts1993/crp/util"
	"github.com/zts1993/crp/log"
	"net"
	"sync/atomic"
	"time"
 )

type Session struct {
	c        net.Conn // tcp for per client
	r        *bufio.Reader
	w        *bufio.Writer
	ps       *CRPServer
	upstream *Connection

	ctime int64
	mtime int64

	closed int32 //1 for closed
}

func NewSession(ps *CRPServer, conn net.Conn) *Session {
	cop := Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",

	}

	upstream, err := NewConnection(&cop)
	if err != nil {
		log.Error(err)
		return nil
	}

	s := &Session{
		c:        conn,
		r:        bufio.NewReaderSize(conn, ps.Conf.ReaderBufSize),
		w:        bufio.NewWriterSize(conn, ps.Conf.WriterBufSize),
		ps:       ps,
		upstream: upstream,
		ctime:    time.Now().Unix(),
		mtime:    time.Now().Unix(),
	}

	return s
}

func (s *Session) readLoop() error {
	//defer util.PanicProcess(s.readLoop, nil)
	//defer s.Cleanup()
	for {
		var buf = make([]byte, 1024)

		n, err := s.r.Read(buf)

		if err != nil {
			log.Error(err)
			break
		}


		if n <= 0 {
			continue
		}

		//log.Info(string(buf))


		s.upstream.Wt.Write(buf[0:n])
		s.upstream.Wt.Flush()

	}

	return nil
}

func (s *Session) writeLoop() error {
	//defer util.PanicProcess(s.readLoop, nil)
	//defer s.Cleanup()

	for {
		var buf = make([]byte, 1024)

		n, err := s.upstream.Rd.Read(buf)

		if err != nil {
			log.Error(err)
			break
		}

		if n <= 0 {
			continue
		}


		s.w.Write(buf[0:n])
		s.w.Flush()

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
