package validator

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/tossp/tsgo/pkg/db"
	"github.com/tossp/tsgo/pkg/null"
	//"github.com/tossp/tsgo/pkg/utils"
	//"github.com/jinzhu/inflection"
)

var (
	validate = new(TsValidator)
)

type TsValidator struct {
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
func (v *TsValidator) Validate(i interface{}) error {
	v.lazyinit()
	return v.validator.Struct(i)
}

func (v *TsValidator) Var(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}
func (v *TsValidator) VarWithValue(field interface{}, other interface{}, tag string) error {
	return v.validator.VarWithValue(field, other, tag)
}
func (v *TsValidator) Struct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validator.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *TsValidator) Engine() interface{} {
	v.lazyinit()
	return v.validator
}

func (v *TsValidator) lazyinit() {
	v.once.Do(func() {
		v.validator = validator.New()
		v.validator.SetTagName("valid")
		v.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("desc")
		})
		v.validator.RegisterValidation("tsdbunique", ValidateUniq, true)
		v.validator.RegisterValidation("tscmpn", ValidateChinaMobilePhoneNum, true)
		v.validator.RegisterCustomTypeFunc(ValidateDBType,
			sql.NullString{}, sql.NullInt64{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{},
			null.Bool{}, null.CIDR{}, null.Float{}, null.Int{}, null.IP{}, null.String{}, null.Time{}, null.UUID{},
		)

		trans := FindTranslator("zh")
		_ = zhTrans.RegisterDefaultTranslations(v.validator, trans)
		_ = v.validator.RegisterTranslation("alphanumunicode", trans, registrationFunc("alphanumunicode", "{0}只能包含字母、数字和汉字", false), translateFunc)
		_ = v.validator.RegisterTranslation("alphaunicode", trans, registrationFunc("alphaunicode", "{0}只能包含字母和汉字", false), translateFunc)
		_ = v.validator.RegisterTranslation("e164", trans, registrationFunc("e164", "{0}必须是一个有效的电话号码", false), translateFunc)
		_ = v.validator.RegisterTranslation("tsdbunique", trans, registrationFunc("tsdbunique", "{0}已经被其他记录使用，请更换", false), translateFunc)
		_ = v.validator.RegisterTranslation("tscmpn", trans, registrationFunc("tscmpn", "{0}必须是一个有效的国内手机号码", false), translateFunc)

	})
}
func (v *TsValidator) RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	v.lazyinit()
	return v.validator.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}
func (v *TsValidator) RegisterTranslation(tag string, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc, locales ...string) error {
	v.lazyinit()
	trans := FindTranslator(locales...)
	return v.validator.RegisterTranslation(tag, trans, registerFn, translationFn)
}

//New 创建新的验证器
func New() *TsValidator {
	return &TsValidator{}
}

var nl = struct{}{}

//ValidateDBType 验证数据库类型
func ValidateDBType(field reflect.Value) (val interface{}) {
	var err error
	valuer, ok := field.Interface().(driver.Valuer)
	if ok {
		val, err = valuer.Value()
		if err == nil {
			if val == nil {
				return nl
			}
			return val

		}
	}
	return nil
}

func ValidateUniq(fl validator.FieldLevel) bool {
	currentField, _, _, _ := fl.GetStructFieldOK2()
	table := currentField.Type().Name() // table name
	value := fl.Field().String()        // value
	column := fl.StructFieldName()      // column name
	var result int64
	q := db.G().
		//Debug().
		Table(db.TableName(table)).
		Where(fmt.Sprintf("%s=?", db.ColumnName(column)), value)
	uid := currentField.FieldByName("UID")
	if !uid.IsZero() {
		q = q.Where("uid!=?", uid.Interface())
	}
	q = q.Count(&result)
	if q.Error != nil {
		return false
	}
	return result == 0
}
func ValidateChinaMobilePhoneNum(fl validator.FieldLevel) bool {
	return regexp.MustCompile("^(?:(?:\\+|00)86)?1(?:(?:3[\\d])|(?:4[5-7|9])|(?:5[0-3|5-9])|(?:6[5-7])|(?:7[0-8])|(?:8[\\d])|(?:9[1|8|9]))\\d{8}$").
		MatchString(fl.Field().String())
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
