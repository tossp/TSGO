package validator

import (
	"database/sql"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/tossp/tsgo/pkg/tstype"
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
		_ = v.validator.RegisterValidation("tsdbunique", ValidateUniq, true)
		_ = v.validator.RegisterValidation("tscmpn", ValidateChinaMobilePhoneNum, true)
		_ = v.validator.RegisterValidation("idcard", ValidateIdCard, true)
		v.validator.RegisterCustomTypeFunc(ValidateDBType,
			sql.NullString{}, sql.NullInt64{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{},
			tstype.UUID{},
			// valid:".{2,}?"
			tstype.Varchar{}, tstype.Text{}, tstype.Bool{}, tstype.Timestamptz{},
			tstype.Bool{}, tstype.Numeric{}, tstype.Timestamptz{},
			tstype.JSON{}, tstype.JSONB{}, tstype.Hstore{},
			tstype.UUID{}, tstype.UUIDArray{}, tstype.Inet{}, tstype.InetArray{},
		)

		trans := FindTranslator("zh")
		_ = zhTrans.RegisterDefaultTranslations(v.validator, trans)
		_ = v.validator.RegisterTranslation("alphanumunicode", trans, registrationFunc("alphanumunicode", "{0}只能包含字母、数字和汉字", false), translateFunc)
		_ = v.validator.RegisterTranslation("alphaunicode", trans, registrationFunc("alphaunicode", "{0}只能包含字母和汉字", false), translateFunc)
		_ = v.validator.RegisterTranslation("e164", trans, registrationFunc("e164", "{0}必须是一个有效的电话号码", false), translateFunc)
		_ = v.validator.RegisterTranslation("tsdbunique", trans, registrationFunc("tsdbunique", "{0}已经被其他记录使用，请更换", false), translateFunc)
		_ = v.validator.RegisterTranslation("tscmpn", trans, registrationFunc("tscmpn", "{0}必须是一个有效的国内手机号码", false), translateFunc)
		_ = v.validator.RegisterTranslation("idcard", trans, registrationFunc("idcard", "{0}无效，请核对", false), translateFunc)

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

func RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

func RegisterTranslation(tag string, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc, locales ...string) error {
	return validate.RegisterTranslation(tag, registerFn, translationFn, locales...)
}
func Validate() *TsValidator {
	return validate
}

//New 创建新的验证器
func New() *TsValidator {
	return &TsValidator{}
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
