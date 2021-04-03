package examples

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/serhiy-t/go-erx"
)

func panicsAtoi(s string) int {
	i, err := strconv.Atoi(s)
	erx.PanicIf(err)
	return i
}

func panicsConsumeSlash(url *string) {
	if len(*url) == 0 || (*url)[0] != '/' {
		panic(fmt.Errorf("expected '/' next, but got: %s", *url))
	}
	*url = (*url)[1:]
}

func panicsConsumeKey(url *string, key string) {
	panicsConsumeSlash(url)
	if !strings.HasPrefix(*url, key) {
		panic(fmt.Errorf("expected '%s' next, but got: %s", key, *url))
	}
	*url = (*url)[len(key):]
}

func panicsConsumeValue(url *string) string {
	panicsConsumeSlash(url)
	slashPos := strings.Index(*url, "/")
	if slashPos < 0 {
		slashPos = len(*url)
	}
	result := (*url)[:slashPos]
	*url = (*url)[slashPos:]
	return result
}

func panicsAtof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	erx.PanicIf(err)
	return f
}
