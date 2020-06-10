package client

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

var (
	hdrUserAgentKey   = http.CanonicalHeaderKey("User-Agent")
	hdrUserAgentValue = "tossp-app/0.0.0 (+https://github.com/tossp)"
)

func New() *resty.Client {
	return resty.New().SetHeader(hdrUserAgentKey, hdrUserAgentValue)
}
func SetUserAgent(value string) {
	hdrUserAgentValue = value
}
