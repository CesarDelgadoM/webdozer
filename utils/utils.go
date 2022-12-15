package utils

import (
	"strings"
)

func ExtractNameUrl(url string) (key string) {

	partsUrl := strings.Split(url, ".")
	key = partsUrl[1]
	return
}
