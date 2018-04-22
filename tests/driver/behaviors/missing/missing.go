package missing

import crossdock "github.com/crossdock/crossdock-go"

// Execute implements the 'missing-trace-context' behavior.
func Execute(t crossdock.T) {
	t.Skipf("not implemented")
}
