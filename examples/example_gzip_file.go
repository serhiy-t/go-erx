package examples

import (
	"compress/gzip"
	"fmt"
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
var noGzipStats = GzipStats{Compressed: -1, Uncompressed: -1}

// GzipFile compresses file srcFilename into dstFilename.
func GzipFile(dstFilename string, srcFilename string) (stats GzipStats, out error) {
	defer erx.WrapFmtErrorw(&out,"error compressing file %s to %s", srcFilename, dstFilename)

	errLogger := erx.ErrorLogger(&out)
	defer errLogger.LogSilentErrors()

	if len(dstFilename) == 0 {
		return noGzipStats, fmt.Errorf("dst file should be specified")
	}
	if len(srcFilename) == 0 {
		return noGzipStats, fmt.Errorf("src file should be specified")
	}

	reader, err := os.Open(srcFilename)
	if err != nil {
		return noGzipStats, err
	}
	defer erx.IgnoreErr(reader.Close, errLogger)

	writer, err := ByteCounterWrap(os.Create(dstFilename))
	if err != nil {
		return noGzipStats, err
	}
	defer erx.OnErrorOrPanic(&out, func() { _ = os.Remove(dstFilename) })
	defer erx.OnSuccess(&out, func() { stats.Compressed = writer.Bytes() })
	defer erx.CheckErr(&out, writer.Close, errLogger)

	gzipWriter := gzip.NewWriter(writer)
	defer erx.CheckErr(&out, gzipWriter.Close, errLogger)

	stats.Uncompressed, err = io.Copy(gzipWriter, reader)
	if err != nil {
		return noGzipStats, err
	}

	return
}
