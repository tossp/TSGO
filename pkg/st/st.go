package st

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/tossp/tsgo/pkg/errors"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/utils"

	"github.com/labstack/echo/v4"
)

const (
	stFilterKey      = "StFilter"
	stPreloadKey     = "StPreload"
	stRelatedKey     = "StRelated"
	stWhereKey       = "StWhere"
	stOrderKey       = "StOrder"
	stOmitKey        = "StOmit"
	stAssociationKey = "StAssociation"
	stGormSetKey     = "StGormSetKey"
)

//获取结构体中字段的名称
func getFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return []string{}
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

//获取结构体中字段的名称
func getFieldName2(structName interface{}) (fields []*gorm.StructField) {
	scope := &gorm.Scope{Value: structName}
	fields = scope.GetStructFields()[:]
	return
}

func mustPtrStruct(structName interface{}) (err error) {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		err = errors.New(fmt.Sprintf("类型不为指针,%s", t))
		return
	}
	if t.Kind() != reflect.Struct {
		err = errors.New(fmt.Sprintf("类型不为实体,%s", t))
		return
	}
	return
}
func makePtrSlice(structName interface{}) (objs interface{}, err error) {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		err = errors.New(fmt.Sprintf("类型不为指针,%s", t))
		return
	}
	if t.Kind() != reflect.Struct {
		err = errors.New(fmt.Sprintf("类型不为实体,%s", t))
		return
	}
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(structName)), 0, 5)
	tmp := reflect.New(slice.Type())
	tmp.Elem().Set(slice)
	objs = tmp.Interface()
	return
}

func Makefilter(c echo.Context, key, v string) {
	var tmp []string
	if exfilter, ok := c.Get(stFilterKey).([]string); ok {
		tmp = exfilter
	}
	tmp = append(tmp, fmt.Sprintf("%s:%s", utils.GonicCasedName(key), v))
	c.Set(stFilterKey, tmp)
}

type ExPreload struct {
	Column     string
	Conditions []interface{}
}

func MakePreload(c echo.Context, column string, conditions ...interface{}) {
	tmp := make([]*ExPreload, 0)
	if exPreload, ok := c.Get(stPreloadKey).([]*ExPreload); ok {
		tmp = exPreload
	}
	tmp = append(tmp, &ExPreload{
		Column:     column,
		Conditions: conditions,
	})
	c.Set(stPreloadKey, tmp)
}

type ExRelated struct {
	Value       interface{}
	ForeignKeys []string
}

func MakeRelated(c echo.Context, value interface{}, foreignKeys ...string) {
	tmp := make([]*ExRelated, 0)
	if exPreload, ok := c.Get(stRelatedKey).([]*ExRelated); ok {
		tmp = exPreload
	}
	tmp = append(tmp, &ExRelated{
		Value:       value,
		ForeignKeys: foreignKeys,
	})
	c.Set(stRelatedKey, tmp)
}

type ExWhere struct {
	Query interface{}
	Args  []interface{}
}

func MakeWhere(c echo.Context, query interface{}, args ...interface{}) {
	tmp := make([]*ExWhere, 0)
	if ex, ok := c.Get(stWhereKey).([]*ExWhere); ok {
		tmp = ex
	}
	tmp = append(tmp, &ExWhere{
		Query: query,
		Args:  args,
	})
	c.Set(stWhereKey, tmp)
}

type ExOrder struct {
	Value   interface{}
	Reorder []bool
}

func MakeOrder(c echo.Context, value interface{}, reorder ...bool) {
	tmp := make([]*ExOrder, 0)
	if ex, ok := c.Get(stOrderKey).([]*ExOrder); ok {
		tmp = ex
	}
	tmp = append(tmp, &ExOrder{
		Value:   value,
		Reorder: reorder,
	})
	c.Set(stOrderKey, tmp)
}

type ExOmit struct {
	Columns []string
}

func MakeOmit(c echo.Context, columns ...string) {
	tmp := make([]*ExOmit, 0)
	if ex, ok := c.Get(stOmitKey).([]*ExOmit); ok {
		tmp = ex
	}
	tmp = append(tmp, &ExOmit{
		Columns: columns,
	})
	c.Set(stOmitKey, tmp)
}

type ExSet struct {
	Name  string
	Value interface{}
}

func MakeSet(c echo.Context, name string, value interface{}) {
	tmp := make([]*ExSet, 0)
	if ex, ok := c.Get(stGormSetKey).([]*ExSet); ok {
		tmp = ex
	}
	tmp = append(tmp, &ExSet{
		Name:  name,
		Value: value,
	})
	c.Set(stGormSetKey, tmp)
}
