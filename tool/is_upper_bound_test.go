package tool

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsUpperBound(t *testing.T) {
	var d, i, s = 5, 0, 0
	for ; IsUpperBound(i, time.Duration(d)); i++ {
		s = (1 << i)
	}
	assert.True(t, s < d)
	d, i, s = 10, 0, 0
	for ; IsUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
	}
	assert.True(t, (1<<i) < d)
	d, i, s = 25, 0, 0
	for ; IsUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
	}
	assert.True(t, (1<<i) < d)
	d, i, s = 50, 0, 0
	for ; IsUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
	}
	assert.True(t, (1<<i) < d)
	for i := 1; IsUpperBound(i, 1000*time.Millisecond); i++ {
		time.Sleep(time.Millisecond * time.Duration(i) * 100)
		fmt.Fprintf(os.Stderr, "i: %d, sleep: %v, UpperBound: %v\n", i, 100*time.Millisecond*time.Duration(i), 1000*time.Millisecond)
	}
}

func TestIsUpperBoundWithSleep(t *testing.T) {
	for i := 1; IsUpperBoundWithSleep(i, 1000, time.Millisecond); i++ {
		time.Sleep(time.Millisecond * time.Duration(i) * 100)
	}
}
