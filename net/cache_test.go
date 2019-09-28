package net

import "testing"

// TestNewCache ...
func TestNewCache(t *testing.T) {
	cache := NewCache("./tmp")
	url := "https://pics.javbus.com/cover/6qx9_b.jpg"
	t.Log(cache.Get(url))

	t.Log(cache.Save(url, "./save/image.jpg"))

}
