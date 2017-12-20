## Sort

Go's native Sort package uses interface elements, the reflection on which considerably slows down the sorting algorithms. Additionally descending order sorting using greater than for the `Less()` function is inefficient.

Herein are ascending and descending implemetations for native number types, with no reflection, less function calls and other optimizations.

Included are also key/value sorting algorithms and stable sorting.

### Example an ascending sort on slice of Ints

     import "github.com/AlasdairF/Sort/Int"
     list := []int{10, 44, 1, 7, 4, 0, -9, 0, 3, 65, 38}
     sortInt.Asc(list)
     
### Example a descending stable sort on slice of Uint32s

     import "github.com/AlasdairF/Sort/Uint32"
     list := []uint32{10, 44, 1, 7, 4, 0, 9, 0, 3, 65, 38}
     sortUint32.StableDesc(list)

### Example sort on key/value pair Uint32/Float64

     import "github.com/AlasdairF/Sort/Uint32Float64"
     keyval := []sortUint32Float64.KeyVal{
         sortUint32Float64.KeyVal{0, 10.5},
         sortUint32Float64.KeyVal{1, 44.1},
         sortUint32Float64.KeyVal{2, 1.9},
         sortUint32Float64.KeyVal{3, 8.5},
     }
     sortUint32Float64.Desc(keyval)

### Example identical to above using the `New()` helper function

     import "github.com/AlasdairF/Sort/Uint32Float64"
     scores := []float64{10.5, 44.1, 1.9, 8.5}
     keyval := sortUint32Float64.New(scores) // keys are filled in automatically starting from 0
     sortUint32Float64.Desc(keyval)
     
     // Then the keyval underlying array can be reused as follows
     scores2 := []float64{5.4, 30.4, 100.5}
     keyval = sortUint32Float64.Fill(scores2, keyval)
     sortUint32Float64.Asc(keyval)
     
     // Maybe after the sorting you only want the indexes?
     keys := sortUint32Float64.Keys([]uint32{}, keyval)

## Proof

     package main
     
     import (
      "time"
      "fmt"
      "sort"
      "github.com/AlasdairF/Sort/Int"
      "math/rand"
     )
     
     func newCopy(ints []int) []int {
     	s := make([]int, len(ints))
     	copy(s, ints)
     	return s
     }
     
     func main() {
     
     // 2x 100 random slices of ints (all the same)
     ints := make([]int, 100000)
     for i:=0; i<100000; i++ {
     	ints[i] = rand.Int()
     }
     copies1 := make([][]int, 100)
     copies2 := make([][]int, 100)
     for i:=0; i<100; i++ {
     	copies1[i] = newCopy(ints)
     	copies2[i] = newCopy(ints)
     }
     
     t1:= time.Now().UnixNano()
     for i:=0; i<100; i++ {
     	sort.Ints(copies1[i])
     }
     t2 := time.Now().UnixNano()
     for i:=0; i<100; i++ {
     	sortInt.Asc(copies2[i])
     }
     t3 := time.Now().UnixNano()
     
     fmt.Println(`native sort took`, t2-t1, `nanoseconds`)
     fmt.Println(`sortInt took`, t3-t2, `nanoseconds`)
     }

The results:     
     
     root /home/root # ./test
     native sort took 4170648735 nanoseconds
     sortInt took     1464291506 nanoseconds
     root /home/root # ./test
     native sort took 3938108749 nanoseconds
     sortInt took     1286878885 nanoseconds
     root /home/root # ./test
     native sort took 4251135285 nanoseconds
     sortInt took     1610303693 nanoseconds
     root /home/root # ./test
     native sort took 4721015337 nanoseconds
     sortInt took     1380711512 nanoseconds
     root /home/root # ./test
     native sort took 4785185202 nanoseconds
     sortInt took     1572508937 nanoseconds
     root /home/root # ./test
     native sort took 3235662481 nanoseconds
     sortInt took     1329911802 nanoseconds
     root /home/root # ./test
     native sort took 3523107884 nanoseconds
     sortInt took     1246381327 nanoseconds
     
