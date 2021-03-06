package stream

import (
	"fmt"
	"testing"
)

func TestDistinct(t *testing.T) {
	Just(1, 1, 2, 2, 3, 5).Distinct(func(item any) any {
		return item
	}).ForEach(func(item any) {
		fmt.Println(item)
	})

	Just(1, 1, 2, 2, 3, 5, 5).Distinct(func(item any) any {
		uid := item.(int)
		// 对大于4的item进行特殊去重逻辑,最终只保留一个>3的item
		if uid > 3 {
			return 4
		}
		return item
	}).ForEach(func(item any) {
		fmt.Println(item)
	})
}

func TestFilter(t *testing.T) {
	Just(1, 2, 3, 4, 5).Filter(func(item any) bool {
		return item.(int)%2 != 0
	}).ForEach(func(item any) {
		fmt.Println(item)
	})

	Just(7, 8, 9, 4, 5).Filter(func(item any) bool {
		return item.(int)%2 != 0
	}, WithWorkers(1)).ForEach(func(item any) {
		fmt.Println(item)
	})
}

func TestSort(t *testing.T) {
	Just(1, 2, 3, 4, 5).Sort(func(a, b any) bool {
		return a.(int) > b.(int)
	}).ForEach(func(item any) {
		fmt.Println(item)
	})
}
