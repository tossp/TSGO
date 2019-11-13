package jwt

import (
	"sync"
	"time"
)

var sess sync.Map

func SessStart() {
	t := time.NewTicker(time.Second * 1)
	for {
		select {
		case now := <-t.C:
			utc := now.Unix()
			f := func(k, v interface{}) bool {
				c, _ := v.(*TsClaims)
				if !c.VerifyExpiresAt(utc, true) {
					sess.Delete(k)
				}
				return true
			}
			sess.Range(f)
		}
	}
}
func SessGet(key string) (data *TsClaims, ok bool) {
	value, ok := sess.Load(key)
	if !ok {
		return
	}
	data, _ = value.(*TsClaims)
	if !data.VerifyExpiresAt(time.Now().Unix(), true) {
		ok = false
		data = nil
		sess.Delete(key)
	}
	return
}

func SessSet(key string, value *TsClaims) {
	sess.Store(key, value)
	return
}
func SessDel(key string) {
	sess.Delete(key)
	return
}
