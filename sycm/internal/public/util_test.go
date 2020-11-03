package public

import (
	"net/url"
	"testing"
)

func TestCookieToMap(t *testing.T) {
	cookieStr := "_samesite_flag_=true; cookie2=22456b=5194e372=1c21368ed0g695f9b; =100; sn=%E5%A4%A9%E7%8C%AB%E6%97%97%E8%88%B0%E5%BA%97%3A%E5%90%89%E5%A7%86"
	cookieMap := CookieToMap(cookieStr)
	for k, v := range cookieMap {
		t.Logf("[%v]\t[%v]", k, v)
	}
}

func TestIsHan(t *testing.T) {
	encodeStr := "%E5%A4%A9%E7%8C%AB%E6%97%97%E8%88%B0%E5%BA%97%3A%E5%90%89%E5%A7%86"
	t.Logf("%v is 汉字 ? %v", encodeStr, IsHan(encodeStr))
	han, err := url.QueryUnescape(encodeStr)
	if err != nil {
		t.Errorf("urldecode error: %v", err.Error())
	}
	t.Logf("%v is 汉字 ? %v", han, IsHan(han))
	if han != "天猫旗舰店:吉姆" {
		t.Errorf("urldecode res unexpected %v", han)
	}
}
