package examples

import (
	"github.com/serhiy-t/go-erx"
	"io"
	"os"
)

func CopyFile(src string, dst string) (out error) {
	defer erx.WrapFmtErrorw(&out, "copy %s %s", src, dst)

	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer erx.IgnoreErr(r.Close)

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer erx.OnError(&out, func() { _ = os.Remove(dst) })
	defer erx.CheckErr(&out, w.Close)

	return erx.ResultErr(io.Copy(w, r))
}
