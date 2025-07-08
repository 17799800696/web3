package main

import (
	"fmt"
	"github.com/test/init_project/Chapter1"
)

func main() {
	fmt.Println(Chapter1.RemoveDuplicates([]int{1, 1, 2}))
	fmt.Println(Chapter1.Task2(121))
	fmt.Println(Chapter1.IsValid("()"))
	fmt.Println(Chapter1.IsValid("(]"))
	fmt.Println(Chapter1.IsValid("([])"))
	fmt.Println(Chapter1.LongestCommonPrefix([]string{"flower", "flow", "flight"}))
	fmt.Println(Chapter1.PlusOne([]int{1, 2, 3}))
	fmt.Println(Chapter1.PlusOne([]int{9, 9, 9}))
	fmt.Println(Chapter1.PlusOne([]int{9, 6, 9}))
	fmt.Println(Chapter1.Merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
	fmt.Println(Chapter1.TwoSum([]int{2, 7, 11, 15}, 9))
}
