package public

import (
	"strings"
	"unicode"
)

//是否包含汉字
func IsHan(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func CookieToMap(cookieStr string) map[string]string {
	mp := map[string]string{}
	for _, cookie := range strings.Split(cookieStr, "; ") {
		idx := strings.Index(cookie, "=")
		if idx < 0 {
			continue
		}
		mp[cookie[:idx]] = cookie[idx+1:]
	}
	return mp
}
