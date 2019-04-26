package slice

import "sort"

// Int64 int64类型切片,支持自动排序,快速查找
type Int64 []int64

// find 在int64切片中快速查找指定值
// 返回 值索引、是否存在
func (i64 Int64) find(v int64) (int, bool) {
	l := len(i64)
	if l == 0 || i64[0] > v {
		return 0, false
	}

	if v > i64[l-1] {
		return l, false
	}

	idx := sort.Search(l, func(i int) bool {
		return v <= i64[i]
	})

	return idx, idx < l && i64[idx] == v
}

// Add 向int64切片中增加值
func (i64 *Int64) Add(v int64) bool {
	idx, found := i64.find(v)
	if found {
		return false
	}

	*i64 = append(*i64, 0)
	copy((*i64)[idx+1:], (*i64)[idx:])
	(*i64)[idx] = v
	return true
}

// Rem 从int64切片中删除指定值
func (i64 *Int64) Rem(v int64) bool {
	idx, found := i64.find(v)
	if !found {
		return false
	}
	if idx == len(*i64)-1 {
		*i64 = (*i64)[:idx]
	} else {
		*i64 = append((*i64)[:idx], (*i64)[idx+1:]...)
	}
	return true
}

// Contains 查找int64切片是否存在值
func (i64 Int64) Contains(v int64) bool {
	_, contains := i64.find(v)
	return contains
}
