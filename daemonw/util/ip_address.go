package util

import (
	"strings"
	"net/http"
	"net"
)

func GetRequestIP(r *http.Request, isForwarded bool) string {
	if isForwarded {

		//forward by LB or proxy
		clientIp := r.Header.Get("X-Forwarded-For")
		if index := strings.IndexByte(clientIp, ','); index >= 0 {
			clientIp = clientIp[:index]
		}
		clientIp = strings.TrimSpace(clientIp)
		if len(clientIp) > 0 {
			return clientIp
		}

		//forward by nginx
		clientIP := strings.TrimSpace(r.Header.Get("X-Real-Ip"))
		if len(clientIP) > 0 {
			return clientIP
		}

		//forward by Google App Engine
		clientIp = strings.TrimSpace(r.Header.Get("X-Appengine-Remote-Addr"))
		if len(clientIp) > 0 {
			return clientIp;
		}
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip;
	}
	return ""
}
