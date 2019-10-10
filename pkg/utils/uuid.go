package utils

import (
	uuid "github.com/satori/go.uuid"
)

const (
	//ZeroUUID uuid起始值或者叫默认值
	ZeroUUID = "00000000-0000-0000-0000-000000000000"
	//namespaceUUID uuid1.1986-09-12T11:30:00+08:00
	namespaceUUID = "b4b02c01-d5c3-11c4-b381-abb0ec78cc40"
)

var (
	NamespaceUUIDMenu = uuid.Must(uuid.FromString("00000001-1986-0912-0511-013419290962"))
)

func NewUuidV4() uuid.UUID {
	// tmp, _ := uuid.NewV4()
	// return tmp.String()
	return uuid.NewV4()
}
func NewUuidV5(ns uuid.UUID, name string) uuid.UUID {
	// tmp, _ := uuid.NewV5(ns,name)
	// return tmp.String()
	return uuid.NewV5(ns, name)
}

func NewUuid(name string) uuid.UUID {
	return uuid.NewV5(NamespaceUUIDMenu, name)
}

//ValidUUID 验证 uuid
func ValidUUID(input string) (uuid.UUID, error) {
	return uuid.FromString(input)
}
