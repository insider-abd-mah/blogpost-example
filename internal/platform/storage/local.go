package storage

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Local storage
type Local struct{}

// Read file from local storage
func (*Local) Read(file string) ([]byte, error) {
	if os.Getenv("GIN_MODE") == "testing" {
		return ioutil.ReadFile(fmt.Sprintf("../../%s", file))
	}

	return ioutil.ReadFile(file)
}
