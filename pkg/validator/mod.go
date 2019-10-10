package validator

func Struct(i interface{}) error {
	return validate.Struct(i)
}
