//// file: testutils.go

// Package gardener ...
// We're adding testutils to gardener since most go apps use gardener for testing anyways
package gardener

import "bytes"

//// ====== Structures ======

// MockRC ...
// Mocks buffer
type MockRC struct {
	*bytes.Buffer
}

//// ====== Public ======

// Close ...
// Closes mock readcloser
func (rc *MockRC) Close() (err error) {
	return
}

func panicCheck(err error) {
	if err != nil {
		panic(err)
	}
}
