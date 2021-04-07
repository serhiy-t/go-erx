package examples

import (
	"github.com/serhiy-t/go-erx"
)

type Location struct {
	Lat  float64
	Long float64
}

// url: /lat/67.452342343/long/67.452342343
func ParseLatLongUrl(url string) (l Location, out error) {
	defer erx.ErrFromPanic(&out)

	consumeKeyOrPanic(&url, "lat")
	l.Lat = atofOrPanic(consumeValueOrPanic(&url))
	consumeKeyOrPanic(&url, "long")
	l.Long = atofOrPanic(consumeValueOrPanic(&url))
	return
}
