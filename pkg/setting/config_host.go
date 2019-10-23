package setting

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/tossp/tsgo/pkg/utils/crypto"
)

type Server struct {
	PubKey     string           `description:"公玥"`
	PassPubKey string           `description:"校验公玥"`
	Addr       string           `description:"远程服务器地址"`
	LastUpdate time.Time        `description:"上次更新时间"`
	Host       map[string]*Host `description:"虚拟主机"`
}

type Host struct {
	Domain     []string  `description:"绑定域名"`
	LastUpdate time.Time `description:"上次更新时间"`
}

func (s Server) SignVer(pt, sign []byte) (ok bool) {
	if s.PubKey == "" {
		return
	}
	fmt.Println("SignVer", s.PubKey)
	return crypto.SignVer(crypto.ToECDSAPub(crypto.HexDecode(s.PubKey)), pt, sign)
}

func GetServers(uid string) (data *Server) {
	data = new(Server)
	m := viper.GetStringMapString("server")
	s := m[uid]
	if len(s) > 0 {
		_ = json.Unmarshal([]byte(s), data)
	}
	return
}
func SetServers(uid string, data *Server) (err error) {
	data.PassPubKey = ""
	m := viper.GetStringMapString("server")
	s, _ := json.Marshal(data)
	m[uid] = string(s)
	viper.Set("server", m)
	err = viper.WriteConfig()
	return
}
func GetAllServers() (data []*Server) {
	data = make([]*Server, 0)
	m := viper.GetStringMapString("server")
	for _, s := range m {
		if len(s) > 0 {
			tmp := new(Server)
			if err := json.Unmarshal([]byte(s), tmp); err != nil {
				data = append(data, tmp)
			}
		}
	}
	return
}
