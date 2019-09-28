package net

import "testing"

// TestRegisterProxy ...
func TestRegisterProxy(t *testing.T) {
	RegisterProxy("socks5://localhost:18080")
}
