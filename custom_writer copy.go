package response_analyser

import (
	"errors"
	"fmt"
)

type MyCustomWC struct {
	isClosed bool
	// You might store an underlying file, buffer, or network connection here
}

func (m *MyCustomWC) Write(p []byte) (n int, err error) {
	if m.isClosed {
		return 0, errors.New("closed writer")
	}
	fmt.Println(string(p))

	return len(p), nil
}

func (m *MyCustomWC) Close() error {
	m.isClosed = true
	// Implement logic to release resources, flush buffers, or close connections
	return nil
}
