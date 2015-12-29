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

	str := con.MustString("server", "outer_addr", "")
	golog.Info("outer_addr:", str)
	g_config.setOuterAddr(str)

	str = con.MustString("server", "inner_addr", "")
	golog.Info("inner_addr:", str)
	g_config.setInnerAddr(str)

	strs := con.MustStringSlice("redis", "redis_info", nil)
	golog.Info("redis_info:", strs)
	g_config.setDBInfos(strs)
}

func GetConfig() *loginServerConfig {
	return g_config
}

var g_config *loginServerConfig

type loginServerConfig struct {
	outer_addr string
	inner_addr string
	db_info    []string
}

func (c *loginServerConfig) setOuterAddr(a string) {
	c.outer_addr = a
}

func (c *loginServerConfig) GetOuterAddr() string {
	return c.outer_addr
}

func (c *loginServerConfig) setInnerAddr(a string) {
	c.inner_addr = a
}

func (c *loginServerConfig) GetInnerAddr() string {
	return c.inner_addr
}

func (c *loginServerConfig) setDBInfos(a []string) {
	c.db_info = a
}

func (c *loginServerConfig) GetDBInfos() []string {
	return c.db_info
}
