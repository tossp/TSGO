package utils

import (
	"net"
	"net/http"
	"strings"
)

// HTTP Header keys
const (
	Age                           string = "Age"
	AltSCV                        string = "Alt-Svc"
	Accept                        string = "Accept"
	AcceptCharset                 string = "Accept-Charset"
	AcceptPatch                   string = "Accept-Patch"
	AcceptRanges                  string = "Accept-Ranges"
	AcceptedLanguage              string = "Accept-Language"
	AcceptEncoding                string = "Accept-Encoding"
	Authorization                 string = "Authorization"
	CrossOriginResourcePolicy     string = "Cross-Origin-Resource-Policy"
	CacheControl                  string = "Cache-Control"
	Connection                    string = "Connection"
	ContentDisposition            string = "Content-Disposition"
	ContentEncoding               string = "Content-Encoding"
	ContentLength                 string = "Content-Length"
	ContentType                   string = "Content-Type"
	ContentLanguage               string = "Content-Language"
	ContentLocation               string = "Content-Location"
	ContentRange                  string = "Content-Range"
	Date                          string = "Date"
	DeltaBase                     string = "Delta-Base"
	ETag                          string = "ETag"
	Expires                       string = "Expires"
	Host                          string = "Host"
	IM                            string = "IM"
	IfMatch                       string = "If-Match"
	IfModifiedSince               string = "If-Modified-Since"
	IfNoneMatch                   string = "If-None-Match"
	IfRange                       string = "If-Range"
	IfUnmodifiedSince             string = "If-Unmodified-Since"
	KeepAlive                     string = "Keep-Alive"
	LastModified                  string = "Last-Modified"
	Link                          string = "Link"
	Pragma                        string = "Pragma"
	ProxyAuthenticate             string = "Proxy-Authenticate"
	ProxyAuthorization            string = "Proxy-Authorization"
	PublicKeyPins                 string = "Public-Key-Pins"
	RetryAfter                    string = "Retry-After"
	Referer                       string = "Referer"
	Server                        string = "Server"
	SetCookie                     string = "Set-Cookie"
	StrictTransportSecurity       string = "Strict-Transport-Security"
	Trailer                       string = "Trailer"
	TK                            string = "Tk"
	TransferEncoding              string = "Transfer-Encoding"
	Location                      string = "Location"
	Upgrade                       string = "Upgrade"
	Vary                          string = "Vary"
	Via                           string = "Via"
	Warning                       string = "Warning"
	WWWAuthenticate               string = "WWW-Authenticate"
	XForwardedFor                 string = "X-Forwarded-For"
	XForwardedHost                string = "X-Forwarded-Host"
	XForwardedProto               string = "X-Forwarded-Proto"
	XRealIP                       string = "X-Real-Ip"
	XContentTypeOptions           string = "X-Content-Type-Options"
	XFrameOptions                 string = "X-Frame-Options"
	XXSSProtection                string = "X-XSS-Protection"
	XDNSPrefetchControl           string = "X-DNS-Prefetch-Control"
	Allow                         string = "Allow"
	Origin                        string = "Origin"
	AccessControlAllowOrigin      string = "Access-Control-Allow-Origin"
	AccessControlAllowCredentials string = "Access-Control-Allow-Credentials"
	AccessControlAllowHeaders     string = "Access-Control-Allow-Headers"
	AccessControlAllowMethods     string = "Access-Control-Allow-Methods"
	AccessControlExposeHeaders    string = "Access-Control-Expose-Headers"
	AccessControlMaxAge           string = "Access-Control-Max-Age"
	AccessControlRequestHeaders   string = "Access-Control-Request-Headers"
	AccessControlRequestMethod    string = "Access-Control-Request-Method"
	TimingAllowOrigin             string = "Timing-Allow-Origin"
	UserAgent                     string = "User-Agent"
)

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func ClientIP(r *http.Request) (clientIP string) {
	values := r.Header[XRealIP]
	if len(values) > 0 {
		clientIP = strings.TrimSpace(values[0])
		if clientIP != "" {
			return
		}
	}
	if values = r.Header[XForwardedFor]; len(values) > 0 {
		clientIP = values[0]
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if clientIP != "" {
			return
		}
	}
	clientIP, _, _ = net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	return
}

// AcceptedLanguages returns an array of accepted languages denoted by
// the Accept-Language header sent by the browser
func AcceptedLanguages(r *http.Request) (languages []string) {
	accepted := r.Header.Get(AcceptedLanguage)
	if accepted == "" {
		return
	}
	options := strings.Split(accepted, ",")
	l := len(options)
	languages = make([]string, l)

	for i := 0; i < l; i++ {
		locale := strings.SplitN(options[i], ";", 2)
		languages[i] = strings.Trim(locale[0], " ")
	}
	return
}
