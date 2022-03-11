package regular

import (
	"fmt"
	"testing"
	"time"

	"github.com/godcong/excavator"
)

// TestNew ...
func TestNew(t *testing.T) {
	before := time.Now()

	excH := excavator.New()

	regular := New(excH)
	regular.Run()

	fmt.Printf("took %v\n", time.Since(before))
}
