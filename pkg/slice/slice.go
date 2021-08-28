package slice

// RemoveDuplicationSort 去重并保持排序
func RemoveDuplicationSort(arr []string) []string {
	length := len(arr)
	if length == 0 {
		return arr
	}

	j := 0
	for i := 1; i < length; i++ {
		if arr[i] != arr[j] {
			j++
			if j < i {
				swap(arr, i, j)
			}
		}
	}

	return arr[:j+1]
}

func swap(arr []string, a, b int) {
	arr[a], arr[b] = arr[b], arr[a]
}
