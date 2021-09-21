package slice

import "testing"

func TestRemoveDuplicationSort(t *testing.T) {
	s := []string{"1", "1", "2", "3"}
	t.Log(RemoveDuplicationSort(s))
}
