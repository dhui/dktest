package mockdockerclient

// MockReader is a mock implementation of the io.Reader interface
type MockReader struct {
	N   int
	Err error
}

// Read is a mock implementation of io.Reader.Read()
func (r MockReader) Read(p []byte) (n int, err error) {
	return r.N, r.Err
}

// MockCloser is a mock implementation of the io.Closer interface
type MockCloser struct {
	Err error
}

// Close is a mock implementation of io.Closer.Close()
func (c MockCloser) Close() error { return c.Err }

// MockReadCloser is a mock implementation of the io.ReadCloser interface
type MockReadCloser struct {
	MockReader
	MockCloser
}
