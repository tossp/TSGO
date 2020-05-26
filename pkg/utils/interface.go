package utils

import "net"

type LikeContextLog interface {
	Log(k, v string)
	Ip() net.IP
}
