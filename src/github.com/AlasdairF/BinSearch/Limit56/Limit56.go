package sortLimit56

/*
	This package is specifically used by github.com/AlasdairF/BinSearch
*/

// ================= COMMON =================

type Slice []KeyVal
type KeyVal struct {
	K int
	V [7]uint64
}

func (a Slice) less(i, j int) bool {
	switch {
		case a[i].V[0] < a[j].V[0]: return true
		case a[i].V[0] > a[j].V[0]: return false
		case a[i].V[1] < a[j].V[1]: return true
		case a[i].V[1] > a[j].V[1]: return false
		case a[i].V[2] < a[j].V[2]: return true
		case a[i].V[2] > a[j].V[2]: return false
		case a[i].V[3] < a[j].V[3]: return true
		case a[i].V[3] > a[j].V[3]: return false
		case a[i].V[4] < a[j].V[4]: return true
		case a[i].V[4] > a[j].V[4]: return false
		case a[i].V[5] < a[j].V[5]: return true
		case a[i].V[5] > a[j].V[5]: return false
		case a[i].V[6] < a[j].V[6]: return true
		default: return false
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ------------- ASCENDING -------------

func heapSortAsc(data Slice, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownAsc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownAsc(data, lo, i, first)
	}
}

func insertionSortAsc(data Slice, a, b int) {
	var j int
	for i := a + 1; i < b; i++ {
		for j = i; j > a && data.less(j, j-1); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDownAsc(data Slice, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.less(first+child, first+child+1) {
			child++
		}
		if !data.less(first+root, first+child) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func medianOfThreeAsc(data Slice, m1, m0, m2 int) {
	// bubble sort on 3 elements
	if data.less(m1, m0) {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data.less(m2, m1) {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data.less(m1, m0) {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRangeAsc(data Slice, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a], data[b] = data[b], data[a]
		a++
		b++
	}
}

func doPivotAsc(data Slice, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeAsc(data, lo, lo+s, lo+2*s)
		medianOfThreeAsc(data, m, m-s, m+s)
		medianOfThreeAsc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeAsc(data, lo, m, hi-1)

	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data.less(b, pivot) {
				b++
			} else if !data.less(pivot, b) {
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data.less(pivot, c-1) {
				c--
			} else if !data.less(c-1, pivot) {
				data[c-1], data[d-1] = data[d-1], data[c-1]
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}

	n := min(b-a, a-lo)
	swapRangeAsc(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeAsc(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortAsc(data Slice, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortAsc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotAsc(data, a, b)
		if mlo-a < b-mhi {
			quickSortAsc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSortAsc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		insertionSortAsc(data, a, b)
	}
}

func Asc(data Slice) {
	maxDepth := 0
	for i := len(data); i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortAsc(data, 0, len(data), maxDepth)
}
