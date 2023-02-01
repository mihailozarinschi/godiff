package godiff

import "io"

type ReaderAt interface {
	io.Reader
	io.ReaderAt
}

// NewReaderAt provides a wrapper around an io.ReadSeeker implementation that doesn't implement io.ReaderAt
func NewReaderAt(rs io.ReadSeeker) ReaderAt {
	return &readerAt{rs}
}

type readerAt struct {
	io.ReadSeeker
}

func (r *readerAt) ReadAt(b []byte, off int64) (n int, err error) {
	_, err = r.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return r.Read(b)
}
