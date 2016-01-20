//@auther foolbread
//@time 2015-11-12
//@file aqua/login_server/config/config.go
package config

import (
	"flag"

	fbconfig "github.com/foolbread/fbcommon/config"

	"github.com/foolbread/fbcommon/golog"
)

func InitConfig() {
	golog.Info("initing login config......")
	flag.StringVar(&config_path, "f", "conf.ini", "config file path")
	flag.Parse()

	g_config = new(loginServerConfig)
	loadConfig()
}

func loadConfig() {
	con, err := fbconfig.LoadConfigByFile(config_path)
	if err != nil {
		golog.Critical(err)
	}

	str := con.MustString("server", "outer_addr", "")
	golog.Info("outer_addr:", str)
	g_config.setOuterAddr(str)

	str = con.MustString("server", "inner_addr", "")
	golog.Info("inner_addr:", str)
	g_config.setInnerAddr(str)

	strs := con.MustStringSlice("session_storage", "redis_info", nil)
	golog.Info("redis_info:", strs)
	g_config.setSessionDBInfos(strs)
}

func GetConfig() *loginServerConfig {
	return g_config
}

var g_config *loginServerConfig
var config_path string

type loginServerConfig struct {
	outer_addr string
	inner_addr string
	session_db []string
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

func (c *loginServerConfig) setSessionDBInfos(a []string) {
	c.session_db = a
}

func (c *loginServerConfig) GetSessionDBInfos() []string {
	return c.session_db
}
