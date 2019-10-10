package setting

import (
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	PubKey     string           `description:"公玥"`
	Addr       string           `description:"远程服务器地址"`
	LastUpdate time.Time        `description:"上次更新时间"`
	Host       map[string]*Host `description:"虚拟主机"`
}

type Host struct {
	Domain     []string  `description:"绑定域名"`
	LastUpdate time.Time `description:"上次更新时间"`
}

func GetServers() (data map[string]*Server) {
	m := viper.GetStringMap("server")
	data = make(map[string]*Server, len(m))
	for k, v := range m {
		data[k] = v.(*Server)
		if data[k] == nil {
			data[k].Host = make(map[string]*Host)
		}
	}
	return
}
