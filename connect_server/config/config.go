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

	str = con.MustString("login", "login_addr", "")
	golog.Info("login_addr", str)
	g_config.setLoginAddr(str)

	val := con.MustInt("server", "id", 0)
	golog.Info("id:", val)
	g_config.setConnectId(val)

	strs := con.MustStringSlice("session_storage", "redis_info", nil)
	golog.Info("redis_info:", strs)
	g_config.setSessionDBInfos(strs)
}

func GetConfig() *connectServerConfig {
	return g_config
}

var g_config *connectServerConfig

type connectServerConfig struct {
	id         int
	inner_addr string
	outer_addr string
	login_addr string
	session_db []string
}

func (c *connectServerConfig) setInnerAddr(a string) {
	c.inner_addr = a
}

func (c *connectServerConfig) setOuterAddr(a string) {
	c.outer_addr = a
}

func (c *connectServerConfig) setLoginAddr(a string) {
	c.login_addr = a
}

func (c *connectServerConfig) setConnectId(id int) {
	c.id = id
}

func (c *connectServerConfig) setSessionDBInfos(a []string) {
	c.session_db = a
}

func (c *connectServerConfig) GetInnerAddr() string {
	return c.inner_addr
}

func (c *connectServerConfig) GetOuterAddr() string {
	return c.outer_addr
}

func (c *connectServerConfig) GetLoginAddr() string {
	return c.login_addr
}

func (c *connectServerConfig) GetConnectId() int {
	return c.id
}

func (c *connectServerConfig) GetSessionDBInfos() []string {
	return c.session_db
}
