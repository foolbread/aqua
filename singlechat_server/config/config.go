//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/config/config.go
package config

import (
	fbconfig "fbcommon/config"
	"fbcommon/golog"
	"flag"
)

func InitConfig() {
	golog.Info("initing singlechat server config......")
	flag.StringVar(&config_path, "f", "conf.ini", "config file path")
	flag.Parse()

	g_config = new(singlechatServerConfig)
	loadConfig()
}

func loadConfig() {
	c, err := fbconfig.LoadConfigByFile(config_path)
	if err != nil {
		golog.Critical(err)
	}

	strs := c.MustStringSlice("server", "connect_server", nil)
	golog.Info("connect_server:", strs)
	g_config.setConnetServer(strs)

	t := c.MustInt("server", "service_type", 0)
	golog.Info("service_type:", t)
	g_config.setServiceType(uint32(t))
}

func GetConfig() *singlechatServerConfig {
	return g_config
}

var g_config *singlechatServerConfig
var config_path string

type singlechatServerConfig struct {
	service_type   uint32
	connect_server []string
}

func (s *singlechatServerConfig) setConnetServer(a []string) {
	s.connect_server = a
}

func (s *singlechatServerConfig) GetConnetServer() []string {
	return s.connect_server
}

func (s *singlechatServerConfig) setServiceType(t uint32) {
	s.service_type = t
}

func (s *singlechatServerConfig) GetServiceType() uint32 {
	return s.service_type
}
