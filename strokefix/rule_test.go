package strokefix

import (
	"fmt"
	"testing"
	"time"

	"github.com/godcong/excavator"
)

func TestNumberChar(t *testing.T) {
	before := time.Now()

	excH := excavator.New()

	NumberChar(excH)

	fmt.Printf("took %v\n", time.Since(before))
}
