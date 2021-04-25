package regular

import (
	"excavator"
	"fmt"
	"testing"
	"time"
)

// TestNew ...
func TestNew(t *testing.T) {
	before := time.Now()

	excH := excavator.New()

	regular := New(excH)
	regular.Run()

	fmt.Printf("took %v\n", time.Since(before))
}
