package config

import (
	"fmt"
	"testing"
)

func TestOutputConfig(t *testing.T) {
	e := OutputConfig(DefaultConfig())
	if e != nil {
		fmt.Println("config wrong", "error")
		panic(e)
	}

}
