package concurrent

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestAtomic(t *testing.T) {
	num := int32(100)
	atomic.AddInt32(&num, 1)
	fmt.Println("num: ", num)

	atomic.CompareAndSwapInt32(&num, 99, num+1)
	fmt.Println("num: ", num)
}
