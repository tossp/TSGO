package casbin

import (
	"github.com/spf13/viper"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/tossp/tsgo/pkg/casbin/gorm-adapter"
	"github.com/tossp/tsgo/pkg/db"
)

const (
	defAcm = `[request_definition]
r = sub, obj, act, service

[policy_definition]
p = sub, obj, act, service, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = (r.service == p.service || p.service=="*") && ( g(r.sub, p.sub) || p.sub=="*") && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && regexMatch(r.act, p.act)`
)

var (
	adapter  *gormAdapter.Adapter
	enforcer *casbin.SyncedEnforcer
)

func init() {
	viper.SetDefault("accessControl.Enable", true)
}

func Start() (err error) {
	if adapter, err = gormAdapter.NewAdapterByDBUsePrefix(db.G(), db.TableName("")); err != nil {
		return
	}
	m := make(model.Model)
	//log.Info("鉴权模型", setting.GetAccessControlModel())
	if err = m.LoadModelFromText(defAcm); err != nil {
		return
	}
	if enforcer, err = casbin.NewSyncedEnforcer(m, adapter); err != nil {
		return
	}
	enforcer.EnableEnforce(viper.GetBool("accessControl.Enable"))
	enforcer.EnableLog(false)
	enforcer.EnableAutoSave(true)
	enforcer.EnableAutoBuildRoleLinks(true)
	enforcer.StartAutoLoadPolicy(time.Minute * 15)
	//err = enforcer.LoadPolicy()
	return
}

func Adapter() *gormAdapter.Adapter {
	return adapter
}

func E() *casbin.SyncedEnforcer {
	return enforcer
}
