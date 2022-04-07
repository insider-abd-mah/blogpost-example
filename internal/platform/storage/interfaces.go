package storage

// Reader interface for storage
type Reader interface {
	Read(file string) ([]byte, error)
}

// Writer interface for storage
type Writer interface {
	Write(file string, content []byte, contentType string) error
}
