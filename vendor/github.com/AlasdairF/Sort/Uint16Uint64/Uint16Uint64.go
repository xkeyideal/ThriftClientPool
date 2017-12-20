package sortUint16Uint64

// ================= COMMON =================

type KeyVal struct {
	K uint16
	V uint64
}

func New(ar []uint64) []KeyVal {
	newar := make([]KeyVal, len(ar))
	for i, v := range ar {
		newar[i] = KeyVal{uint16(i), v}
	}
	return newar
}

func Fill(ar []uint64, newar []KeyVal) []KeyVal {
	if cap(newar) < len(ar) {
		newar = make([]KeyVal, len(ar))
	} else {
		newar = newar[0:len(ar)]
	}
	for i, v := range ar {
		newar[i] = KeyVal{uint16(i), v}
	}
	return newar
}

func Keys(ar []uint16, newar []KeyVal) []uint16 {
	l := len(newar)
	if len(ar) < l {
		ar = make([]uint16, l)
	}
	for i, v := range newar {
		ar[i] = v.K
	}
	return ar[0:l]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ------------- ASCENDING -------------

func heapSortAsc(data []KeyVal, a, b int) {
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

func insertionSortAsc(data []KeyVal, a, b int) {
	var j int
	for i := a + 1; i < b; i++ {
		for j = i; j > a && data[j].V < data[j-1].V; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDownAsc(data []KeyVal, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child].V < data[first+child+1].V {
			child++
		}
		if data[first+root].V >= data[first+child].V {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func medianOfThreeAsc(data []KeyVal, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if data[m1].V < data[m0].V {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data[m2].V < data[m1].V {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data[m1].V < data[m0].V {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRangeAsc(data []KeyVal, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a], data[b] = data[b], data[a]
		a++
		b++
	}
}

func doPivotAsc(data []KeyVal, lo, hi int) (midlo, midhi int) {
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
			if data[b].V < data[pivot].V {
				b++
			} else if data[pivot].V >= data[b].V {
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data[pivot].V < data[c-1].V {
				c--
			} else if data[c-1].V >= data[pivot].V {
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

func quickSortAsc(data []KeyVal, a, b, maxDepth int) {
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

func Asc(data []KeyVal) {
	maxDepth := 0
	for i := len(data); i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortAsc(data, 0, len(data), maxDepth)
}

func IsSortedAsc(data []KeyVal) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i].V < data[i-1].V {
			return false
		}
	}
	return true
}

func StableAsc(data []KeyVal) {
	n := len(data)
	blockSize := 20
	a, b := 0, blockSize
	for b <= n {
		insertionSortAsc(data, a, b)
		a = b
		b += blockSize
	}
	insertionSortAsc(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMergeAsc(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		symMergeAsc(data, a, a+blockSize, n)
		blockSize *= 2
	}
}

func symMergeAsc(data []KeyVal, a, m, b int) {
	if a >= m || m >= b {
		return
	}
	mid := a + (b-a)/2
	n := mid + m
	var start, c, r, p int
	if m > mid {
		start = n - b
		r, p = mid, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].V >= data[c].V {
				start = c + 1
			} else {
				r = c
			}
		}
	} else {
		start = a
		r, p = m, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].V >= data[c].V {
				start = c + 1
			} else {
				r = c
			}
		}
	}
	end := n - start
	rotateAsc(data, start, m, end)
	symMergeAsc(data, a, start, mid)
	symMergeAsc(data, mid, end, b)
}

func rotateAsc(data []KeyVal, a, m, b int) {
	i := m - a
	if i == 0 {
		return
	}
	j := b - m
	if j == 0 {
		return
	}
	if i == j {
		swapRangeAsc(data, a, m, i)
		return
	}
	p := a + i
	for i != j {
		if i > j {
			swapRangeAsc(data, p-i, p, j)
			i -= j
		} else {
			swapRangeAsc(data, p-i, p+j-i, i)
			j -= i
		}
	}
	swapRangeAsc(data, p-i, p, i)
}

// ------------- DESCENDING -------------

func heapSortDesc(data []KeyVal, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownDesc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownDesc(data, lo, i, first)
	}
}

func insertionSortDesc(data []KeyVal, a, b int) {
	var j int
	for i := a + 1; i < b; i++ {
		for j = i; j > a && data[j].V > data[j-1].V; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDownDesc(data []KeyVal, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child].V > data[first+child+1].V {
			child++
		}
		if data[first+root].V <= data[first+child].V {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func medianOfThreeDesc(data []KeyVal, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if data[m1].V > data[m0].V {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data[m2].V > data[m1].V {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data[m1].V > data[m0].V {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRangeDesc(data []KeyVal, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a], data[b] = data[b], data[a]
		a++
		b++
	}
}

func doPivotDesc(data []KeyVal, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeDesc(data, lo, lo+s, lo+2*s)
		medianOfThreeDesc(data, m, m-s, m+s)
		medianOfThreeDesc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeDesc(data, lo, m, hi-1)

	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data[b].V > data[pivot].V {
				b++
			} else if data[pivot].V <= data[b].V {
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data[pivot].V > data[c-1].V {
				c--
			} else if data[c-1].V <= data[pivot].V {
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
	swapRangeDesc(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeDesc(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortDesc(data []KeyVal, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortDesc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotDesc(data, a, b)
		if mlo-a < b-mhi {
			quickSortDesc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSortDesc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		insertionSortDesc(data, a, b)
	}
}

func Desc(data []KeyVal) {
	maxDepth := 0
	for i := len(data); i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortDesc(data, 0, len(data), maxDepth)
}

func IsSortedDesc(data []KeyVal) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i].V > data[i-1].V {
			return false
		}
	}
	return true
}

func StableDesc(data []KeyVal) {
	n := len(data)
	blockSize := 20
	a, b := 0, blockSize
	for b <= n {
		insertionSortDesc(data, a, b)
		a = b
		b += blockSize
	}
	insertionSortDesc(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMergeDesc(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		symMergeDesc(data, a, a+blockSize, n)
		blockSize *= 2
	}
}

func symMergeDesc(data []KeyVal, a, m, b int) {
	if a >= m || m >= b {
		return
	}
	mid := a + (b-a)/2
	n := mid + m
	var start, c, r, p int
	if m > mid {
		start = n - b
		r, p = mid, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].V < data[c].V {
				start = c + 1
			} else {
				r = c
			}
		}
	} else {
		start = a
		r, p = m, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].V < data[c].V {
				start = c + 1
			} else {
				r = c
			}
		}
	}
	end := n - start
	rotateDesc(data, start, m, end)
	symMergeDesc(data, a, start, mid)
	symMergeDesc(data, mid, end, b)
}

func rotateDesc(data []KeyVal, a, m, b int) {
	i := m - a
	if i == 0 {
		return
	}
	j := b - m
	if j == 0 {
		return
	}
	if i == j {
		swapRangeDesc(data, a, m, i)
		return
	}
	p := a + i
	for i != j {
		if i > j {
			swapRangeDesc(data, p-i, p, j)
			i -= j
		} else {
			swapRangeDesc(data, p-i, p+j-i, i)
			j -= i
		}
	}
	swapRangeDesc(data, p-i, p, i)
}

// ------------- ASCENDING BY KEY -------------

func heapSortByKeyAsc(data []KeyVal, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownByKeyAsc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownByKeyAsc(data, lo, i, first)
	}
}

func insertionSortByKeyAsc(data []KeyVal, a, b int) {
	var j int
	for i := a + 1; i < b; i++ {
		for j = i; j > a && data[j].K < data[j-1].K; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDownByKeyAsc(data []KeyVal, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child].K < data[first+child+1].K {
			child++
		}
		if data[first+root].K >= data[first+child].K {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func medianOfThreeByKeyAsc(data []KeyVal, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if data[m1].K < data[m0].K {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data[m2].K < data[m1].K {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data[m1].K < data[m0].K {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRangeByKeyAsc(data []KeyVal, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a], data[b] = data[b], data[a]
		a++
		b++
	}
}

func doPivotByKeyAsc(data []KeyVal, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeByKeyAsc(data, lo, lo+s, lo+2*s)
		medianOfThreeByKeyAsc(data, m, m-s, m+s)
		medianOfThreeByKeyAsc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeByKeyAsc(data, lo, m, hi-1)

	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data[b].K < data[pivot].K {
				b++
			} else if data[pivot].K >= data[b].K {
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data[pivot].K < data[c-1].K {
				c--
			} else if data[c-1].K >= data[pivot].K {
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
	swapRangeByKeyAsc(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeByKeyAsc(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortByKeyAsc(data []KeyVal, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortByKeyAsc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotByKeyAsc(data, a, b)
		if mlo-a < b-mhi {
			quickSortByKeyAsc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSortByKeyAsc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		insertionSortByKeyAsc(data, a, b)
	}
}

func AscByKey(data []KeyVal) {
	maxDepth := 0
	for i := len(data); i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortByKeyAsc(data, 0, len(data), maxDepth)
}

func IsSortedAscByKey(data []KeyVal) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i].K < data[i-1].K {
			return false
		}
	}
	return true
}

func StableAscByKey(data []KeyVal) {
	n := len(data)
	blockSize := 20
	a, b := 0, blockSize
	for b <= n {
		insertionSortByKeyAsc(data, a, b)
		a = b
		b += blockSize
	}
	insertionSortByKeyAsc(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMergeByKeyAsc(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		symMergeByKeyAsc(data, a, a+blockSize, n)
		blockSize *= 2
	}
}

func symMergeByKeyAsc(data []KeyVal, a, m, b int) {
	if a >= m || m >= b {
		return
	}
	mid := a + (b-a)/2
	n := mid + m
	var start, c, r, p int
	if m > mid {
		start = n - b
		r, p = mid, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].K >= data[c].K {
				start = c + 1
			} else {
				r = c
			}
		}
	} else {
		start = a
		r, p = m, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].K >= data[c].K {
				start = c + 1
			} else {
				r = c
			}
		}
	}
	end := n - start
	rotateByKeyAsc(data, start, m, end)
	symMergeByKeyAsc(data, a, start, mid)
	symMergeByKeyAsc(data, mid, end, b)
}

func rotateByKeyAsc(data []KeyVal, a, m, b int) {
	i := m - a
	if i == 0 {
		return
	}
	j := b - m
	if j == 0 {
		return
	}
	if i == j {
		swapRangeByKeyAsc(data, a, m, i)
		return
	}
	p := a + i
	for i != j {
		if i > j {
			swapRangeByKeyAsc(data, p-i, p, j)
			i -= j
		} else {
			swapRangeByKeyAsc(data, p-i, p+j-i, i)
			j -= i
		}
	}
	swapRangeByKeyAsc(data, p-i, p, i)
}

// ------------- DESCENDING BY KEY -------------

func heapSortByKeyDesc(data []KeyVal, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownByKeyDesc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownByKeyDesc(data, lo, i, first)
	}
}

func insertionSortByKeyDesc(data []KeyVal, a, b int) {
	var j int
	for i := a + 1; i < b; i++ {
		for j = i; j > a && data[j].K > data[j-1].K; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDownByKeyDesc(data []KeyVal, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child].K > data[first+child+1].K {
			child++
		}
		if data[first+root].K <= data[first+child].K {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func medianOfThreeByKeyDesc(data []KeyVal, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if data[m1].K > data[m0].K {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data[m2].K > data[m1].K {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data[m1].K > data[m0].K {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRangeByKeyDesc(data []KeyVal, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a], data[b] = data[b], data[a]
		a++
		b++
	}
}

func doPivotByKeyDesc(data []KeyVal, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeByKeyDesc(data, lo, lo+s, lo+2*s)
		medianOfThreeByKeyDesc(data, m, m-s, m+s)
		medianOfThreeByKeyDesc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeByKeyDesc(data, lo, m, hi-1)

	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data[b].K > data[pivot].K {
				b++
			} else if data[pivot].K <= data[b].K {
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data[pivot].K > data[c-1].K {
				c--
			} else if data[c-1].K <= data[pivot].K {
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
	swapRangeByKeyDesc(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeByKeyDesc(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortByKeyDesc(data []KeyVal, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortByKeyDesc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotByKeyDesc(data, a, b)
		if mlo-a < b-mhi {
			quickSortByKeyDesc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSortByKeyDesc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		insertionSortByKeyDesc(data, a, b)
	}
}

func DescByKey(data []KeyVal) {
	maxDepth := 0
	for i := len(data); i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortByKeyDesc(data, 0, len(data), maxDepth)
}

func IsSortedByKeyDesc(data []KeyVal) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i].K > data[i-1].K {
			return false
		}
	}
	return true
}

func StableDescByKey(data []KeyVal) {
	n := len(data)
	blockSize := 20
	a, b := 0, blockSize
	for b <= n {
		insertionSortByKeyDesc(data, a, b)
		a = b
		b += blockSize
	}
	insertionSortByKeyDesc(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMergeByKeyDesc(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		symMergeByKeyDesc(data, a, a+blockSize, n)
		blockSize *= 2
	}
}

func symMergeByKeyDesc(data []KeyVal, a, m, b int) {
	if a >= m || m >= b {
		return
	}
	mid := a + (b-a)/2
	n := mid + m
	var start, c, r, p int
	if m > mid {
		start = n - b
		r, p = mid, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].K < data[c].K {
				start = c + 1
			} else {
				r = c
			}
		}
	} else {
		start = a
		r, p = m, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c].K < data[c].K {
				start = c + 1
			} else {
				r = c
			}
		}
	}
	end := n - start
	rotateByKeyDesc(data, start, m, end)
	symMergeByKeyDesc(data, a, start, mid)
	symMergeByKeyDesc(data, mid, end, b)
}

func rotateByKeyDesc(data []KeyVal, a, m, b int) {
	i := m - a
	if i == 0 {
		return
	}
	j := b - m
	if j == 0 {
		return
	}
	if i == j {
		swapRangeByKeyDesc(data, a, m, i)
		return
	}
	p := a + i
	for i != j {
		if i > j {
			swapRangeByKeyDesc(data, p-i, p, j)
			i -= j
		} else {
			swapRangeByKeyDesc(data, p-i, p+j-i, i)
			j -= i
		}
	}
	swapRangeByKeyDesc(data, p-i, p, i)
}
