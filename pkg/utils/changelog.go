package utils

import (
	"fmt"
	"github.com/tossp/tsgo/pkg/utils/structs"
	"time"
)

type Change struct {
	Key  string      `json:"key"`
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

type Changelog []Change

type valueOrZeroTime interface {
	ValueOrZero() time.Time
}
type valueOrZeroString interface {
	ValueOrZero() string
}

func MakeChangelog(m, m1 interface{}) (cl Changelog) {
	s := structs.New(m)
	s1 := structs.New(m1)
	for _, v := range s.Fields() {
		if v.Tag("diff") == "" {
			continue
		}
		if vs, vs1 := makeChangeV(v.Value()), makeChangeV(s1.Field(v.Name()).Value()); vs != vs1 {
			cl = append(cl, Change{v.Tag("desc"), v.Value(), s1.Field(v.Name()).Value()})
			//log.Debug("makeChangelog", v.Name(), v.Tag("desc"), vs, vs1)
		}
	}
	return
}
func makeChangeV(v interface{}) (d string) {
	if s, ok := v.(fmt.Stringer); ok {
		d = s.String()
	} else if s, ok := v.(valueOrZeroString); ok {
		d = s.ValueOrZero()
	} else if s, ok := v.(valueOrZeroTime); ok {
		d = fmt.Sprintf("%d", s.ValueOrZero().UnixNano())
	} else {
		d = fmt.Sprintf("%v", v)
	}
	return
}
