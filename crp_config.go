package crp

import (
	"github.com/zts1993/crp/config"
	"github.com/zts1993/crp/log"
	"runtime"
)

type CRPConfig struct {
	Port          int // proxy listen port
	TargetPort    int
	TargetString  string
	ReaderBufSize int
	WriterBufSize int
	FileName      string
	Config        config.Configer
}

func NewCRPConfig(filename string) *CRPConfig {
	c, err := config.NewConfig("ini", filename)
	if err != nil {
		log.Fatal("read config file failed ", err)
	}

	loglevel := c.DefaultString("log::loglevel", "info")
	log.SetLevelByString(loglevel)

	highlighting := c.DefaultBool("log::highlighting", false)
	log.SetHighlighting(highlighting)

	logfile := c.DefaultString("log::logfile", "")
	if logfile != "" {
		err := log.SetOutputByName(logfile)
		if err != nil {
			log.Fatal("Set log Output failed ", err)
		}
		log.Info("Set log Output to file ", logfile)
		log.SetRotateByDay()
	}

	cpus := c.DefaultInt("proxy::cpus", runtime.NumCPU())
	runtime.GOMAXPROCS(cpus)
	log.Warningf("Set runtime GOMAXPROCS to %d, total cpu num %d", cpus, runtime.NumCPU())

	crpConfig := &CRPConfig{
		Port:          c.DefaultInt("proxy::port", 10086),
		ReaderBufSize: c.DefaultInt("proxy::readerbuffer", 1048576),
		WriterBufSize: c.DefaultInt("proxy::writerbuffer", 1048576),
		TargetString:  c.DefaultString("target::ip", "127.0.0.1"),
		TargetPort:    c.DefaultInt("target::port", 6379),
	}

	log.Infof("Config %+v", crpConfig)

	return crpConfig
}
