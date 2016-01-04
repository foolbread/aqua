//@auther foolbread
//@time 2015-12-01
//@file aqua/singlechat_server/config/config.go
package config

import (
	fbconfig "fbcommon/config"
	"fbcommon/golog"
	"flag"
)

func init() {
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
}

func GetConfig() *singlechatServerConfig {
	return g_config
}

var g_config *singlechatServerConfig
var config_path string

type singlechatServerConfig struct {
	connect_server []string
}

func (s *singlechatServerConfig) setConnetServer(a []string) {
	s.connect_server = a
}

func (s *singlechatServerConfig) GetConnetServer() []string {
	return s.connect_server
}
