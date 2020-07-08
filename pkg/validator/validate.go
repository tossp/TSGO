package validator

import (
	"database/sql/driver"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tossp/tsgo/pkg/db"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	idCardCoefficient []int32 = []int32{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	idCardCode        []byte  = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	nl = struct{}{}
)

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

func ValidateIdCard(fl validator.FieldLevel) bool {
	idCardNo := fl.Field().String()
	if len(idCardNo) != 18 {
		return false
	}

	birthDay, err := time.Parse("20060102150405", idCardNo[6:14]+"000001")
	if err != nil {
		return false
	}
	if birthDay.After(time.Now()) {
		return false
	}
	idByte := []byte(strings.ToUpper(idCardNo))

	sum := int32(0)
	for i := 0; i < 17; i++ {
		sum += int32(byte(idByte[i])-byte('0')) * idCardCoefficient[i]
	}
	return idCardCode[sum%11] == idByte[17]
}
