package setting

import (
	"github.com/spf13/viper"
	"github.com/tossp/tsgo/pkg/null"
)

var restrictedAccessList = make([]*null.CIDR, 0, 0)

func IsRestrictedAccess() bool {
	return viper.GetBool("RestrictedAccess.enable")
}
func RestrictedAccessList() (data []*null.CIDR) {
	if len(restrictedAccessList) > 0 {
		data = restrictedAccessList
		return
	}
	s := viper.GetStringSlice("RestrictedAccess.list")
	for _, v := range s {
		ipnet, err := null.ParseCIDR(v)
		if err != nil {
			continue
		}
		restrictedAccessList = append(restrictedAccessList, ipnet)
	}
	data = restrictedAccessList
	return
}
func SetRestrictedAccess(t bool) {
	viper.Set("RestrictedAccess.enable", t)
}
func SetRestrictedAccessList(data []*null.CIDR) {
	s := make([]string, 0, len(data))
	for _, v := range data {
		if !v.IsZero() {
			s = append(s, v.CIDR.String())
		}

	}
	viper.Set("RestrictedAccess.list", s)
	write()
	return
}
