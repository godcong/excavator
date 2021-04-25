package strokefix

import (
	"excavator"
	"fmt"
	"testing"
	"time"
)

func TestNumberChar(t *testing.T) {
	before := time.Now()

	excH := excavator.New()

	NumberChar(excH)

	fmt.Printf("took %v\n", time.Since(before))
}
