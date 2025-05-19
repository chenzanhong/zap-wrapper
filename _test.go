// zap_wrapper_test.go
package zap_wrapper

import "testing"

func TestInfow(t *testing.T) {
	msg := Infow("test message", "key", "value")
	if msg == "" {
		t.Fail()
	}
}