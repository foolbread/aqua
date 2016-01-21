//@auther foolbread
//@time 2016-01-21
//@file aqua/relation_server/config/config.go
package config

import (
	"flag"

	fbconfig "github.com/foolbread/fbcommon/config"
	"github.com/foolbread/fbcommon/golog"
)

func InitConfig() {
	golog.Info("initing relation server config......")

	flag.StringVar(&config_path, "f", "conf.ini", "config file path")
	flag.Parse()

	g_config = new(relationServerConfig)

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

	strs = c.MustStringSlice("session_storage", "redis_info", nil)
	golog.Info("session_storage:", strs)
	g_config.setSessionDBInfos(strs)

	strs = c.MustStringSlice("relation_storage", "redis_info", nil)
	golog.Info("relation_storage:", strs)
	g_config.setRelationDBInfos(strs)
}

func GetConfig() *relationServerConfig {
	return g_config
}

var g_config *relationServerConfig
var config_path string

type relationServerConfig struct {
	service_type   uint32
	connect_server []string
	session_db     []string
	relation_db    []string
}

func (s *relationServerConfig) setConnetServer(a []string) {
	s.connect_server = a
}

func (s *relationServerConfig) GetConnetServer() []string {
	return s.connect_server
}

func (s *relationServerConfig) setServiceType(t uint32) {
	s.service_type = t
}

func (s *relationServerConfig) GetServiceType() uint32 {
	return s.service_type
}

func (s *relationServerConfig) setSessionDBInfos(a []string) {
	s.session_db = a
}

func (s *relationServerConfig) GetSessionDBInfos() []string {
	return s.session_db
}

func (s *relationServerConfig) setRelationDBInfos(a []string) {
	s.relation_db = a
}

func (s *relationServerConfig) GetRelationDBInfos() []string {
	return s.relation_db
}
