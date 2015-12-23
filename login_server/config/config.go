//@auther foolbread
//@time 2015-11-12
//@file aqua/login_server/config/config.go
package config

import (
	fbconfig "fbcommon/config"
	"fbcommon/golog"
)

func InitConfig() {
	golog.Info("initing login config......")
	g_config = new(loginServerConfig)
	loadConfig()
}

func loadConfig() {
	con, err := fbconfig.LoadConfigByFile("conf.ini")
	if err != nil {
		golog.Critical(err)
	}

	str := con.MustString("server", "listen_addr", "")
	golog.Info("listen_addr:", str)
	g_config.setListenAddr(str)

	strs := con.MustStringSlice("redis", "redis_info", nil)
	golog.Info("redis_info:", strs)
	g_config.setDBInfos(strs)
}

func GetConfig() *loginServerConfig {
	return g_config
}

var g_config *loginServerConfig

type loginServerConfig struct {
	addr    string
	db_info []string
}

func (c *loginServerConfig) setListenAddr(a string) {
	c.addr = a
}

func (c *loginServerConfig) GetListenAddr() string {
	return c.addr
}

func (c *loginServerConfig) setDBInfos(a []string) {
	c.db_info = a
}

func (c *loginServerConfig) GetDBInfos() []string {
	return c.db_info
}
