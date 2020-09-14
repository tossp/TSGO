package st

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tossp/tsgo/pkg/errors"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/utils"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
)

func StXorm(c echo.Context, obj interface{}, omit ...string) (data map[string]interface{}, err error) {
	if err = mustPtrStruct(obj); err != nil {
		return
	}
	orgPi := c.QueryParam("pi")          //分页数
	orgPs := c.QueryParam("ps")          //每页数量
	orgSort := c.QueryParam("sort")      //排序
	filter := make(map[string]string, 0) //筛选
	for _, v := range getFieldName(obj) {
		tmp := c.QueryParam(v)
		if tmp == "" {
			continue
		}
		if tmp == "descend" || tmp == "ascend" {
			// TODO 单项排序的处理逻辑
			orgSort = fmt.Sprintf("%s-%s.%s", orgSort, v, tmp)
			continue
		}
		filter[v] = c.QueryParam(v)
	}
	if exfilter, ok := c.Get(stFilterKey).([]string); ok {
		for _, v := range exfilter {
			tmp := strings.Split(v, ":")
			if len(tmp) < 2 || tmp[0] == "" {
				continue
			}
			filter[tmp[0]] = strings.Join(tmp[1:], ":")
		}
	}

	pi, _ := strconv.Atoi(orgPi)
	if pi < 1 {
		pi = 1
	}

	ps, _ := strconv.Atoi(orgPs)
	if ps < 50 {
		ps = 50
	}
	var m *xorm.Session
	//if len(omit) == 0 {
	//	m = db.X()
	//} else {
	//	m = db.X().Omit(omit...)
	//}

	//multiSort := strings.Split(orgSort, "-")
	m = makeXorm(m, orgSort, filter)
	data = make(map[string]interface{})
	objs, err := makePtrSlice(obj, 5)
	if err != nil {
		err = errors.NewMessageError(err, 0, "反射数据实体错误")
		return
	}
	if err = m.Limit(ps, ps*(pi-1)).Find(objs); err != nil {
		sql, age := m.LastSQL()
		log.Debugf("ST查询语句 Find %s %#v", sql, age)
		err = errors.NewMessageError(err, 0, "获取列表错误")
		return
	}

	m = makeXorm(m, orgSort, filter)
	total, err := m.Count(obj)
	if err != nil {
		err = errors.NewMessageError(err, 0, "获取数量错误")
		return
	}
	sql, age := m.LastSQL()
	log.Debugf("ST查询语句 Count %s %#v", sql, age)

	data["total"] = total
	data["list"] = objs
	return
}

func makeXorm(sess *xorm.Session, orgSort string, filter map[string]string) *xorm.Session {
	for _, v := range strings.Split(orgSort, "-") {
		tmp := strings.Split(v, ".")
		if len(tmp) < 2 || tmp[0] == "" {
			continue
		}
		key := utils.GonicCasedName(strings.Join(tmp[:len(tmp)-1], "."))
		value := tmp[len(tmp)-1]
		if value == "ascend" {
			sess = sess.Asc(key)
		}
		if value == "descend" {
			sess = sess.Desc(key)
		}
	}

	for k, v := range filter {
		tmp := strings.Split(v, ",")
		if len(tmp) < 1 || tmp[0] == "" {
			continue
		}
		for _, v := range tmp {
			sess = sess.And(fmt.Sprintf("%s = ?", utils.GonicCasedName(k)), v)
		}
	}
	return sess
}
