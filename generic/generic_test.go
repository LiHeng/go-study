package generic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/samber/lo"
)

func TestGenericMin(t *testing.T) {
	fmt.Println(min(1, 2))
	fmt.Println(min(3.4, 3.5))
}

func TestGenericSet(t *testing.T) {
	set := NewSet[string]()
	set.Add("lili")
	set.Add("lucy")
	fmt.Println(set.Values())

	intSet := NewSet[int]()
	intSet.Add(2)
	intSet.Add(1)
	fmt.Println(intSet.Values())
}

func TestGeneticLibrary(t *testing.T) {
	fmt.Println(lo.Uniq([]string{"Samuel", "Marc", "Samuel"}))
	fmt.Println(lo.FlatMap([]int64{0, 1, 2}, func(x int64, _ int) []string {
		return []string{
			strconv.FormatInt(x, 10),
			strconv.FormatInt(x, 10),
		}
	}))
}
