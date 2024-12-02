package sys

import (
	"fmt"
	"strconv"
	"strings"
)

// 递归生成多维矩阵的下标
func generateIndices_v3(sizes []int, result *[][]int) {
	var helper func([]int, int)
	helper = func(current []int, depth int) {
		if depth == len(sizes) {
			// 复制当前下标并添加到结果中
			index := make([]int, len(current))
			copy(index, current)
			*result = append(*result, index)
			return
		}

		for i := 0; i < sizes[depth]; i++ {
			current[depth] = i
			helper(current, depth+1)
		}
	}

	helper(make([]int, len(sizes)), 0)
}

// 递归生成多维矩阵的下标
func GenerateIndices(sizes []int, current []int, result *[][]int) {
	if len(current) == len(sizes) {
		// 复制当前下标并添加到结果中
		index := make([]int, len(current))
		copy(index, current)
		*result = append(*result, index)
		return
	}

	for i := 1; i <= sizes[len(current)]; i++ {
		GenerateIndices(sizes, append(current, i), result)
	}
}

func Indices(sizes []int, current []int) (res [][]int) {
	if len(current) == len(sizes) {
		// 复制当前下标并添加到结果中
		index := make([]int, len(current))
		copy(index, current)
		res = append(res, index)
		return
	}

	for i := 1; i <= sizes[len(current)]; i++ {
		res = append(res, Indices(sizes, append(current, i))...)
	}
	return
}

func GenArryList(dimensions [][]int) (res [][]int) {
	var helper func(current []int, depth int)
	helper = func(current []int, depth int) {
		if depth == len(dimensions) {
			index := make([]int, len(current))
			copy(index, current)
			res = append(res, index)
			return
		}

		for _, val := range dimensions[depth] {
			helper(append(current, val), depth+1)
		}
	}

	helper([]int{}, 0)
	return
}

func ParseRange(input string) ([]int, error) {
	var result []int
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			start, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err != nil {
				return nil, err
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			result = append(result, num)
		}
	}
	return result, nil
}
