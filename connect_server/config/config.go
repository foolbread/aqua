//@auther foolbread
//@time 2015-11-12
//@file aqua/connect_server/config/config.go
package config

import (
	fbconfig "fbcommon/config"
	"fbcommon/golog"
)

func InitConfig() {
	golog.Info("initing connect config......")
	g_config = new(connectServerConfig)
	loadConfig()
}

func loadConfig() {
	con, err := fbconfig.LoadConfigByFile("conf.ini")
	if err != nil {
		golog.Critical(err)
	}

	str := con.MustString("server", "inner_addr", "")
	golog.Info("inner_addr:", str)
	g_config.setInnerAddr(str)

	str = con.MustString("server", "outer_addr", "")
	golog.Info("outer_addr:", str)
	g_config.setOuterAddr(str)
}

func GetConfig() *connectServerConfig {
	return g_config
}

var g_config *connectServerConfig

type connectServerConfig struct {
	inner_addr string
	outer_addr string
}

func (c *connectServerConfig) setInnerAddr(a string) {
	c.inner_addr = a
}

func (c *connectServerConfig) setOuterAddr(a string) {
	c.outer_addr = a
}

func (c *connectServerConfig) GetInnerAddr() string {
	return c.inner_addr
}

func (c *connectServerConfig) GetOuterAddr() string {
	return c.outer_addr
}
