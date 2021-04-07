package examples

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/serhiy-t/go-erx"
)

func consumeSlashOrPanic(url *string) {
	if len(*url) == 0 || (*url)[0] != '/' {
		panic(fmt.Errorf("expected '/' next, but got: %s", *url))
	}
	*url = (*url)[1:]
}

func consumeKeyOrPanic(url *string, key string) {
	consumeSlashOrPanic(url)
	if !strings.HasPrefix(*url, key) {
		panic(fmt.Errorf("expected '%s' next, but got: %s", key, *url))
	}
	*url = (*url)[len(key):]
}

func consumeValueOrPanic(url *string) string {
	consumeSlashOrPanic(url)
	slashPos := strings.Index(*url, "/")
	if slashPos < 0 {
		slashPos = len(*url)
	}
	result := (*url)[:slashPos]
	*url = (*url)[slashPos:]
	return result
}

func atofOrPanic(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	erx.PanicIf(err)
	return f
}
