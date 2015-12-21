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

	strs := con.MustStringSlice("redis", "redis_addr", nil)
	golog.Info("redis_addr:", strs)
	g_config.setDBAddrs(strs)
}

func GetConfig() *loginServerConfig {
	return g_config
}

var g_config *loginServerConfig

type loginServerConfig struct {
	addr    string
	db_addr []string
}

func (c *loginServerConfig) setListenAddr(a string) {
	c.addr = a
}

func (c *loginServerConfig) GetListenAddr() string {
	return c.addr
}

func (c *loginServerConfig) setDBAddrs(a []string) {
	c.db_addr = a
}

func (c *loginServerConfig) GetDBAddrs() []string {
	return c.db_addr
}
