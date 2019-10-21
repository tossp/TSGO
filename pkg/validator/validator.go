package validator

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"sync"

	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/null"

	enLocales "github.com/go-playground/locales/en"
	zhLocales "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
	zhTrans "gopkg.in/go-playground/validator.v9/translations/zh"
)

var (
	validate = new(defaultValidator)
)

type defaultValidator struct {
	once      sync.Once
	validator *validator.Validate
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// Validate For echo
func (v *defaultValidator) Validate(i interface{}) error {
	return v.Struct(i)
}

func (v *defaultValidator) Var(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}
func (v *defaultValidator) VarWithValue(field interface{}, other interface{}, tag string) error {
	return v.validator.VarWithValue(field, other, tag)
}
func (v *defaultValidator) Struct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validator.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validator
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		zH := zhLocales.New()
		uni = ut.New(zH, zH, enLocales.New())
		trans, has := FindTranslator("zh")
		if !has {
			log.Warn("验证翻译默认语种查找失败")
		}
		v.validator = validator.New()
		v.validator.SetTagName("valid")
		_ = zhTrans.RegisterDefaultTranslations(v.validator, trans)
		v.validator.RegisterCustomTypeFunc(ValidateDBType,
			sql.NullString{}, sql.NullInt64{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{},
			null.String{}, null.Time{}, null.Int{}, null.Float{}, null.Bool{},
		)
	})
}

func New() *defaultValidator {
	return &defaultValidator{}
}

func ValidateDBType(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}
	return nil
}

//func UserStructLevelValidation(sl validator.StructLevel) {
//
//	user := sl.Current().Interface().(User)
//
//	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
//		sl.ReportError(user.FirstName, "FirstName", "fname", "fnameorlname", "")
//		sl.ReportError(user.LastName, "LastName", "lname", "fnameorlname", "")
//	}
//
//	// plus can do more, even with different tag than "fnameorlname"
//}

//func ValidateUniq(fl validator.FieldLevel) bool {
//	var result struct{ Count int }
//	currentField, _, _ := fl.GetStructFieldOK()
//	table := modelTableNameMap[currentField.Type().Name()] // table name
//	value := fl.Field().String()                           // value
//	column := fl.FieldName()                               // column name
//	sql := fmt.Sprintf("select count(*) from %s where %s='%s'", table, column, value)
//	db.PG.Raw(sql).Scan(&result)
//	dup := result.Count > 0
//	return !dup
//}
