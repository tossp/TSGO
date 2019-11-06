package validator

import (
	enLocales "github.com/go-playground/locales/en"
	zhLocales "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/tossp/tsgo/pkg/log"
)

var (
	uni *ut.UniversalTranslator
)

func init() {
	zH := zhLocales.New()
	eN := enLocales.New()
	uni = ut.New(zH, zH, eN)
}

func FindTranslator(locales ...string) (trans ut.Translator) {
	if locales == nil {
		locales = []string{}
	}
	locales = append(locales, "zh")
	trans, has := uni.FindTranslator(locales...)
	if !has {
		log.Warn("验证翻译默认语种查找失败")
	}
	return
}
