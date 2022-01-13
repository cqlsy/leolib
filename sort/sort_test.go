package sort

import (
	"fmt"
	"testing"
)

func TestQuickSort(t *testing.T) {
	arr := []int{1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12,1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12}
	//arr := []int{1, 9,}

	fmt.Println(QuickSort(arr))
	fmt.Println(count)
}

func TestBubbleSort(t *testing.T) {
	arr := []int{1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12,1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12}

	fmt.Println(BubbleSort(arr))
	fmt.Println(count)
}
