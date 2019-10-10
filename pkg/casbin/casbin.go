package casbin

import (
	"time"

	gormAdapter "github.com/tossp/tsgo/pkg/casbin/gorm-adapter"
	"github.com/tossp/tsgo/pkg/db"
	"github.com/tossp/tsgo/pkg/setting"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

var adapter *gormAdapter.Adapter
var enforcer *casbin.SyncedEnforcer

func Start() (err error) {
	if adapter, err = gormAdapter.NewAdapterByDBUsePrefix(db.G(), setting.GetAccessControlPrefix()); err != nil {
		return
	}
	m := make(model.Model)
	//log.Info("鉴权模型", setting.GetAccessControlModel())
	if err = m.LoadModelFromText(setting.GetAccessControlModel()); err != nil {
		return
	}
	if enforcer, err = casbin.NewSyncedEnforcer(m, adapter); err != nil {
		return
	}
	enforcer.EnableEnforce(setting.GetAccessControlEnable())
	enforcer.EnableLog(true)
	enforcer.EnableAutoSave(true)
	enforcer.EnableAutoBuildRoleLinks(true)
	enforcer.StartAutoLoadPolicy(time.Minute * 5)
	//err = enforcer.LoadPolicy()
	return
}

func Adapter() *gormAdapter.Adapter {
	return adapter
}

func E() *casbin.SyncedEnforcer {
	return enforcer
}
