package sys

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
