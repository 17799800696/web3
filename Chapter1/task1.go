package main

import (
	"fmt"
	"sort"
)

func RemoveDuplicates(nums []int) int {
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

func IsValid(s string) bool {
	mapping := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	stack := []rune{}

	for _, char := range s {

		if left, ok := mapping[char]; ok {

			if len(stack) == 0 || stack[len(stack)-1] != left {
				return false
			}

			stack = stack[:len(stack)-1]
		} else {

			stack = append(stack, char)
		}
	}

	return len(stack) == 0
}

func LongestCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }

    // 以第一个字符串为基准
    first := strs[0]
    
    // 逐列比较每个字符串的对应字符
    for i := 0; i < len(first); i++ {
        char := first[i]
        for j := 1; j < len(strs); j++ {
            // 如果当前字符串长度不足i，或者字符不匹配
            if i >= len(strs[j]) ||  strs[j][i]!= char {
                return first[:i] // 返回前缀，不包含索引i的字符
            }
        }
    }
    
    // 所有字符都匹配，第一个字符串本身就是最长公共前缀
    return first
}

func PlusOne(digits []int) []int {
    n := len(digits)
    
    for i := n - 1; i >= 0; i-- {
        digits[i]++
        if digits[i] < 10 {
            return digits
        }
        digits[i] = 0
    }
    result := make([]int, n+1)
    result[0] = 1
    return result
}

func Merge(intervals [][]int) [][]int {
    if len(intervals) <= 1 {
        return intervals
    }
    
    // 按区间的起始点排序
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })
    
    merged := [][]int{intervals[0]} // 初始化合并后的数组
    
    for _, current := range intervals[1:] {
        // 获取最后一个合并区间的结束点
        last := merged[len(merged)-1]
        
        if current[0] <= last[1] { // 有重叠，更新结束点
            if current[1] > last[1] {
                last[1] = current[1]
            }
        } else { // 无重叠，添加新区间
            merged = append(merged, current)
        }
    }
    
    return merged
}

func TwoSum(nums []int, target int) []int {
    // 哈希表存储值到索引的映射
    numMap := make(map[int]int)
    
    for i, num := range nums {
        // 计算需要的另一个数
        complement := target - num
        
        // 检查哈希表中是否存在该数，且不是当前索引
        if idx, exists := numMap[complement]; exists {
            return []int{idx, i}
        }
        
        // 将当前数存入哈希表
        numMap[num] = i
    }
    
    return nil
}

func Task2(x int) bool {
	if x < 0 {
		return false
	}
	original := x
	reversed := 0
	for x != 0 {
		reversed = reversed*10 + x%10
		x /= 10
	}
	return original == reversed
}

func main() {
	fmt.Println(RemoveDuplicates([]int{1, 1, 2}))
	fmt.Println(Task2(121))
	fmt.Println(IsValid("()"))
	fmt.Println(IsValid("(]"))
	fmt.Println(IsValid("([])"))
	fmt.Println(LongestCommonPrefix([]string{"flower", "flow", "flight"}))
	fmt.Println(PlusOne([]int{1, 2, 3}))
	fmt.Println(PlusOne([]int{9, 9, 9}))
	fmt.Println(PlusOne([]int{9, 6, 9}))
	fmt.Println(Merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
	fmt.Println(TwoSum([]int{2, 7, 11, 15}, 9))	
}
