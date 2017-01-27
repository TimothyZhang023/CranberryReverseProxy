package crp

import (
	"time"
	"net"
)

type Options struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string

	MaxPipe int

	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func() (net.Conn, error)

	// The maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int

	ReaderBufferSize int
	WriterBufferSize int

	// Sets the deadline for establishing new connections. If reached,
	// dial will fail with a timeout.
	DialTimeout time.Duration
	// Sets the deadline for socket reads. If reached, commands will
	// fail with a timeout instead of blocking.
	ReadTimeout time.Duration
	// Sets the deadline for socket writes. If reached, commands will
	// fail with a timeout instead of blocking.
	WriteTimeout time.Duration

	// Specifies amount of time after which client closes idle
	// connections. Should be less than server's timeout.
	// Default is to not close idle connections.
	IdleTimeout time.Duration
}

func (opt *Options) getNetwork() string {
	if opt.Network == "" {
		return "tcp"
	}
	return opt.Network
}

func (opt *Options) getDialer() func() (net.Conn, error) {
	if opt.Dialer == nil {
		opt.Dialer = func() (net.Conn, error) {
			conn, err := net.DialTimeout(opt.getNetwork(), opt.Addr, opt.getDialTimeout())
			return conn, err
		}
	}
	return opt.Dialer
}

func (opt *Options) getDialTimeout() time.Duration {
	if opt.DialTimeout == 0 {
		return 2 * time.Second
	}
	return opt.DialTimeout
}

func (opt *Options) getIdleTimeout() time.Duration {
	return opt.IdleTimeout
}
