package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expeceted T) {
	// Go test runner will report the filename and line number of the code which called our Equal() function in the output.
	t.Helper()

	if actual != expeceted {
		t.Errorf("Expected: %v => Actual: %v", expeceted, actual)	
	}
}
