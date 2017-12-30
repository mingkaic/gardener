//// file: testutils.go

// We're adding testutils to gardener since most go apps use gardener for testing anyways
package gardener

import "bytes"

//// ====== Structures ======

type MockRC struct {
	*bytes.Buffer
}

//// ====== Public ======

// close mock readcloser
func (rc *MockRC) Close() (err error) {
	return
}

func panicCheck(err error) {
	if err != nil {
		panic(err)
	}
}
