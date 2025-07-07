package Chapter1

import "fmt"


func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	k := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[k-1] {
			nums[k] = nums[i]
			k++
		}
	}
	return k
}

func Task1() {
	var num = removeDuplicates([]int{1, 1, 2})
	fmt.Println(num)
}
