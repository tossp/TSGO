package validator

import ut "github.com/go-playground/universal-translator"

var (
	uni *ut.UniversalTranslator
)

func FindTranslator(locales ...string) (trans ut.Translator, found bool) {
	if locales == nil {
		locales = []string{}
	}
	locales = append(locales, "zh")
	return uni.FindTranslator(locales...)
}
