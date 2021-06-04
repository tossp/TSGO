package utils

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/tossp/tsgo/pkg/log"
)

var (
	namespaceTS = uuid.Must(uuid.FromString("19860912-9dad-11d1-80b4-00c04fd430c8"))
	namespaceZH = uuid.Must(uuid.FromString("00000001-1986-0912-0511-013419290962"))
)

func NamespaceTS() uuid.UUID {
	return namespaceTS
}
func NamespaceZH() uuid.UUID {
	return namespaceZH
}
func NewUuidV4() uuid.UUID {
	u, err := uuid.NewV4()
	if err != nil {
		log.Warn("生成uuid V4错误", err)
		time.Sleep(time.Second)
		return NewUuidV4()
	}
	return u
}
func NewUuidV5(ns uuid.UUID, name string) uuid.UUID {
	return uuid.NewV5(ns, name)
}
func UuidFromInterface(input interface{}) uuid.UUID {
	if value, ok := input.(interface{ Get() interface{} }); ok {
		return UuidFromInterface(value.Get())
	}
	switch value := input.(type) {
	case [16]byte:
		return uuid.FromBytesOrNil(value[:])
	case []byte:
		return uuid.FromBytesOrNil(value)
	case string:
		return uuid.FromStringOrNil(value)
	case *string:
		return uuid.FromStringOrNil(*value)
	default:
		return uuid.FromStringOrNil(fmt.Sprintf("%s", value))
	}
}
func UuidIsZero(input interface{}) bool {
	return bytes.Equal(uuid.Nil.Bytes(), UuidFromInterface(input).Bytes())
}
