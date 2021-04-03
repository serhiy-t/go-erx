package examples

import (
	"compress/gzip"
	"github.com/serhiy-t/go-erx"
	"io"
	"os"
)

type CountBytesWriter struct {
	writer io.WriteCloser
	bytes  int64
}

func (w *CountBytesWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	w.bytes += int64(n)
	return
}

func (w *CountBytesWriter) Close() error {
	return w.writer.Close()
}

func (w *CountBytesWriter) Bytes() int64 {
	return w.bytes
}

func ByteCounterWrap(writer io.WriteCloser, err error) (*CountBytesWriter, error) {
	return &CountBytesWriter{writer: writer}, err
}

type GzipStats struct {
	Compressed   int64
	Uncompressed int64
}

// GzipFile compresses file srcFilename into dstFilename.
func GzipFile(dstFilename string, srcFilename string) (stats GzipStats, out error) {
	e, ret := erx.ErrorRef(&out)
	defer e.OnError(func() { stats.Compressed = -1; stats.Uncompressed = -1 })
	defer e.WrapFmtErrorw("error compressing file %s to %s", srcFilename, dstFilename)
	defer e.LogSilentErrors()

	if ret(
		erx.AssertErr(len(dstFilename) > 0, "dst file should be specified"),
		erx.AssertErr(len(srcFilename) > 0, "src file should be specified")) {
		return
	}

	reader, err := os.Open(srcFilename)
	if ret(err) {
		return
	}
	defer e.IgnoreErr(reader.Close)

	writer, err := ByteCounterWrap(os.Create(dstFilename))
	if ret(err) {
		return
	}
	defer e.OnErrorOrPanic(func() { _ = os.Remove(dstFilename) })
	defer e.OnSuccess(func() { stats.Compressed = writer.Bytes() })
	defer e.CheckErr(writer.Close)

	gzipWriter := gzip.NewWriter(writer)
	defer e.CheckErr(gzipWriter.Close)

	stats.Uncompressed, err = io.Copy(gzipWriter, reader)
	if ret(err) {
		return
	}

	return
}
