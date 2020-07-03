package st

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tossp/tsgo/pkg/db"
	"github.com/tossp/tsgo/pkg/errors"
	"github.com/tossp/tsgo/pkg/utils"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func StGorm(c echo.Context, obj interface{}, omit ...string) (data map[string]interface{}, err error) {
	if err = mustPtrStruct(obj); err != nil {
		return
	}
	var m = db.G()
	if ex, ok := c.Get(stGormSetKey).([]*ExSet); ok {
		for _, v := range ex {
			m.Set(v.Name, v.Value)
		}
	}

	if len(omit) > 0 {
		m = m.Omit(omit...)
	}

	orgPi := c.QueryParam("pi")          //分页数
	orgPs := c.QueryParam("ps")          //每页数量
	orgSort := c.QueryParam("sort")      //排序
	filter := make(map[string]string, 0) //筛选
	for _, v := range getFieldName2(obj) {
		tmp := c.QueryParam(v.Name)
		if tmp == "" {
			continue
		}
		if tmp == "ascend" {
			m = m.Order(v.DBName)
			continue
		}
		if tmp == "descend" {
			m = m.Order(fmt.Sprintf("%s desc", v.DBName))
			continue
		}
		filter[v.DBName] = c.QueryParam(v.Name)
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
	if ps < 5 {
		ps = 5
	}
	if ps > 1000 {
		ps = 1000
	}

	//multiSort := strings.Split(orgSort, "-")
	m = makeGorm(m, orgSort, filter)
	data = make(map[string]interface{})
	objs, err := makePtrSlice(obj)
	if err != nil {
		err = errors.NewMessageError(err, 0, "反射数据实体错误")
		return
	}
	if ex, ok := c.Get(stWhereKey).([]*ExWhere); ok {
		for _, v := range ex {
			m = m.Where(v.Query, v.Args...)
		}
	}
	if ex, ok := c.Get(stOrderKey).([]*ExOrder); ok {
		for _, v := range ex {
			m = m.Order(v.Value)
		}
	}
	if ex, ok := c.Get(stOmitKey).([]*ExOmit); ok {
		var tmp []string
		for _, v := range ex {
			tmp = append(tmp, v.Columns...)
		}
		m = m.Omit(tmp...)
	}
	if ex, ok := c.Get(stPreloadKey).([]*ExPreload); ok {
		for _, v := range ex {
			m = m.Preload(v.Column, v.Conditions...)
		}
	}

	var total int64
	if err = m.Model(obj).Count(&total).Error; err != nil {
		err = errors.NewMessageError(err, 0, "获取列表数量错误")
		return
	}
	if err = m.Offset(ps * (pi - 1)).Limit(ps).Find(objs).Error; err != nil {
		err = errors.NewMessageError(err, 0, "获取列表错误")
		return
	}

	data["total"] = total
	data["list"] = objs
	return
}

func makeGorm(sess *gorm.DB, orgSort string, filter map[string]string) *gorm.DB {
	for _, v := range strings.Split(orgSort, "-") {
		tmp := strings.Split(v, ".")
		if len(tmp) < 2 || tmp[0] == "" {
			continue
		}
		key := utils.GonicCasedName(strings.Join(tmp[:len(tmp)-1], "."))
		value := tmp[len(tmp)-1]
		if value == "ascend" {
			sess = sess.Order(key + " asc")
		}
		if value == "descend" {
			sess = sess.Order(key + " desc")
		}
	}

	for k, v := range filter {
		tmp := strings.Split(v, ",")
		if len(tmp) < 1 || tmp[0] == "" {
			continue
		}
		for _, v := range tmp {
			sess = sess.Where(fmt.Sprintf("%s = ?", k), v)
		}
	}
	return sess
}
