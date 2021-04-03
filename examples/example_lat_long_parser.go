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

	panicsConsumeKey(&url, "lat")
	l.Lat = panicsAtof(panicsConsumeValue(&url))
	panicsConsumeKey(&url, "long")
	l.Long = panicsAtof(panicsConsumeValue(&url))
	return
}
