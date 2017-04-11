package pushy

import (
	"net/url"
)

func findPushy(serviceName string, useTLS bool) (url.URL, error) {
	ret, err := url.Parse("https://api.pushy.me/")
	if err != nil || ret == nil {
		return url.URL{}, err
	}
	return *ret, err
}