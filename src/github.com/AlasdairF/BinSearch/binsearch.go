package binsearch

import (
// Customized sorting algorithms for KeyBytes
 "github.com/AlasdairF/BinSearch/Limit16"
 "github.com/AlasdairF/BinSearch/Limit24"
 "github.com/AlasdairF/BinSearch/Limit32"
 "github.com/AlasdairF/BinSearch/Limit40"
 "github.com/AlasdairF/BinSearch/Limit48"
 "github.com/AlasdairF/BinSearch/Limit56"
 "github.com/AlasdairF/BinSearch/Limit64"
// Customized sorting algorithms for KeyValBytes
 "github.com/AlasdairF/BinSearch/LimitVal8"
 "github.com/AlasdairF/BinSearch/LimitVal16"
 "github.com/AlasdairF/BinSearch/LimitVal24"
 "github.com/AlasdairF/BinSearch/LimitVal32"
 "github.com/AlasdairF/BinSearch/LimitVal40"
 "github.com/AlasdairF/BinSearch/LimitVal48"
 "github.com/AlasdairF/BinSearch/LimitVal56"
 "github.com/AlasdairF/BinSearch/LimitVal64"
// General optimized sorting algorithms using primitives
 "github.com/AlasdairF/Sort/IntUint64"
 "github.com/AlasdairF/Sort/IntUint32"
 "github.com/AlasdairF/Sort/IntUint16"
 "github.com/AlasdairF/Sort/IntUint8"
 "github.com/AlasdairF/Sort/IntInt"
// Read/write custom file format
 "github.com/AlasdairF/Custom"
// Error handling
 "errors"
)

/*

	INFORMATION
	
	Maximum key size is 64 bytes.

	TYPES
	
	'Key' type (KeyBytes, KeyUint64 etc.) is a key-only store where the index of each key can be used to store any number of associated values.
	'KeyVal' type (KeyValBytes, KeyValUint64, etc.) is a key-value store where the vals are ints. For custom value types use 'Key' (above).
	'Counter' type (CounterBytes, CounterUint64, etc.) counts the number of occurances (equivalent to map[key]++) and allows very fast lookups.

	INDEX
	
	KeyBytes, KeyRunes
		func (t *KeyBytes) Len() int
		func (t *KeyBytes) Find(thekey []byte) (int, bool)					Returns: index, exists.
		func (t *KeyBytes) Add(thekey []byte) (int, bool)					Returns: index, exists. Adds the key if it does not already exist and returns the new index, otherwise returns the current index of the existing key.
		func (t *KeyBytes) AddAt(thekey []byte, i int) error				Returns error if thekey > 64 bytes
		func (t *KeyBytes) AddUnsorted(thekey []byte) error					Returns error if thekey > 64 bytes
		func (t *KeyBytes) Build() ([]int, error)							Returns slice mapping old indexes to new indexes. Can only be used after AddUnsorted, otherwise returns an error.
		func (t *KeyBytes) Optimize()										Copies all the data to new slices with capacity equal to length.
		func (t *KeyBytes) Reset() bool										Returns false if the structure is empty (Len() == 0)
		func (t *KeyBytes) Next() ([]byte, bool)							Returns: original slice of bytes, EOF (true = EOF)
		func (t *KeyBytes) Keys() [][]byte									Returns slice containing all the keys in order
		func (t *KeyBytes) Write(w *custom.Writer)							Writes built structure out to custom.Writer (requires github.com/AlasdairF/Custom)
		func (t *KeyBytes) Read(r *custom.Reader)							Reads structure in from custom.Reader (requires github.com/AlasdairF/Custom)
		
	KeyValBytes, KeyValRunes
		func (t *KeyValBytes) Len() int
		func (t *KeyValBytes) Find(thekey []byte) (int, bool)				Returns: value, exists
		func (t *KeyValBytes) Update(thekey []byte, fn func(int) int) bool	Returns boolean value for whether the key exists or not, if it exists the value is modified according to the fn function
		func (t *KeyValBytes) UpdateAll(fn func(int) int)					Modifies all values by the fn function
		func (t *KeyValBytes) Add(thekey []byte, theval int) bool			Returns whether it exists. Replaces old value with the new value if it exists, otherwise adds it in place.
		func (t *KeyValBytes) AddUnsorted(thekey []byte, theval int) error	Returns error if thekey > 64 bytes
		func (t *KeyValBytes) Build()										Only required to be called after AddUnsorted, otherwise it will shrink array capacity to length.
		func (t *KeyValBytes) Optimize()									Copies all the data to new slices with capacity equal to length.
		func (t *KeyValBytes) Reset() bool									Returns false if the structure is empty (Len() == 0)
		func (t *KeyValBytes) Next() ([]byte, int, bool)					Returns: original slice of bytes, value, EOF (true = EOF)
		func (t *KeyValBytes) Keys() [][]byte								Returns slice containing all the keys in order
		func (t *KeyValBytes) Write(w *custom.Writer)						Writes built structure out to custom.Writer (requires github.com/AlasdairF/Custom)
		func (t *KeyValBytes) Read(r *custom.Reader)						Reads structure in from custom.Reader (requires github.com/AlasdairF/Custom)
		
	CounterBytes, CounterRunes
		func (t *CounterBytes) Len() int									Len() is only accurate after Build()
		func (t *CounterBytes) Find(thekey []byte) (int, bool)				Returns: frequency, exists. Will return nonsensical results if used before Build() is executed; only use after Build.
		func (t *CounterBytes) Update(thekey []byte, fn func(int) int) bool	Returns boolean value for whether the key exists or not, if it exists the value is modified according to the fn function
		func (t *CounterBytes) UpdateAll(fn func(int) int)					Modifies all values by the fn function
		func (t *CounterBytes) Add(thekey []byte, theval int) error			Returns an error if thekey > 64 bytes
		func (t *CounterBytes) Build()										Always required before Find.
		func (t *CounterBytes) Optimize()									Copies all the data to new slices with capacity equal to length.
		func (t *CounterBytes) Reset() bool									Returns false if the structure is empty (Len() == 0)
		func (t *CounterBytes) Next() ([]byte, int, bool)					Returns: original slice of bytes, value, EOF (true = EOF)
		func (t *CounterBytes) Keys() [][]byte								Returns slice containing all the keys in order
		func (t *CounterBytes) Write(w *custom.Writer)						Writes built structure out to custom.Writer (requires github.com/AlasdairF/Custom)
		func (t *CounterBytes) Read(r *custom.Reader)						Reads structure in from custom.Reader (requires github.com/AlasdairF/Custom)
		func (t *CounterBytes) KeyBytes() *KeyBytes							Copies keys to a KeyBytes structure
		func (t *CounterBytes) KeyValBytes() *KeyBytes						Copies keys and values to a KeyValBytes structure
		
	KeyInt, KeyUint64, KeyUint32, KeyUint16, KeyUint8
		func (t *KeyInt) Len() int
		func (t *KeyInt) Find(thekey []byte) (int, bool)					Returns: index, exists.
		func (t *KeyInt) Add(thekey []byte) (int, bool)						Returns: index, exists.
		func (t *KeyInt) AddAt(thekey []byte, i int)
		func (t *KeyInt) AddUnsorted(thekey []byte)
		func (t *KeyInt) Build() []int										Returns slice mapping old indexes to new indexes. Only required if AddUnsorted was used, otherwise it will shrink array capacity to length.
		func (t *KeyInt) Optimize()											Copies all the data to new slices with capacity equal to length.
		func (t *KeyInt) Reset() bool										Returns false if the structure is empty (Len() == 0)
		func (t *KeyInt) Next() (uint64, bool)								Returns: key, EOF (true = EOF)
		func (t *KeyInt) Keys() []uint64									Returns slice containing all the keys in order
		func (t *KeyInt) Write(w *custom.Writer)							Writes built structure out to custom.Writer (requires github.com/AlasdairF/Custom)
		func (t *KeyInt) Read(r *custom.Reader)								Reads structure in from custom.Reader (requires github.com/AlasdairF/Custom)
		
	KeyValInt, KeyValUint64, KeyValUint32, KeyValUint16, KeyValUint8
		func (t *KeyValInt) Len() int
		func (t *KeyValInt) Find(thekey uint64) (int, bool)					Returns: value, exists
		func (t *KeyValInt) Update(thekey uint64, fn func(int) int) bool	Returns boolean value for whether the key exists or not, if it exists the value is modified according to the fn function
		func (t *KeyValInt) UpdateAll(fn func(int) int)						Modifies all values by the fn function
		func (t *KeyValInt) Add(thekey uint64, theval int) bool				Returns whether it exists. Replaces old value with the new value if it exists, otherwise adds it in place.
		func (t *KeyValInt) AddUnsorted(thekey uint64, theval int)
		func (t *KeyValInt) Build()											Only required to be called after AddUnsorted, otherwise it will shrink array capacity to length.
		func (t *KeyValInt) Optimize()										Copies all the data to new slices with capacity equal to length.
		func (t *KeyValInt) Reset() bool									Returns false if the structure is empty (Len() == 0)
		func (t *KeyValInt) Next() ([]byte, int, bool)						Returns: original slice of bytes, value, EOF (true = EOF)
		func (t *KeyValInt) Keys() []uint64									Returns slice containing all the keys in order
		func (t *KeyValInt) Write(w *custom.Writer)							Writes built structure out to custom.Writer (requires github.com/AlasdairF/Custom)
		func (t *KeyValInt) Read(r *custom.Reader)							Reads structure in from custom.Reader (requires github.com/AlasdairF/Custom))
		
	CounterInt, CounterUint64, CounterUint32, CounterUint16, CounterUint8
		func (t *CounterInt) Len() int										Len() is only accurate after Build()
		func (t *CounterInt) Find(thekey uint64) (int, bool)				Returns: frequency, exists. Will return nonsensical results if used before Build() is executed; only use after Build.
		func (t *CounterInt) Update(thekey uint64, fn func(int) int) bool	Returns boolean value for whether the key exists or not, if it exists the value is modified according to the fn function
		func (t *CounterInt) UpdateAll(fn func(int) int)					Modifies all values by the fn function
		func (t *CounterInt) Add(thekey uint64, theval int)
		func (t *CounterInt) Build()										Always required before Find.
		func (t *CounterInt) Optimize()										Copies all the data to new slices with capacity equal to length.
		func (t *CounterInt) Reset() bool									Returns false if the structure is empty (Len() == 0)
		func (t *CounterInt) Next() ([]byte, int, bool)						Returns: original slice of bytes, value, EOF (true = EOF)
		func (t *CounterInt) Keys() []uint64								Returns slice containing all the keys in order
		func (t *CounterInt) Write(w *custom.Writer)						Writes built structure out to custom.Writer (requires github.com/AlasdairF/Custom)
		func (t *CounterInt) Read(r *custom.Reader)							Reads structure in from custom.Reader (requires github.com/AlasdairF/Custom)
		func (t *CounterInt) Copy() *KeyInt									Copies keys to a KeyInt structure

*/

// ====================== bytes ======================

// ---------- KeyBytes ----------
// Key bytes has around 5KB of memory overhead for the structures, beyond this it stores all keys as efficiently as possible.

// Add this to any struct to make it binary searchable.
type KeyBytes struct {
 limit8 [8][]uint64 // where len(word) <= 8
 limit16 [8][][2]uint64
 limit24 [8][][3]uint64
 limit32 [8][][4]uint64
 limit40 [8][][5]uint64
 limit48 [8][][6]uint64
 limit56 [8][][7]uint64
 limit64 [8][][8]uint64
// The order vars are used only when using AddSorted & Build. Build clears them. They are used for remembering the order that the keys were added in so the remap can be returned to the user by Build.
 order8 [8][]int
 order16 [8][]int
 order24 [8][]int
 order32 [8][]int
 order40 [8][]int
 order48 [8][]int
 order56 [8][]int
 order64 [8][]int
 count [64]int // Used to convert limit maps to the 1D array value indicating where the value exists
 total int
// Used for iterating through all of it
 onlimit int
 on8 int
 oncursor int
}

func bytes2uint64(word []byte) (uint64, int) {
	switch len(word) {
		case 0:
			return 0, 0 // an empty slice is sorted with the single characters
		case 1:
			return uint64(word[0]), 0
		case 2:
			return (uint64(word[0]) << 8) | uint64(word[1]), 1
		case 3:
			return (uint64(word[0]) << 16) | (uint64(word[1]) << 8) | uint64(word[2]), 2
		case 4:
			return (uint64(word[0]) << 24) | (uint64(word[1]) << 16) | (uint64(word[2]) << 8) | uint64(word[3]), 3
		case 5:
			return (uint64(word[0]) << 32) | (uint64(word[1]) << 24) | (uint64(word[2]) << 16) | (uint64(word[3]) << 8) | uint64(word[4]), 4
		case 6:
			return (uint64(word[0]) << 40) | (uint64(word[1]) << 32) | (uint64(word[2]) << 24) | (uint64(word[3]) << 16) | (uint64(word[4]) << 8) | uint64(word[5]), 5
		case 7:
			return (uint64(word[0]) << 48) | (uint64(word[1]) << 40) | (uint64(word[2]) << 32) | (uint64(word[3]) << 24) | (uint64(word[4]) << 16) | (uint64(word[5]) << 8) | uint64(word[6]), 6
		default:
			return (uint64(word[0]) << 56) | (uint64(word[1]) << 48) | (uint64(word[2]) << 40) | (uint64(word[3]) << 32) | (uint64(word[4]) << 24) | (uint64(word[5]) << 16) | (uint64(word[6]) << 8) | uint64(word[7]), 7
	}
}

func uint642bytes(word []byte, v uint64) {
	word[7] = byte(v & 255)
	word[6] = byte((v >> 8) & 255)
	word[5] = byte((v >> 16) & 255)
	word[4] = byte((v >> 24) & 255)
	word[3] = byte((v >> 32) & 255)
	word[2] = byte((v >> 40) & 255)
	word[1] = byte((v >> 48) & 255)
	word[0] = byte((v >> 56) & 255)
}


func uint642bytesend(word []byte, v uint64) int {
	if v >> 56 != 0 {
		word[7] = byte(v & 255)
		word[6] = byte((v >> 8) & 255)
		word[5] = byte((v >> 16) & 255)
		word[4] = byte((v >> 24) & 255)
		word[3] = byte((v >> 32) & 255)
		word[2] = byte((v >> 40) & 255)
		word[1] = byte((v >> 48) & 255)
		word[0] = byte((v >> 56) & 255)
		return 8
	}
	if v >> 48 != 0 {
		word[6] = byte(v & 255)
		word[5] = byte((v >> 8) & 255)
		word[4] = byte((v >> 16) & 255)
		word[3] = byte((v >> 24) & 255)
		word[2] = byte((v >> 32) & 255)
		word[1] = byte((v >> 40) & 255)
		word[0] = byte((v >> 48) & 255)
		return 7
	}
	if v >> 40 != 0 {
		word[5] = byte(v & 255)
		word[4] = byte((v >> 8) & 255)
		word[3] = byte((v >> 16) & 255)
		word[2] = byte((v >> 24) & 255)
		word[1] = byte((v >> 32) & 255)
		word[0] = byte((v >> 40) & 255)
		return 6
	}
	if v >> 32 != 0 {
		word[4] = byte(v & 255)
		word[3] = byte((v >> 8) & 255)
		word[2] = byte((v >> 16) & 255)
		word[1] = byte((v >> 24) & 255)
		word[0] = byte((v >> 32) & 255)
		return 5
	}
	if v >> 24 != 0 {
		word[3] = byte(v & 255)
		word[2] = byte((v >> 8) & 255)
		word[1] = byte((v >> 16) & 255)
		word[0] = byte((v >> 24) & 255)
		return 4
	}
	if v >> 16 != 0 {
		word[2] = byte(v & 255)
		word[1] = byte((v >> 8) & 255)
		word[0] = byte((v >> 16) & 255)
		return 3
	}
	if v >> 8 != 0 {
		word[1] = byte(v & 255)
		word[0] = byte((v >> 8) & 255)
		return 2
	}
	if v != 0 {
		word[0] = byte(v & 255)
		return 1
	}
	return 0
}

func reverse8(v uint64) []byte {
	word := make([]byte, 8)
	i := uint642bytesend(word, v)
	return word[0:i]
}

func reverse8b(v [2]uint64) []byte {
	word := make([]byte, 8)
	i := uint642bytesend(word, v[0])
	return word[0:i]
}

func reverse16(v [2]uint64) []byte {
	word := make([]byte, 16)
	uint642bytes(word, v[0])
	i := uint642bytesend(word[8:], v[1])
	return word[0 : 8 + i]
}

func reverse16b(v [3]uint64) []byte {
	word := make([]byte, 16)
	uint642bytes(word, v[0])
	i := uint642bytesend(word[8:], v[1])
	return word[0 : 8 + i]
}

func reverse24(v [3]uint64) []byte {
	word := make([]byte, 24)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	i := uint642bytesend(word[16:], v[2])
	return word[0 : 16 + i]
}

func reverse24b(v [4]uint64) []byte {
	word := make([]byte, 24)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	i := uint642bytesend(word[16:], v[2])
	return word[0 : 16 + i]
}

func reverse32(v [4]uint64) []byte {
	word := make([]byte, 32)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	i := uint642bytesend(word[24:], v[3])
	return word[0 : 24 + i]
}

func reverse32b(v [5]uint64) []byte {
	word := make([]byte, 32)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	i := uint642bytesend(word[24:], v[3])
	return word[0 : 24 + i]
}

func reverse40(v [5]uint64) []byte {
	word := make([]byte, 40)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	i := uint642bytesend(word[32:], v[4])
	return word[0 : 32 + i]
}

func reverse40b(v [6]uint64) []byte {
	word := make([]byte, 40)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	i := uint642bytesend(word[32:], v[4])
	return word[0 : 32 + i]
}

func reverse48(v [6]uint64) []byte {
	word := make([]byte, 48)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	uint642bytes(word[32:], v[4])
	i := uint642bytesend(word[40:], v[5])
	return word[0 : 40 + i]
}

func reverse48b(v [7]uint64) []byte {
	word := make([]byte, 48)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	uint642bytes(word[32:], v[4])
	i := uint642bytesend(word[40:], v[5])
	return word[0 : 40 + i]
}

func reverse56(v [7]uint64) []byte {
	word := make([]byte, 56)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	uint642bytes(word[32:], v[4])
	uint642bytes(word[40:], v[5])
	i := uint642bytesend(word[48:], v[6])
	return word[0 : 48 + i]
}

func reverse56b(v [8]uint64) []byte {
	word := make([]byte, 56)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	uint642bytes(word[32:], v[4])
	uint642bytes(word[40:], v[5])
	i := uint642bytesend(word[48:], v[6])
	return word[0 : 48 + i]
}

func reverse64(v [8]uint64) []byte {
	word := make([]byte, 64)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	uint642bytes(word[32:], v[4])
	uint642bytes(word[40:], v[5])
	uint642bytes(word[48:], v[6])
	i := uint642bytesend(word[56:], v[7])
	return word[0 : 56 + i]
}

func reverse64b(v [9]uint64) []byte {
	word := make([]byte, 64)
	uint642bytes(word, v[0])
	uint642bytes(word[8:], v[1])
	uint642bytes(word[16:], v[2])
	uint642bytes(word[24:], v[3])
	uint642bytes(word[32:], v[4])
	uint642bytes(word[40:], v[5])
	uint642bytes(word[48:], v[6])
	i := uint642bytesend(word[56:], v[7])
	return word[0 : 56 + i]
}

func (t *KeyBytes) Len() int {
	return t.total
}

// Find returns the index based on the key.
func (t *KeyBytes) Find(thekey []byte) (int, bool) {
	
	var at, min int
	var compare uint64
	switch (len(thekey) - 1) / 8 {
	
		case 0: // 0 - 8 bytes
			a, l := bytes2uint64(thekey)
			cur := t.limit8[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				return at + t.count[l], true // found
			}
			return min + t.count[l], false // doesn't exist
			
		case 1: // 9 - 16 bytes
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			cur := t.limit16[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 8], true // found
			}
			return min + t.count[l + 8], false // doesn't exist
			
		case 2: // 17 - 24 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			cur := t.limit24[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 16], true // found
			}
			return min + t.count[l + 16], false // doesn't exist
			
		case 3: // 25 - 32 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			cur := t.limit32[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 24], true // found
			}
			return min + t.count[l + 24], false // doesn't exist
			
		case 4: // 33 - 40 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			cur := t.limit40[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 32], true // found
			}
			return min + t.count[l + 32], false // doesn't exist
			
		case 5: // 41 - 48 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			cur := t.limit48[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 40], true // found
			}
			return min + t.count[l + 40], false // doesn't exist
			
		case 6: // 49 - 56 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			cur := t.limit56[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 48], true // found
			}
			return min + t.count[l + 48], false // doesn't exist
			
		case 7: // 57 - 64 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			cur := t.limit64[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][7]; h < compare {
					max = at - 1
					continue
				}
				if h > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 56], true // found
			}
			return min + t.count[l + 56], false // doesn't exist
		
		default: // > 64 bytes
			return t.total, false
	}
}

// Add is equivalent to Find and then AddAt
func (t *KeyBytes) Add(thekey []byte) (int, bool) {
	
	var at, min int
	var compare uint64
	switch (len(thekey) - 1) / 8 {
	
		case 0: // 0 - 8 bytes
			a, l := bytes2uint64(thekey)
			cur := t.limit8[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				return at + t.count[l], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = a
			t.limit8[l] = cur
			for l++; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
			
		case 1: // 9 - 16 bytes
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			cur := t.limit16[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 8], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l + 8]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][2]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [2]uint64{a, b}
			t.limit16[l] = cur
			for l+=9; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
			
		case 2: // 17 - 24 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			cur := t.limit24[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 16], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l + 16]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][3]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [3]uint64{a, b, c}
			t.limit24[l] = cur
			for l+=17; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
			
		case 3: // 25 - 32 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			cur := t.limit32[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 24], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l + 24]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][4]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [4]uint64{a, b, c, d}
			t.limit32[l] = cur
			for l+=25; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
			
		case 4: // 33 - 40 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			cur := t.limit40[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 32], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l + 32]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][5]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [5]uint64{a, b, c, d, e}
			t.limit40[l] = cur
			for l+=33; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
			
		case 5: // 41 - 48 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			cur := t.limit48[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 40], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l + 40]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][6]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [6]uint64{a, b, c, d, e, f}
			t.limit48[l] = cur
			for l+=41; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
			
		case 6: // 49 - 56 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			cur := t.limit56[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 48], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l + 48]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][7]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [7]uint64{a, b, c, d, e, f, g}
			t.limit56[l] = cur
			for l+=49; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
			
		case 7: // 57 - 64 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			cur := t.limit64[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][7]; h < compare {
					max = at - 1
					continue
				}
				if h > compare {
					min = at + 1
					continue
				}
				return at + t.count[l + 56], true // found
			}
			// Doesn't exist so add it >
			at = min
			min += t.count[l + 56]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][8]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [8]uint64{a, b, c, d, e, f, g, h}
			t.limit64[l] = cur
			for l+=57; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return min, false
		
		default: // > 64 bytes
			return t.total, false
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyBytes) AddUnsorted(thekey []byte) error {
	switch (len(thekey) - 1) / 8 {
		case 0:
			a, i := bytes2uint64(thekey)
			t.limit8[i] = append(t.limit8[i], a)
			t.order8[i] = append(t.order8[i], t.total)
			t.count[i + 1]++
			t.total++
			return nil
		case 1:
			a, _ := bytes2uint64(thekey)
			b, i := bytes2uint64(thekey[8:])
			t.limit16[i] = append(t.limit16[i], [2]uint64{a, b})
			t.order16[i] = append(t.order16[i], t.total)
			t.count[i + 9]++
			t.total++
			return nil
		case 2:
			
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, i := bytes2uint64(thekey[16:])
			t.limit24[i] = append(t.limit24[i], [3]uint64{a, b, c})
			t.order24[i] = append(t.order24[i], t.total)
			t.count[i + 17]++
			t.total++
			return nil
		case 3:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, i := bytes2uint64(thekey[24:])
			t.limit32[i] = append(t.limit32[i], [4]uint64{a, b, c, d})
			t.order32[i] = append(t.order32[i], t.total)
			t.count[i + 25]++
			t.total++
			return nil
		case 4:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, i := bytes2uint64(thekey[32:])
			t.limit40[i] = append(t.limit40[i], [5]uint64{a, b, c, d, e})
			t.order40[i] = append(t.order40[i], t.total)
			t.count[i + 33]++
			t.total++
			return nil
		case 5:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, i := bytes2uint64(thekey[40:])
			t.limit48[i] = append(t.limit48[i], [6]uint64{a, b, c, d, e, f})
			t.order48[i] = append(t.order48[i], t.total)
			t.count[i + 41]++
			t.total++
			return nil
		case 6:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, i := bytes2uint64(thekey[48:])
			t.limit56[i] = append(t.limit56[i], [7]uint64{a, b, c, d, e, f, g})
			t.order56[i] = append(t.order56[i], t.total)
			t.count[i + 49]++
			t.total++
			return nil
		case 7:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, i := bytes2uint64(thekey[56:])
			t.limit64[i] = append(t.limit64[i], [8]uint64{a, b, c, d, e, f, g, h})
			t.order64[i] = append(t.order64[i], t.total)
			if i < 7 {
				t.count[i + 57]++
			}
			t.total++
			return nil
		default:
			return errors.New(`Maximum key length is 64 bytes`)
	}
}

// AddAt adds this key to the index in this exact position, so it does not require later rebuilding.
func (t *KeyBytes) AddAt(thekey []byte, i int) error {

	switch (len(thekey) - 1) / 8 {
		case 0:
			a, l := bytes2uint64(thekey)
			i -= t.count[l]
			cur := t.limit8[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = a
			t.limit8[l] = cur
			for l++; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		case 1:
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			i -= t.count[l + 8]
			cur := t.limit16[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][2]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = [2]uint64{a, b}
			t.limit16[l] = cur
			for l+=9; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		case 2:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			i -= t.count[l + 16]
			cur := t.limit24[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][3]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = [3]uint64{a, b, c}
			t.limit24[l] = cur
			for l+=17; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		case 3:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			i -= t.count[l + 24]
			cur := t.limit32[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][4]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = [4]uint64{a, b, c, d}
			t.limit32[l] = cur
			for l+=25; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		case 4:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			i -= t.count[l + 32]
			cur := t.limit40[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][5]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = [5]uint64{a, b, c, d, e}
			t.limit40[l] = cur
			for l+=33; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		case 5:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			i -= t.count[l + 40]
			cur := t.limit48[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][6]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = [6]uint64{a, b, c, d, e, f}
			t.limit48[l] = cur
			for l+=41; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		case 6:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			i -= t.count[l + 48]
			cur := t.limit56[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][7]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = [7]uint64{a, b, c, d, e, f, g}
			t.limit56[l] = cur
			for l+=49; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		case 7:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			i -= t.count[l + 56]
			cur := t.limit64[l]
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][8]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:i])
				copy(tmp[i+1:], cur[i:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[i+1:], cur[i:])
			}
			cur[i] = [8]uint64{a, b, c, d, e, f, g, h}
			t.limit64[l] = cur
			for l+=57; l<64; l++ {
				t.count[l]++
			}
			t.total++
			return nil
			
		default:
			return errors.New(`Maximum key length is 64 bytes`)
	}
}

// Build sorts the keys and returns an array telling you how to sort the values, you must do this yourself.
func (t *KeyBytes) Build() ([]int, error) {

	var l, on, run int
	imap := make([]int, t.total)
	
	{
	var temp []sortIntUint64.KeyVal
	for run=0; run<8; run++ {
		if l = len(t.limit8[run]); l > 0 {
			m := t.order8[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make([]sortIntUint64.KeyVal, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit8[run] {
				temp[z] = sortIntUint64.KeyVal{m[z], k}
			}
			t.order8[run] = nil
			sortIntUint64.Asc(temp)
			newkey := t.limit8[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	{
	var temp sortLimit16.Slice
	for run=0; run<8; run++ {
		if l = len(t.limit16[run]); l > 0 {
			m := t.order16[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make(sortLimit16.Slice, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit16[run] {
				temp[z] = sortLimit16.KeyVal{m[z], k}
			}
			t.order16[run] = nil
			sortLimit16.Asc(temp)
			newkey := t.limit16[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	{
	var temp sortLimit24.Slice
	for run=0; run<8; run++ {
		if l = len(t.limit24[run]); l > 0 {
			m := t.order24[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make(sortLimit24.Slice, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit24[run] {
				temp[z] = sortLimit24.KeyVal{m[z], k}
			}
			t.order24[run] = nil
			sortLimit24.Asc(temp)
			newkey := t.limit24[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	{
	var temp sortLimit32.Slice
	for run=0; run<8; run++ {
		if l = len(t.limit32[run]); l > 0 {
			m := t.order32[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make(sortLimit32.Slice, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit32[run] {
				temp[z] = sortLimit32.KeyVal{m[z], k}
			}
			t.order32[run] = nil
			sortLimit32.Asc(temp)
			newkey := t.limit32[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	{
	var temp sortLimit40.Slice
	for run=0; run<8; run++ {
		if l = len(t.limit40[run]); l > 0 {
			m := t.order40[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make(sortLimit40.Slice, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit40[run] {
				temp[z] = sortLimit40.KeyVal{m[z], k}
			}
			t.order40[run] = nil
			sortLimit40.Asc(temp)
			newkey := t.limit40[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	{
	var temp sortLimit48.Slice
	for run=0; run<8; run++ {
		if l = len(t.limit48[run]); l > 0 {
			m := t.order48[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make(sortLimit48.Slice, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit48[run] {
				temp[z] = sortLimit48.KeyVal{m[z], k}
			}
			t.order48[run] = nil
			sortLimit48.Asc(temp)
			newkey := t.limit48[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	{
	var temp sortLimit56.Slice
	for run=0; run<8; run++ {
		if l = len(t.limit56[run]); l > 0 {
			m := t.order56[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make(sortLimit56.Slice, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit56[run] {
				temp[z] = sortLimit56.KeyVal{m[z], k}
			}
			t.order56[run] = nil
			sortLimit56.Asc(temp)
			newkey := t.limit56[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	{
	var temp sortLimit64.Slice
	for run=0; run<8; run++ {
		if l = len(t.limit64[run]); l > 0 {
			m := t.order64[run]
			if l != len(m) {
				return nil, errors.New(`Build can only be run once. After the first time use AddAt.`)
			}
			if cap(temp) < l {
				temp = make(sortLimit64.Slice, l)
			} else {
				temp = temp[0:l]
			}
			for z, k := range t.limit64[run] {
				temp[z] = sortLimit64.KeyVal{m[z], k}
			}
			t.order64[run] = nil
			sortLimit64.Asc(temp)
			newkey := t.limit64[run]
			for i, obj := range temp {
				imap[on] = obj.K
				on++
				newkey[i] = obj.V
			}
		}
	}
	}
	
	// Correct all the counts
	for run=2; run<64; run++ {
		t.count[run] += t.count[run-1]
	}
	
	return imap, nil
}

func (t *KeyBytes) Optimize() {

	var l, run int
	
	for run=0; run<8; run++ {
		if l = len(t.limit8[run]); l > 0 {
			newkey := make([]uint64, l)
			copy(newkey, t.limit8[run])
			t.limit8[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit16[run]); l > 0 {
			newkey := make([][2]uint64, l)
			copy(newkey, t.limit16[run])
			t.limit16[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit24[run]); l > 0 {
			newkey := make([][3]uint64, l)
			copy(newkey, t.limit24[run])
			t.limit24[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit32[run]); l > 0 {
			newkey := make([][4]uint64, l)
			copy(newkey, t.limit32[run])
			t.limit32[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit40[run]); l > 0 {
			newkey := make([][5]uint64, l)
			copy(newkey, t.limit40[run])
			t.limit40[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit48[run]); l > 0 {
			newkey := make([][6]uint64, l)
			copy(newkey, t.limit48[run])
			t.limit48[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit56[run]); l > 0 {
			newkey := make([][7]uint64, l)
			copy(newkey, t.limit56[run])
			t.limit56[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit64[run]); l > 0 {
			newkey := make([][8]uint64, l)
			copy(newkey, t.limit64[run])
			t.limit64[run] = newkey
		}
	}
}

// Reset() must be called before Next(). Returns whether there are any entries.
func (t *KeyBytes) Reset() bool {
	t.onlimit = 0
	t.on8 = 0
	t.oncursor = 0
	if len(t.limit8[0]) == 0 {
		if t.total == 0 {
			return false
		} else {
			t.forward(0)
		}
	}
	return true
}

func (t *KeyBytes) forward(l int) bool {
	t.oncursor++
	for t.oncursor >= l {
		t.oncursor = 0
		if t.on8++; t.on8 == 8 {
			t.on8 = 0
			if t.onlimit++; t.onlimit == 8 {
				t.Reset()
				return true
			}
		}
		switch t.onlimit {
			case 0: l = len(t.limit8[t.on8])
			case 1: l = len(t.limit16[t.on8])
			case 2: l = len(t.limit24[t.on8])
			case 3: l = len(t.limit32[t.on8])
			case 4: l = len(t.limit40[t.on8])
			case 5: l = len(t.limit48[t.on8])
			case 6: l = len(t.limit56[t.on8])
			case 7: l = len(t.limit64[t.on8])
		}
	}
	return false
}

func (t *KeyBytes) Next() ([]byte, bool) {
	switch t.onlimit {
		case 0:
			v := t.limit8[t.on8][t.oncursor]
			eof := t.forward(len(t.limit8[t.on8]))
			return reverse8(v), eof
		case 1:
			v := t.limit16[t.on8][t.oncursor]
			eof := t.forward(len(t.limit16[t.on8]))
			return reverse16(v), eof
		case 2:
			v := t.limit24[t.on8][t.oncursor]
			eof := t.forward(len(t.limit24[t.on8]))
			return reverse24(v), eof
		case 3:
			v := t.limit32[t.on8][t.oncursor]
			eof := t.forward(len(t.limit32[t.on8]))
			return reverse32(v), eof
		case 4:
			v := t.limit40[t.on8][t.oncursor]
			eof := t.forward(len(t.limit40[t.on8]))
			return reverse40(v), eof
		case 5:
			v := t.limit48[t.on8][t.oncursor]
			eof := t.forward(len(t.limit48[t.on8]))
			return reverse48(v), eof
		case 6:
			v := t.limit56[t.on8][t.oncursor]
			eof := t.forward(len(t.limit56[t.on8]))
			return reverse56(v), eof
		default:
			v := t.limit64[t.on8][t.oncursor]
			eof := t.forward(len(t.limit64[t.on8]))
			return reverse64(v), eof
	}
}

func (t *KeyBytes) Keys() [][]byte {

	var on, run int
	keys := make([][]byte, t.total)
	
	for run=0; run<8; run++ {
		for _, v := range t.limit8[run] {
			keys[on] = reverse8(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit16[run] {
			keys[on] = reverse16(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit24[run] {
			keys[on] = reverse24(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit32[run] {
			keys[on] = reverse32(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit40[run] {
			keys[on] = reverse40(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit48[run] {
			keys[on] = reverse48(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit56[run] {
			keys[on] = reverse56(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit64[run] {
			keys[on] = reverse64(v)
			on++
		}
	}
	
	return keys
}

func (t *KeyBytes) Write(w *custom.Writer) {
	var i, run int

	// Write total
	w.WriteUint64Variable(uint64(t.total))
	
	// Write count
	for i=0; i<64; i++ {
		w.WriteUint64Variable(uint64(t.count[i]))
	}
	
	// Write t.limit8
	for run=0; run<8; run++ {
		tmp := t.limit8[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v)
		}
	}
	// Write t.limit16
	for run=0; run<8; run++ {
		tmp := t.limit16[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
		}
	}
	// Write t.limit24
	for run=0; run<8; run++ {
		tmp := t.limit24[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
		}
	}
	// Write t.limit32
	for run=0; run<8; run++ {
		tmp := t.limit32[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
		}
	}
	// Write t.limit40
	for run=0; run<8; run++ {
		tmp := t.limit40[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
		}
	}
	// Write t.limit48
	for run=0; run<8; run++ {
		tmp := t.limit48[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
		}
	}
	// Write t.limit56
	for run=0; run<8; run++ {
		tmp := t.limit56[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
		}
	}
	// Write t.limit64
	for run=0; run<8; run++ {
		tmp := t.limit64[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
			w.WriteUint64(v[7])
		}
	}
	
}

func (t *KeyBytes) Read(r *custom.Reader) {
	var run int
	var i, l, a, b, c, d, e, f, g, h uint64

	// Write total
	t.total = int(r.ReadUint64Variable())
	
	// Read count
	for i=0; i<64; i++ {
		t.count[i] = int(r.ReadUint64Variable())
	}
	
	// Read t.limit8
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([]uint64, l)
		for i=0; i<l; i++ {
			tmp[i] = r.ReadUint64()
		}
		t.limit8[run] = tmp
	}
	// Read t.limit16
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][2]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			tmp[i] = [2]uint64{a, b}
		}
		t.limit16[run] = tmp
	}
	// Read t.limit24
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][3]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			tmp[i] = [3]uint64{a, b, c}
		}
		t.limit24[run] = tmp
	}
	// Read t.limit32
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][4]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			tmp[i] = [4]uint64{a, b, c, d}
		}
		t.limit32[run] = tmp
	}
	// Read t.limit40
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][5]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			tmp[i] = [5]uint64{a, b, c, d, e}
		}
		t.limit40[run] = tmp
	}
	// Read t.limit48
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][6]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			tmp[i] = [6]uint64{a, b, c, d, e, f}
		}
		t.limit48[run] = tmp
	}
	// Read t.limit56
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][7]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			tmp[i] = [7]uint64{a, b, c, d, e, f, g}
		}
		t.limit56[run] = tmp
	}
	// Read t.limit64
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][8]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			h = r.ReadUint64()
			tmp[i] = [8]uint64{a, b, c, d, e, f, g, h}
		}
		t.limit64[run] = tmp
	}
}

// ---------- KeyValBytes ----------
// Key bytes has around 2KB of memory overhead for the structures, beyond this it stores all keys as efficiently as possible.

// Add this to any struct to make it binary searchable.
type KeyValBytes struct {
 limit8 [8][][2]uint64 // where len(word) <= 8
 limit16 [8][][3]uint64
 limit24 [8][][4]uint64
 limit32 [8][][5]uint64
 limit40 [8][][6]uint64
 limit48 [8][][7]uint64
 limit56 [8][][8]uint64
 limit64 [8][][9]uint64
 total int
// Used for iterating through all of it
 onlimit int
 on8 int
 oncursor int
}

func (t *KeyValBytes) Len() int {
	return t.total
}

func (t *KeyValBytes) GreatestVal() int {
	var l, i2, this int
	var max int = -9223372036854775808
	for i:=0; i<8; i++ {
		l = len(t.limit8[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit8[i][i2][1])
			if this > max {
				max = this
			}
		}
		l = len(t.limit16[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit16[i][i2][2])
			if this > max {
				max = this
			}
		}
		l = len(t.limit24[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit24[i][i2][3])
			if this > max {
				max = this
			}
		}
		l = len(t.limit32[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit32[i][i2][4])
			if this > max {
				max = this
			}
		}
		l = len(t.limit40[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit40[i][i2][5])
			if this > max {
				max = this
			}
		}
		l = len(t.limit48[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit48[i][i2][6])
			if this > max {
				max = this
			}
		}
		l = len(t.limit56[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit56[i][i2][7])
			if this > max {
				max = this
			}
		}
		l = len(t.limit64[i])
		for i2=0; i2<l; i2++ {
			this = int(t.limit64[i][i2][8])
			if this > max {
				max = this
			}
		}
	}
	return max
}

// Find returns the index based on the key.
func (t *KeyValBytes) Find(thekey []byte) (int, bool) {
	
	var at, min int
	var compare uint64
	switch (len(thekey) - 1) / 8 {
	
		case 0: // 0 - 8 bytes
			a, l := bytes2uint64(thekey)
			cur := t.limit8[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				return int(cur[at][1]), true // found
			}
			return 0, false // doesn't exist
			
		case 1: // 9 - 16 bytes
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			cur := t.limit16[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				return int(cur[at][2]), true // found
			}
			return 0, false // doesn't exist
			
		case 2: // 17 - 24 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			cur := t.limit24[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				return int(cur[at][3]), true // found
			}
			return 0, false // doesn't exist
			
		case 3: // 25 - 32 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			cur := t.limit32[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				return int(cur[at][4]), true // found
			}
			return 0, false // doesn't exist
			
		case 4: // 33 - 40 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			cur := t.limit40[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				return int(cur[at][5]), true // found
			}
			return 0, false // doesn't exist
			
		case 5: // 41 - 48 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			cur := t.limit48[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				return int(cur[at][6]), true // found
			}
			return 0, false // doesn't exist
			
		case 6: // 49 - 56 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			cur := t.limit56[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				return int(cur[at][7]), true // found
			}
			return 0, false // doesn't exist
			
		case 7: // 57 - 64 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			cur := t.limit64[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][7]; h < compare {
					max = at - 1
					continue
				}
				if h > compare {
					min = at + 1
					continue
				}
				return int(cur[at][8]), true // found
			}
			return 0, false // doesn't exist
		
		default: // > 64 bytes
			return t.total, false
	}
}

// Modifies the value of the key by running it through the provided function.
func (t *KeyValBytes) Update(thekey []byte, fn func(int) int) bool {
	
	var at, min int
	var compare uint64
	switch (len(thekey) - 1) / 8 {
	
		case 0: // 0 - 8 bytes
			a, l := bytes2uint64(thekey)
			cur := t.limit8[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				cur[at][1] = uint64(fn(int(cur[at][1])))
				return true // found
			}
			return false // doesn't exist
			
		case 1: // 9 - 16 bytes
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			cur := t.limit16[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				cur[at][2] = uint64(fn(int(cur[at][2])))
				return true // found
			}
			return false // doesn't exist
			
		case 2: // 17 - 24 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			cur := t.limit24[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				cur[at][3] = uint64(fn(int(cur[at][3])))
				return true // found
			}
			return false // doesn't exist
			
		case 3: // 25 - 32 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			cur := t.limit32[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				cur[at][4] = uint64(fn(int(cur[at][4])))
				return true // found
			}
			return false // doesn't exist
			
		case 4: // 33 - 40 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			cur := t.limit40[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				cur[at][5] = uint64(fn(int(cur[at][5])))
				return true // found
			}
			return false // doesn't exist
			
		case 5: // 41 - 48 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			cur := t.limit48[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				cur[at][6] = uint64(fn(int(cur[at][6])))
				return true // found
			}
			return false // doesn't exist
			
		case 6: // 49 - 56 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			cur := t.limit56[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				cur[at][7] = uint64(fn(int(cur[at][7])))
				return true // found
			}
			return false // doesn't exist
			
		case 7: // 57 - 64 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			cur := t.limit64[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][7]; h < compare {
					max = at - 1
					continue
				}
				if h > compare {
					min = at + 1
					continue
				}
				cur[at][8] = uint64(fn(int(cur[at][8])))
				return true // found
			}
			return false // doesn't exist
		
		default: // > 64 bytes
			return false
	}
}

// Modifies all values by running each through the provided function.
func (t *KeyValBytes) UpdateAll(fn func(int) int) {
	var run, l, i int
	for run=0; run<8; run++ {
		tmp := t.limit8[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][1] = uint64(fn(int(tmp[i][1])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit16[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][2] = uint64(fn(int(tmp[i][2])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit24[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][3] = uint64(fn(int(tmp[i][3])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit32[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][4] = uint64(fn(int(tmp[i][4])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit40[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][5] = uint64(fn(int(tmp[i][5])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit48[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][6] = uint64(fn(int(tmp[i][6])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit56[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][7] = uint64(fn(int(tmp[i][7])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit64[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][8] = uint64(fn(int(tmp[i][8])))
		}
	}
}

// Add is equivalent to Find and then AddAt
func (t *KeyValBytes) Add(thekey []byte, theval int) bool {
	
	var at, min int
	var compare uint64
	switch (len(thekey) - 1) / 8 {
	
		case 0: // 0 - 8 bytes
			a, l := bytes2uint64(thekey)
			cur := t.limit8[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if cur[at][1] != uint64(theval) {
					cur[at][1] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][2]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [2]uint64{a, uint64(theval)}
			t.limit8[l] = cur
			t.total++
			return false
			
		case 1: // 9 - 16 bytes
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			cur := t.limit16[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if cur[at][2] != uint64(theval) {
					cur[at][2] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][3]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [3]uint64{a, b, uint64(theval)}
			t.limit16[l] = cur
			t.total++
			return false
			
		case 2: // 17 - 24 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			cur := t.limit24[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if cur[at][3] != uint64(theval) {
					cur[at][3] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][4]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [4]uint64{a, b, c, uint64(theval)}
			t.limit24[l] = cur
			t.total++
			return false
			
		case 3: // 25 - 32 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			cur := t.limit32[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if cur[at][4] != uint64(theval) {
					cur[at][4] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][5]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [5]uint64{a, b, c, d, uint64(theval)}
			t.limit32[l] = cur
			t.total++
			return false
			
		case 4: // 33 - 40 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			cur := t.limit40[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if cur[at][5] != uint64(theval) {
					cur[at][5] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][6]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [6]uint64{a, b, c, d, e, uint64(theval)}
			t.limit40[l] = cur
			t.total++
			return false
			
		case 5: // 41 - 48 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			cur := t.limit48[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if cur[at][6] != uint64(theval) {
					cur[at][6] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][7]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [7]uint64{a, b, c, d, e, f, uint64(theval)}
			t.limit48[l] = cur
			t.total++
			return false
			
		case 6: // 49 - 56 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			cur := t.limit56[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if cur[at][7] != uint64(theval) {
					cur[at][7] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][8]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [8]uint64{a, b, c, d, e, f, g, uint64(theval)}
			t.limit56[l] = cur
			t.total++
			return false
			
		case 7: // 57 - 64 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			cur := t.limit64[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][7]; h < compare {
					max = at - 1
					continue
				}
				if h > compare {
					min = at + 1
					continue
				}
				if cur[at][8] != uint64(theval) {
					cur[at][8] = uint64(theval)
				}
				return true // found
			}
			// Doesn't exist so add it >
			at = min
			lc := len(cur)
			if lc == cap(cur) {
				tmp := make([][9]uint64, lc + 1, (lc * 2) + 1)
				copy(tmp, cur[0:at])
				copy(tmp[at+1:], cur[at:])
				cur = tmp
			} else {
				cur = cur[0:lc+1]
				copy(cur[at+1:], cur[at:])
			}
			cur[at] = [9]uint64{a, b, c, d, e, f, g, h, uint64(theval)}
			t.limit64[l] = cur
			t.total++
			return false
		
		default: // > 64 bytes
			return false
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValBytes) AddUnsorted(thekey []byte, theval int) error {
	switch (len(thekey) - 1) / 8 {
		case 0:
			a, i := bytes2uint64(thekey)
			t.limit8[i] = append(t.limit8[i], [2]uint64{a, uint64(theval)})
			t.total++
			return nil
		case 1:
			a, _ := bytes2uint64(thekey)
			b, i := bytes2uint64(thekey[8:])
			t.limit16[i] = append(t.limit16[i], [3]uint64{a, b, uint64(theval)})
			t.total++
			return nil
		case 2:
			
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, i := bytes2uint64(thekey[16:])
			t.limit24[i] = append(t.limit24[i], [4]uint64{a, b, c, uint64(theval)})
			t.total++
			return nil
		case 3:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, i := bytes2uint64(thekey[24:])
			t.limit32[i] = append(t.limit32[i], [5]uint64{a, b, c, d, uint64(theval)})
			t.total++
			return nil
		case 4:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, i := bytes2uint64(thekey[32:])
			t.limit40[i] = append(t.limit40[i], [6]uint64{a, b, c, d, e, uint64(theval)})
			t.total++
			return nil
		case 5:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, i := bytes2uint64(thekey[40:])
			t.limit48[i] = append(t.limit48[i], [7]uint64{a, b, c, d, e, f, uint64(theval)})
			t.total++
			return nil
		case 6:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, i := bytes2uint64(thekey[48:])
			t.limit56[i] = append(t.limit56[i], [8]uint64{a, b, c, d, e, f, g, uint64(theval)})
			t.total++
			return nil
		case 7:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, i := bytes2uint64(thekey[56:])
			t.limit64[i] = append(t.limit64[i], [9]uint64{a, b, c, d, e, f, g, h, uint64(theval)})
			t.total++
			return nil
		default:
			return errors.New(`Maximum key length is 64 bytes`)
	}
}

// Build sorts the keys
func (t *KeyValBytes) Build() {

	var run int
	
	for run=0; run<8; run++ {
		if len(t.limit8[run]) > 0 {
			sortLimitVal8.Asc(sortLimitVal8.Slice(t.limit8[run]))
		}
	}
	
	for run=0; run<8; run++ {
		if len(t.limit16[run]) > 0 {
			sortLimitVal16.Asc(sortLimitVal16.Slice(t.limit16[run]))
		}
	}
	
	for run=0; run<8; run++ {
		if len(t.limit24[run]) > 0 {
			sortLimitVal24.Asc(sortLimitVal24.Slice(t.limit24[run]))
		}
	}
	
	for run=0; run<8; run++ {
		if len(t.limit32[run]) > 0 {
			sortLimitVal32.Asc(sortLimitVal32.Slice(t.limit32[run]))
		}
	}
	
	for run=0; run<8; run++ {
		if len(t.limit40[run]) > 0 {
			sortLimitVal40.Asc(sortLimitVal40.Slice(t.limit40[run]))
		}
	}
	
	for run=0; run<8; run++ {
		if len(t.limit48[run]) > 0 {
			sortLimitVal48.Asc(sortLimitVal48.Slice(t.limit48[run]))
		}
	}
	
	for run=0; run<8; run++ {
		if len(t.limit56[run]) > 0 {
			sortLimitVal56.Asc(sortLimitVal56.Slice(t.limit56[run]))
		}
	}
	
	for run=0; run<8; run++ {
		if len(t.limit64[run]) > 0 {
			sortLimitVal64.Asc(sortLimitVal64.Slice(t.limit64[run]))
		}
	}
}

func (t *KeyValBytes) Optimize() {

	var l, run int
	
	for run=0; run<8; run++ {
		if l = len(t.limit8[run]); l > 0 {
			newkey := make([][2]uint64, l)
			copy(newkey, t.limit8[run])
			t.limit8[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit16[run]); l > 0 {
			newkey := make([][3]uint64, l)
			copy(newkey, t.limit16[run])
			t.limit16[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit24[run]); l > 0 {
			newkey := make([][4]uint64, l)
			copy(newkey, t.limit24[run])
			t.limit24[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit32[run]); l > 0 {
			newkey := make([][5]uint64, l)
			copy(newkey, t.limit32[run])
			t.limit32[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit40[run]); l > 0 {
			newkey := make([][6]uint64, l)
			copy(newkey, t.limit40[run])
			t.limit40[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit48[run]); l > 0 {
			newkey := make([][7]uint64, l)
			copy(newkey, t.limit48[run])
			t.limit48[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit56[run]); l > 0 {
			newkey := make([][8]uint64, l)
			copy(newkey, t.limit56[run])
			t.limit56[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit64[run]); l > 0 {
			newkey := make([][9]uint64, l)
			copy(newkey, t.limit64[run])
			t.limit64[run] = newkey
		}
	}
}

// Reset() must be called before Next(). Returns whether there are any entries.
func (t *KeyValBytes) Reset() bool {
	t.onlimit = 0
	t.on8 = 0
	t.oncursor = 0
	if len(t.limit8[0]) == 0 {
		if t.total == 0 {
			return false
		} else {
			t.forward(0)
		}
	}
	return true
}

func (t *KeyValBytes) forward(l int) bool {
	t.oncursor++
	for t.oncursor >= l {
		t.oncursor = 0
		if t.on8++; t.on8 == 8 {
			t.on8 = 0
			if t.onlimit++; t.onlimit == 8 {
				t.Reset()
				return true
			}
		}
		switch t.onlimit {
			case 0: l = len(t.limit8[t.on8])
			case 1: l = len(t.limit16[t.on8])
			case 2: l = len(t.limit24[t.on8])
			case 3: l = len(t.limit32[t.on8])
			case 4: l = len(t.limit40[t.on8])
			case 5: l = len(t.limit48[t.on8])
			case 6: l = len(t.limit56[t.on8])
			case 7: l = len(t.limit64[t.on8])
		}
	}
	return false
}

func (t *KeyValBytes) Next() ([]byte, int, bool) {
	switch t.onlimit {
		case 0:
			v := t.limit8[t.on8][t.oncursor]
			eof := t.forward(len(t.limit8[t.on8]))
			return reverse8b(v), int(v[1]), eof
		case 1:
			v := t.limit16[t.on8][t.oncursor]
			eof := t.forward(len(t.limit16[t.on8]))
			return reverse16b(v), int(v[2]), eof
		case 2:
			v := t.limit24[t.on8][t.oncursor]
			eof := t.forward(len(t.limit24[t.on8]))
			return reverse24b(v), int(v[3]), eof
		case 3:
			v := t.limit32[t.on8][t.oncursor]
			eof := t.forward(len(t.limit32[t.on8]))
			return reverse32b(v), int(v[4]), eof
		case 4:
			v := t.limit40[t.on8][t.oncursor]
			eof := t.forward(len(t.limit40[t.on8]))
			return reverse40b(v), int(v[5]), eof
		case 5:
			v := t.limit48[t.on8][t.oncursor]
			eof := t.forward(len(t.limit48[t.on8]))
			return reverse48b(v), int(v[6]), eof
		case 6:
			v := t.limit56[t.on8][t.oncursor]
			eof := t.forward(len(t.limit56[t.on8]))
			return reverse56b(v), int(v[7]), eof
		default:
			v := t.limit64[t.on8][t.oncursor]
			eof := t.forward(len(t.limit64[t.on8]))
			return reverse64b(v), int(v[8]), eof
	}
}

func (t *KeyValBytes) Keys() [][]byte {

	var on, run int
	keys := make([][]byte, t.total)
	
	for run=0; run<8; run++ {
		for _, v := range t.limit8[run] {
			keys[on] = reverse8b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit16[run] {
			keys[on] = reverse16b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit24[run] {
			keys[on] = reverse24b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit32[run] {
			keys[on] = reverse32b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit40[run] {
			keys[on] = reverse40b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit48[run] {
			keys[on] = reverse48b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit56[run] {
			keys[on] = reverse56b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit64[run] {
			keys[on] = reverse64b(v)
			on++
		}
	}
	
	return keys
}


func (t *KeyValBytes) Write(w *custom.Writer) {
	var run int

	// Write total
	w.WriteUint64Variable(uint64(t.total))
	
	// Write t.limit8
	for run=0; run<8; run++ {
		tmp := t.limit8[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
		}
	}
	// Write t.limit16
	for run=0; run<8; run++ {
		tmp := t.limit16[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
		}
	}
	// Write t.limit24
	for run=0; run<8; run++ {
		tmp := t.limit24[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
		}
	}
	// Write t.limit32
	for run=0; run<8; run++ {
		tmp := t.limit32[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
		}
	}
	// Write t.limit40
	for run=0; run<8; run++ {
		tmp := t.limit40[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
		}
	}
	// Write t.limit48
	for run=0; run<8; run++ {
		tmp := t.limit48[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
		}
	}
	// Write t.limit56
	for run=0; run<8; run++ {
		tmp := t.limit56[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
			w.WriteUint64(v[7])
		}
	}
	// Write t.limit64
	for run=0; run<8; run++ {
		tmp := t.limit64[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
			w.WriteUint64(v[7])
			w.WriteUint64(v[8])
		}
	}
}

func (t *KeyValBytes) Read(r *custom.Reader) {
	var run int
	var i, l, a, b, c, d, e, f, g, h, z uint64

	// Write total
	t.total = int(r.ReadUint64Variable())
	
	// Read t.limit8
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][2]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			tmp[i] = [2]uint64{a, b}
		}
		t.limit8[run] = tmp
	}
	// Read t.limit16
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][3]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			tmp[i] = [3]uint64{a, b, c}
		}
		t.limit16[run] = tmp
	}
	// Read t.limit24
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][4]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			tmp[i] = [4]uint64{a, b, c, d}
		}
		t.limit24[run] = tmp
	}
	// Read t.limit32
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][5]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			tmp[i] = [5]uint64{a, b, c, d, e}
		}
		t.limit32[run] = tmp
	}
	// Read t.limit40
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][6]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			tmp[i] = [6]uint64{a, b, c, d, e, f}
		}
		t.limit40[run] = tmp
	}
	// Read t.limit48
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][7]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			tmp[i] = [7]uint64{a, b, c, d, e, f, g}
		}
		t.limit48[run] = tmp
	}
	// Read t.limit56
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][8]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			h = r.ReadUint64()
			tmp[i] = [8]uint64{a, b, c, d, e, f, g, h}
		}
		t.limit56[run] = tmp
	}
	// Read t.limit64
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][9]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			h = r.ReadUint64()
			z = r.ReadUint64()
			tmp[i] = [9]uint64{a, b, c, d, e, f, g, h, z}
		}
		t.limit64[run] = tmp
	}
}

// ---------- CounterBytes ----------
// CounterBytes bytes has around 2KB of memory overhead for the structures, beyond this it stores all keys as efficiently as possible.

// Add this to any struct to make it binary searchable.
type CounterBytes struct {
 limit8 [8][][2]uint64 // where len(word) <= 8
 limit16 [8][][3]uint64
 limit24 [8][][4]uint64
 limit32 [8][][5]uint64
 limit40 [8][][6]uint64
 limit48 [8][][7]uint64
 limit56 [8][][8]uint64
 limit64 [8][][9]uint64
 total int
// Used for iterating through all of it
 onlimit int
 on8 int
 oncursor int
}

func (t *CounterBytes) KeyBytes() *KeyBytes {
	obj := new(KeyBytes)
	obj.total = t.total
	var run int
	for run=0; run<8; run++ {
		tmp := t.limit8[run]
		cpy := make([]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i] = v[0]
		}
		obj.limit8[run] = cpy
		obj.count[run + 1] = len(cpy)
	}
	for run=0; run<8; run++ {
		tmp := t.limit16[run]
		cpy := make([][2]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i][0] = v[0]
			cpy[i][1] = v[1]
		}
		obj.limit16[run] = cpy
		obj.count[run + 9] = len(cpy)
	}
	for run=0; run<8; run++ {
		tmp := t.limit24[run]
		cpy := make([][3]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i][0] = v[0]
			cpy[i][1] = v[1]
			cpy[i][2] = v[2]
		}
		obj.limit24[run] = cpy
		obj.count[run + 17] = len(cpy)
	}
	for run=0; run<8; run++ {
		tmp := t.limit32[run]
		cpy := make([][4]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i][0] = v[0]
			cpy[i][1] = v[1]
			cpy[i][2] = v[2]
			cpy[i][3] = v[3]
		}
		obj.limit32[run] = cpy
		obj.count[run + 25] = len(cpy)
	}
	for run=0; run<8; run++ {
		tmp := t.limit40[run]
		cpy := make([][5]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i][0] = v[0]
			cpy[i][1] = v[1]
			cpy[i][2] = v[2]
			cpy[i][3] = v[3]
			cpy[i][4] = v[4]
		}
		obj.limit40[run] = cpy
		obj.count[run + 33] = len(cpy)
	}
	for run=0; run<8; run++ {
		tmp := t.limit48[run]
		cpy := make([][6]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i][0] = v[0]
			cpy[i][1] = v[1]
			cpy[i][2] = v[2]
			cpy[i][3] = v[3]
			cpy[i][4] = v[4]
			cpy[i][5] = v[5]
		}
		obj.limit48[run] = cpy
		obj.count[run + 41] = len(cpy)
	}
	for run=0; run<8; run++ {
		tmp := t.limit56[run]
		cpy := make([][7]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i][0] = v[0]
			cpy[i][1] = v[1]
			cpy[i][2] = v[2]
			cpy[i][3] = v[3]
			cpy[i][4] = v[4]
			cpy[i][5] = v[5]
			cpy[i][6] = v[6]
		}
		obj.limit56[run] = cpy
		obj.count[run + 49] = len(cpy)
	}
	for run=0; run<8; run++ {
		tmp := t.limit64[run]
		cpy := make([][8]uint64, len(tmp))
		for i, v := range tmp {
			cpy[i][0] = v[0]
			cpy[i][1] = v[1]
			cpy[i][2] = v[2]
			cpy[i][3] = v[3]
			cpy[i][4] = v[4]
			cpy[i][5] = v[5]
			cpy[i][6] = v[6]
			cpy[i][7] = v[7]
		}
		obj.limit64[run] = cpy
		if run < 7 {
			obj.count[run + 57] = len(cpy)
		}
	}
	// Correct all the counts
	for run=2; run<64; run++ {
		obj.count[run] += obj.count[run-1]
	}
	return obj
}

func (t *CounterBytes) KeyValBytes() *KeyValBytes {
	obj := new(KeyValBytes)
	obj.total = t.total
	var run int
	for run=0; run<8; run++ {
		cpy := make([][2]uint64, len(t.limit8[run]))
		copy(cpy, t.limit8[run])
		obj.limit8[run] = cpy
	}
	for run=0; run<8; run++ {
		cpy := make([][3]uint64, len(t.limit16[run]))
		copy(cpy, t.limit16[run])
		obj.limit16[run] = cpy
	}
	for run=0; run<8; run++ {
		cpy := make([][4]uint64, len(t.limit24[run]))
		copy(cpy, t.limit24[run])
		obj.limit24[run] = cpy
	}
	for run=0; run<8; run++ {
		cpy := make([][5]uint64, len(t.limit32[run]))
		copy(cpy, t.limit32[run])
		obj.limit32[run] = cpy
	}
	for run=0; run<8; run++ {
		cpy := make([][6]uint64, len(t.limit40[run]))
		copy(cpy, t.limit40[run])
		obj.limit40[run] = cpy
	}
	for run=0; run<8; run++ {
		cpy := make([][7]uint64, len(t.limit48[run]))
		copy(cpy, t.limit48[run])
		obj.limit48[run] = cpy
	}
	for run=0; run<8; run++ {
		cpy := make([][8]uint64, len(t.limit56[run]))
		copy(cpy, t.limit56[run])
		obj.limit56[run] = cpy
	}
	for run=0; run<8; run++ {
		cpy := make([][9]uint64, len(t.limit64[run]))
		copy(cpy, t.limit64[run])
		obj.limit64[run] = cpy
	}
	return obj
}

func (t *CounterBytes) Len() int {
	return t.total
}

// Find returns the index based on the key.
func (t *CounterBytes) Find(thekey []byte) (int, bool) {
	
	var at, min int
	var compare uint64
	switch (len(thekey) - 1) / 8 {
	
		case 0: // 0 - 8 bytes
			a, l := bytes2uint64(thekey)
			cur := t.limit8[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				return int(cur[at][1]), true // found
			}
			return 0, false // doesn't exist
			
		case 1: // 9 - 16 bytes
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			cur := t.limit16[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				return int(cur[at][2]), true // found
			}
			return 0, false // doesn't exist
			
		case 2: // 17 - 24 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			cur := t.limit24[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				return int(cur[at][3]), true // found
			}
			return 0, false // doesn't exist
			
		case 3: // 25 - 32 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			cur := t.limit32[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				return int(cur[at][4]), true // found
			}
			return 0, false // doesn't exist
			
		case 4: // 33 - 40 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			cur := t.limit40[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				return int(cur[at][5]), true // found
			}
			return 0, false // doesn't exist
			
		case 5: // 41 - 48 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			cur := t.limit48[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				return int(cur[at][6]), true // found
			}
			return 0, false // doesn't exist
			
		case 6: // 49 - 56 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			cur := t.limit56[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				return int(cur[at][7]), true // found
			}
			return 0, false // doesn't exist
			
		case 7: // 57 - 64 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			cur := t.limit64[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][7]; h < compare {
					max = at - 1
					continue
				}
				if h > compare {
					min = at + 1
					continue
				}
				return int(cur[at][8]), true // found
			}
			return 0, false // doesn't exist
		
		default: // > 64 bytes
			return t.total, false
	}
}

// Modifies the value of the key by running it through the provided function.
func (t *CounterBytes) Update(thekey []byte, fn func(int) int) bool {
	
	var at, min int
	var compare uint64
	switch (len(thekey) - 1) / 8 {
	
		case 0: // 0 - 8 bytes
			a, l := bytes2uint64(thekey)
			cur := t.limit8[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				cur[at][1] = uint64(fn(int(cur[at][1])))
				return true // found
			}
			return false // doesn't exist
			
		case 1: // 9 - 16 bytes
			a, _ := bytes2uint64(thekey)
			b, l := bytes2uint64(thekey[8:])
			cur := t.limit16[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				cur[at][2] = uint64(fn(int(cur[at][2])))
				return true // found
			}
			return false // doesn't exist
			
		case 2: // 17 - 24 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, l := bytes2uint64(thekey[16:])
			cur := t.limit24[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				cur[at][3] = uint64(fn(int(cur[at][3])))
				return true // found
			}
			return false // doesn't exist
			
		case 3: // 25 - 32 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, l := bytes2uint64(thekey[24:])
			cur := t.limit32[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				cur[at][4] = uint64(fn(int(cur[at][4])))
				return true // found
			}
			return false // doesn't exist
			
		case 4: // 33 - 40 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, l := bytes2uint64(thekey[32:])
			cur := t.limit40[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				cur[at][5] = uint64(fn(int(cur[at][5])))
				return true // found
			}
			return false // doesn't exist
			
		case 5: // 41 - 48 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, l := bytes2uint64(thekey[40:])
			cur := t.limit48[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				cur[at][6] = uint64(fn(int(cur[at][6])))
				return true // found
			}
			return false // doesn't exist
			
		case 6: // 49 - 56 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, l := bytes2uint64(thekey[48:])
			cur := t.limit56[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				cur[at][7] = uint64(fn(int(cur[at][7])))
				return true // found
			}
			return false // doesn't exist
			
		case 7: // 57 - 64 bytes
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, l := bytes2uint64(thekey[56:])
			cur := t.limit64[l]
			max := len(cur) - 1
			for min <= max {
				at = min + ((max - min) / 2)
				if compare = cur[at][0]; a < compare {
					max = at - 1
					continue
				}
				if a > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][1]; b < compare {
					max = at - 1
					continue
				}
				if b > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][2]; c < compare {
					max = at - 1
					continue
				}
				if c > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][3]; d < compare {
					max = at - 1
					continue
				}
				if d > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][4]; e < compare {
					max = at - 1
					continue
				}
				if e > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][5]; f < compare {
					max = at - 1
					continue
				}
				if f > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][6]; g < compare {
					max = at - 1
					continue
				}
				if g > compare {
					min = at + 1
					continue
				}
				if compare = cur[at][7]; h < compare {
					max = at - 1
					continue
				}
				if h > compare {
					min = at + 1
					continue
				}
				cur[at][8] = uint64(fn(int(cur[at][8])))
				return true // found
			}
			return false // doesn't exist
		
		default: // > 64 bytes
			return false
	}
}

// Modifies all values by running each through the provided function.
func (t *CounterBytes) UpdateAll(fn func(int) int) {
	var run, l, i int
	for run=0; run<8; run++ {
		tmp := t.limit8[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][1] = uint64(fn(int(tmp[i][1])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit16[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][2] = uint64(fn(int(tmp[i][2])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit24[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][3] = uint64(fn(int(tmp[i][3])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit32[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][4] = uint64(fn(int(tmp[i][4])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit40[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][5] = uint64(fn(int(tmp[i][5])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit48[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][6] = uint64(fn(int(tmp[i][6])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit56[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][7] = uint64(fn(int(tmp[i][7])))
		}
	}
	for run=0; run<8; run++ {
		tmp := t.limit64[run]
		l = len(tmp)
		for i=0; i<l; i++ {
			tmp[i][8] = uint64(fn(int(tmp[i][8])))
		}
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *CounterBytes) Add(thekey []byte, theval int) error {
	switch (len(thekey) - 1) / 8 {
		case 0:
			a, i := bytes2uint64(thekey)
			t.limit8[i] = append(t.limit8[i], [2]uint64{a, uint64(theval)})
			t.total++
			return nil
		case 1:
			a, _ := bytes2uint64(thekey)
			b, i := bytes2uint64(thekey[8:])
			t.limit16[i] = append(t.limit16[i], [3]uint64{a, b, uint64(theval)})
			t.total++
			return nil
		case 2:
			
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, i := bytes2uint64(thekey[16:])
			t.limit24[i] = append(t.limit24[i], [4]uint64{a, b, c, uint64(theval)})
			t.total++
			return nil
		case 3:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, i := bytes2uint64(thekey[24:])
			t.limit32[i] = append(t.limit32[i], [5]uint64{a, b, c, d, uint64(theval)})
			t.total++
			return nil
		case 4:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, i := bytes2uint64(thekey[32:])
			t.limit40[i] = append(t.limit40[i], [6]uint64{a, b, c, d, e, uint64(theval)})
			t.total++
			return nil
		case 5:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, i := bytes2uint64(thekey[40:])
			t.limit48[i] = append(t.limit48[i], [7]uint64{a, b, c, d, e, f, uint64(theval)})
			t.total++
			return nil
		case 6:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, i := bytes2uint64(thekey[48:])
			t.limit56[i] = append(t.limit56[i], [8]uint64{a, b, c, d, e, f, g, uint64(theval)})
			t.total++
			return nil
		case 7:
			a, _ := bytes2uint64(thekey)
			b, _ := bytes2uint64(thekey[8:])
			c, _ := bytes2uint64(thekey[16:])
			d, _ := bytes2uint64(thekey[24:])
			e, _ := bytes2uint64(thekey[32:])
			f, _ := bytes2uint64(thekey[40:])
			g, _ := bytes2uint64(thekey[48:])
			h, i := bytes2uint64(thekey[56:])
			t.limit64[i] = append(t.limit64[i], [9]uint64{a, b, c, d, e, f, g, h, uint64(theval)})
			t.total++
			return nil
		default:
			return errors.New(`Maximum key length is 64 bytes`)
	}
}

func (t *CounterBytes) Build() {

	var l, run, n, total, on int
	
	for run=0; run<8; run++ {
		if l = len(t.limit8[run]); l > 0 {
			var temp sortLimitVal8.Slice = t.limit8[run]
			sortLimitVal8.Asc(temp)
			this := temp[0]
			n = int(temp[0][1])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] {
					n += int(k[1])
				} else {
					this[1] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[1])
				}
			}
			this[1] = uint64(n)
			temp[on] = this
			on++
			t.limit8[run] = temp[0:on]
			total += on
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit16[run]); l > 0 {
			var temp sortLimitVal16.Slice = t.limit16[run]
			sortLimitVal16.Asc(temp)
			this := temp[0]
			n = int(temp[0][2])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] && k[1] == this[1] {
					n += int(k[2])
				} else {
					this[2] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[2])
				}
			}
			this[2] = uint64(n)
			temp[on] = this
			on++
			t.limit16[run] = temp[0:on]
			total += on
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit24[run]); l > 0 {
			var temp sortLimitVal24.Slice = t.limit24[run]
			sortLimitVal24.Asc(temp)
			this := temp[0]
			n = int(temp[0][3])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] && k[1] == this[1] && k[2] == this[2] {
					n += int(k[3])
				} else {
					this[3] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[3])
				}
			}
			this[3] = uint64(n)
			temp[on] = this
			on++
			t.limit24[run] = temp[0:on]
			total += on
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit32[run]); l > 0 {
			var temp sortLimitVal32.Slice = t.limit32[run]
			sortLimitVal32.Asc(temp)
			this := temp[0]
			n = int(temp[0][4])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] && k[1] == this[1] && k[2] == this[2] && k[3] == this[3] {
					n += int(k[4])
				} else {
					this[4] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[4])
				}
			}
			this[4] = uint64(n)
			temp[on] = this
			on++
			t.limit32[run] = temp[0:on]
			total += on
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit40[run]); l > 0 {
			var temp sortLimitVal40.Slice = t.limit40[run]
			sortLimitVal40.Asc(temp)
			this := temp[0]
			n = int(temp[0][5])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] && k[1] == this[1] && k[2] == this[2] && k[3] == this[3] && k[4] == this[4] {
					n += int(k[5])
				} else {
					this[5] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[5])
				}
			}
			this[5] = uint64(n)
			temp[on] = this
			on++
			t.limit40[run] = temp[0:on]
			total += on
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit48[run]); l > 0 {
			var temp sortLimitVal48.Slice = t.limit48[run]
			sortLimitVal48.Asc(temp)
			this := temp[0]
			n = int(temp[0][6])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] && k[1] == this[1] && k[2] == this[2] && k[3] == this[3] && k[4] == this[4] && k[5] == this[5] {
					n += int(k[6])
				} else {
					this[6] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[6])
				}
			}
			this[6] = uint64(n)
			temp[on] = this
			on++
			t.limit48[run] = temp[0:on]
			total += on
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit56[run]); l > 0 {
			var temp sortLimitVal56.Slice = t.limit56[run]
			sortLimitVal56.Asc(temp)
			this := temp[0]
			n = int(temp[0][7])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] && k[1] == this[1] && k[2] == this[2] && k[3] == this[3] && k[4] == this[4] && k[5] == this[5] && k[6] == this[6] {
					n += int(k[7])
				} else {
					this[7] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[7])
				}
			}
			this[7] = uint64(n)
			temp[on] = this
			on++
			t.limit56[run] = temp[0:on]
			total += on
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit64[run]); l > 0 {
			var temp sortLimitVal64.Slice = t.limit64[run]
			sortLimitVal64.Asc(temp)
			this := temp[0]
			n = int(temp[0][8])
			on = 0
			for _, k := range temp[1:] {
				if k[0] == this[0] && k[1] == this[1] && k[2] == this[2] && k[3] == this[3] && k[4] == this[4] && k[5] == this[5] && k[6] == this[6] && k[7] == this[7] {
					n += int(k[8])
				} else {
					this[8] = uint64(n)
					temp[on] = this
					on++
					this = k
					n = int(k[8])
				}
			}
			this[8] = uint64(n)
			temp[on] = this
			on++
			t.limit64[run] = temp[0:on]
			total += on
		}
	}
	t.total = total
	
}

func (t *CounterBytes) Optimize() {

	var l, run int
	
	for run=0; run<8; run++ {
		if l = len(t.limit8[run]); l > 0 {
			newkey := make([][2]uint64, l)
			copy(newkey, t.limit8[run])
			t.limit8[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit16[run]); l > 0 {
			newkey := make([][3]uint64, l)
			copy(newkey, t.limit16[run])
			t.limit16[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit24[run]); l > 0 {
			newkey := make([][4]uint64, l)
			copy(newkey, t.limit24[run])
			t.limit24[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit32[run]); l > 0 {
			newkey := make([][5]uint64, l)
			copy(newkey, t.limit32[run])
			t.limit32[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit40[run]); l > 0 {
			newkey := make([][6]uint64, l)
			copy(newkey, t.limit40[run])
			t.limit40[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit48[run]); l > 0 {
			newkey := make([][7]uint64, l)
			copy(newkey, t.limit48[run])
			t.limit48[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit56[run]); l > 0 {
			newkey := make([][8]uint64, l)
			copy(newkey, t.limit56[run])
			t.limit56[run] = newkey
		}
	}
	
	for run=0; run<8; run++ {
		if l = len(t.limit64[run]); l > 0 {
			newkey := make([][9]uint64, l)
			copy(newkey, t.limit64[run])
			t.limit64[run] = newkey
		}
	}
}

// Reset() must be called before Next(). Returns whether there are any entries.
func (t *CounterBytes) Reset() bool {
	t.onlimit = 0
	t.on8 = 0
	t.oncursor = 0
	if len(t.limit8[0]) == 0 {
		if t.total == 0 {
			return false
		} else {
			t.forward(0)
		}
	}
	return true
}

func (t *CounterBytes) forward(l int) bool {
	t.oncursor++
	for t.oncursor >= l {
		t.oncursor = 0
		if t.on8++; t.on8 == 8 {
			t.on8 = 0
			if t.onlimit++; t.onlimit == 8 {
				t.Reset()
				return true
			}
		}
		switch t.onlimit {
			case 0: l = len(t.limit8[t.on8])
			case 1: l = len(t.limit16[t.on8])
			case 2: l = len(t.limit24[t.on8])
			case 3: l = len(t.limit32[t.on8])
			case 4: l = len(t.limit40[t.on8])
			case 5: l = len(t.limit48[t.on8])
			case 6: l = len(t.limit56[t.on8])
			case 7: l = len(t.limit64[t.on8])
		}
	}
	return false
}

func (t *CounterBytes) Next() ([]byte, int, bool) {
	switch t.onlimit {
		case 0:
			v := t.limit8[t.on8][t.oncursor]
			eof := t.forward(len(t.limit8[t.on8]))
			return reverse8b(v), int(v[1]), eof
		case 1:
			v := t.limit16[t.on8][t.oncursor]
			eof := t.forward(len(t.limit16[t.on8]))
			return reverse16b(v), int(v[2]), eof
		case 2:
			v := t.limit24[t.on8][t.oncursor]
			eof := t.forward(len(t.limit24[t.on8]))
			return reverse24b(v), int(v[3]), eof
		case 3:
			v := t.limit32[t.on8][t.oncursor]
			eof := t.forward(len(t.limit32[t.on8]))
			return reverse32b(v), int(v[4]), eof
		case 4:
			v := t.limit40[t.on8][t.oncursor]
			eof := t.forward(len(t.limit40[t.on8]))
			return reverse40b(v), int(v[5]), eof
		case 5:
			v := t.limit48[t.on8][t.oncursor]
			eof := t.forward(len(t.limit48[t.on8]))
			return reverse48b(v), int(v[6]), eof
		case 6:
			v := t.limit56[t.on8][t.oncursor]
			eof := t.forward(len(t.limit56[t.on8]))
			return reverse56b(v), int(v[7]), eof
		default:
			v := t.limit64[t.on8][t.oncursor]
			eof := t.forward(len(t.limit64[t.on8]))
			return reverse64b(v), int(v[8]), eof
	}
}

func (t *CounterBytes) Keys() [][]byte {

	var on, run int
	keys := make([][]byte, t.total)
	
	for run=0; run<8; run++ {
		for _, v := range t.limit8[run] {
			keys[on] = reverse8b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit16[run] {
			keys[on] = reverse16b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit24[run] {
			keys[on] = reverse24b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit32[run] {
			keys[on] = reverse32b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit40[run] {
			keys[on] = reverse40b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit48[run] {
			keys[on] = reverse48b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit56[run] {
			keys[on] = reverse56b(v)
			on++
		}
	}
	for run=0; run<8; run++ {
		for _, v := range t.limit64[run] {
			keys[on] = reverse64b(v)
			on++
		}
	}
	
	return keys
}


func (t *CounterBytes) Write(w *custom.Writer) {
	var run int

	// Write total
	w.WriteUint64Variable(uint64(t.total))
	
	// Write t.limit8
	for run=0; run<8; run++ {
		tmp := t.limit8[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
		}
	}
	// Write t.limit16
	for run=0; run<8; run++ {
		tmp := t.limit16[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
		}
	}
	// Write t.limit24
	for run=0; run<8; run++ {
		tmp := t.limit24[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
		}
	}
	// Write t.limit32
	for run=0; run<8; run++ {
		tmp := t.limit32[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
		}
	}
	// Write t.limit40
	for run=0; run<8; run++ {
		tmp := t.limit40[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
		}
	}
	// Write t.limit48
	for run=0; run<8; run++ {
		tmp := t.limit48[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
		}
	}
	// Write t.limit56
	for run=0; run<8; run++ {
		tmp := t.limit56[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
			w.WriteUint64(v[7])
		}
	}
	// Write t.limit64
	for run=0; run<8; run++ {
		tmp := t.limit64[run]
		w.WriteUint64Variable(uint64(len(tmp)))
		for _, v := range tmp {
			w.WriteUint64(v[0])
			w.WriteUint64(v[1])
			w.WriteUint64(v[2])
			w.WriteUint64(v[3])
			w.WriteUint64(v[4])
			w.WriteUint64(v[5])
			w.WriteUint64(v[6])
			w.WriteUint64(v[7])
			w.WriteUint64(v[8])
		}
	}
}

func (t *CounterBytes) Read(r *custom.Reader) {
	var run int
	var i, l, a, b, c, d, e, f, g, h, z uint64

	// Write total
	t.total = int(r.ReadUint64Variable())
	
	// Read t.limit8
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][2]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			tmp[i] = [2]uint64{a, b}
		}
		t.limit8[run] = tmp
	}
	// Read t.limit16
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][3]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			tmp[i] = [3]uint64{a, b, c}
		}
		t.limit16[run] = tmp
	}
	// Read t.limit24
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][4]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			tmp[i] = [4]uint64{a, b, c, d}
		}
		t.limit24[run] = tmp
	}
	// Read t.limit32
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][5]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			tmp[i] = [5]uint64{a, b, c, d, e}
		}
		t.limit32[run] = tmp
	}
	// Read t.limit40
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][6]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			tmp[i] = [6]uint64{a, b, c, d, e, f}
		}
		t.limit40[run] = tmp
	}
	// Read t.limit48
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][7]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			tmp[i] = [7]uint64{a, b, c, d, e, f, g}
		}
		t.limit48[run] = tmp
	}
	// Read t.limit56
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][8]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			h = r.ReadUint64()
			tmp[i] = [8]uint64{a, b, c, d, e, f, g, h}
		}
		t.limit56[run] = tmp
	}
	// Read t.limit64
	for run=0; run<8; run++ {
		l = r.ReadUint64Variable()
		tmp := make([][9]uint64, l)
		for i=0; i<l; i++ {
			a = r.ReadUint64()
			b = r.ReadUint64()
			c = r.ReadUint64()
			d = r.ReadUint64()
			e = r.ReadUint64()
			f = r.ReadUint64()
			g = r.ReadUint64()
			h = r.ReadUint64()
			z = r.ReadUint64()
			tmp[i] = [9]uint64{a, b, c, d, e, f, g, h, z}
		}
		t.limit64[run] = tmp
	}
}

// ====================== runes ======================
// ---------- KeyRunes ----------

/*
	KeyRunes simply wraps KeyBytes with a custom runes to bytes converter.
	It is therefore more efficient to use KeyBytes directly if you can.
	It is, however, more efficient to use KeyRunes than converting from runes to bytes with Go's encoding packages or via string.
*/

func runes2bytes(word []rune) []byte {
	// Count how many bytes are needed to represent the slice of runes
	var l, i int
	var r rune
	for _, r = range word {
		if r < 256 {
			l++
		} else {
			if r >= 65536 {
				l += 4
			} else {
				l += 3
			}
		}
	}
	newword := make([]byte, l)
	
	// If each unicode value will fit into one byte
	if len(word) == l {
		for i, r = range word {
			newword[i] = byte(r)
		}
		return newword
	}
	
	// It's a bit more fancy
	for _, r = range word {
		if r < 256 {
			newword[i] = byte(r)
			i++
		} else {
			if r >= 65536 {
				newword[i] = 3 // code points 2 & 3 represent 2-3 additional characters needed
				newword[i+1] = byte(r % 256)
				r /= 256
				newword[i+2] = byte(r % 256)
				r /= 256
				newword[i+3] = byte(r % 256)
				i += 4
			} else {
				newword[i] = 2
				newword[i+1] = byte(r % 256)
				r /= 256
				newword[i+2] = byte(r % 256)
				i += 3
			}
		}
	}
	return newword
}

func bytes2runes(word []byte) []rune {
	l := len(word)
	newword := make([]rune, l)
	var b byte
	var i, on int
	for i=0; i<l; i++ {
		b = word[i]
		switch b {
			case 2:
				newword[on] = (rune(word[i+2]) * 256) + rune(word[i+1])
				i += 2
			case 3:
				newword[on] = (rune(word[i+3]) * 65536) + (rune(word[i+2]) * 256) + rune(word[i+1])
				i += 3
			default:
				newword[on] = rune(b)
		}
		on++
	}
	return newword[0:on]
}

// Add this to any struct to make it binary searchable.
type KeyRunes struct {
 child KeyBytes
}

// Find returns the index based on the key.
func (t *KeyRunes) Find(thekey []rune) (int, bool) {
	return t.child.Find(runes2bytes(thekey))
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyRunes) Add(thekey []rune) (int, bool) {
	return t.child.Add(runes2bytes(thekey))
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyRunes) AddUnsorted(thekey []rune) error {
	return t.child.AddUnsorted(runes2bytes(thekey))
}

// AddAt adds this key to the index in this exact position, so it does not require later rebuilding.
func (t *KeyRunes) AddAt(thekey []rune, i int) error {
	return t.child.AddAt(runes2bytes(thekey), i)
}

func (t *KeyRunes) Build() ([]int, error) {
	return t.child.Build()
}

func (t *KeyRunes) Optimize() {
	t.child.Optimize()
}

func (t *KeyRunes) Len() int {
	return t.child.Len()
}

func (t *KeyRunes) Reset() bool {
	return t.child.Reset()
}

func (t *KeyRunes) Next() ([]rune, bool) {
	a, b := t.child.Next()
	return bytes2runes(a), b
}

func (t *KeyRunes) Keys() [][]rune {
	keys := t.child.Keys()
	newkeys := make([][]rune, len(keys))
	for i, v := range keys {
		newkeys[i] = bytes2runes(v)
	}
	return newkeys
}

func (t *KeyRunes) Write(w *custom.Writer) {
	t.child.Write(w)
}

func (t *KeyRunes) Read(r *custom.Reader) {
	t.child.Read(r)
}

// Add this to any struct to make it binary searchable.
type KeyValRunes struct {
 child KeyValBytes
}

// Find returns the index based on the key.
func (t *KeyValRunes) Find(thekey []rune) (int, bool) {
	return t.child.Find(runes2bytes(thekey))
}

func (t *KeyValRunes) Update(thekey []rune, fn func(int) int) bool {
	return t.child.Update(runes2bytes(thekey), fn)
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValRunes) Add(thekey []rune, theval int) bool {
	return t.child.Add(runes2bytes(thekey), theval)
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValRunes) AddUnsorted(thekey []rune, theval int) error {
	return t.child.AddUnsorted(runes2bytes(thekey), theval)
}

func (t *KeyValRunes) Build() {
	t.child.Build()
}

func (t *KeyValRunes) Optimize() {
	t.child.Optimize()
}

func (t *KeyValRunes) Len() int {
	return t.child.Len()
}

func (t *KeyValRunes) Reset() bool {
	return t.child.Reset()
}

func (t *KeyValRunes) Next() ([]rune, int, bool) {
	a, b, c := t.child.Next()
	return bytes2runes(a), b, c
}

func (t *KeyValRunes) Keys() [][]rune {
	keys := t.child.Keys()
	newkeys := make([][]rune, len(keys))
	for i, v := range keys {
		newkeys[i] = bytes2runes(v)
	}
	return newkeys
}

func (t *KeyValRunes) Write(w *custom.Writer) {
	t.child.Write(w)
}

func (t *KeyValRunes) Read(r *custom.Reader) {
	t.child.Read(r)
}

// Add this to any struct to make it binary searchable.
type CounterRunes struct {
 child CounterBytes
}

// Find returns the index based on the key.
func (t *CounterRunes) Find(thekey []rune) (int, bool) {
	return t.child.Find(runes2bytes(thekey))
}

func (t *CounterRunes) Update(thekey []rune, fn func(int) int) bool {
	return t.child.Update(runes2bytes(thekey), fn)
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *CounterRunes) Add(thekey []rune, theval int) error {
	return t.child.Add(runes2bytes(thekey), theval)
}

func (t *CounterRunes) Build() {
	t.child.Build()
}

func (t *CounterRunes) Optimize() {
	t.child.Optimize()
}

func (t *CounterRunes) Len() int {
	return t.child.Len()
}

func (t *CounterRunes) Reset() bool {
	return t.child.Reset()
}

func (t *CounterRunes) Next() ([]rune, int, bool) {
	a, b, c := t.child.Next()
	return bytes2runes(a), b, c
}

func (t *CounterRunes) Keys() [][]rune {
	keys := t.child.Keys()
	newkeys := make([][]rune, len(keys))
	for i, v := range keys {
		newkeys[i] = bytes2runes(v)
	}
	return newkeys
}

func (t *CounterRunes) Write(w *custom.Writer) {
	t.child.Write(w)
}

func (t *CounterRunes) Read(r *custom.Reader) {
	t.child.Read(r)
}

func (t *CounterRunes) KeyRunes() *KeyRunes {
	obj := new(KeyRunes)
	child := t.child.KeyBytes()
	obj.child = *child
	return obj
}

func (t *CounterRunes) KeyValRunes() *KeyValRunes {
	obj := new(KeyValRunes)
	child := t.child.KeyValBytes()
	obj.child = *child
	return obj
}

// ====================== uint64 ======================
// ---------- KeyUint64 ----------

// Add this to any struct to make it binary searchable.
type KeyUint64 struct {
 key []uint64
 cursor int
}

func (t *KeyUint64) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyUint64) Find(thekey uint64) (int, bool) {
	var min, at int
	var current uint64
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at]; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return at, true // found
			}
		}
	}
	return min, false // doesn't exist
}

// Add is equivalent to Find and then AddAt
func (t *KeyUint64) Add(thekey uint64) (int, bool) {
	i, ok := t.Find(thekey)
	if !ok {
		t.AddAt(thekey, i)
	}
	return i, ok
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyUint64) AddUnsorted(thekey uint64) {
	t.key = append(t.key, thekey)
	return
}

// AddAt adds this key to the index in this exact position, so it does not require later rebuilding.
func (t *KeyUint64) AddAt(thekey uint64, i int) {
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]uint64, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:i])
		copy(tmp[i+1:], cur[i:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[i+1:], cur[i:])
	}
	cur[i] = thekey
	t.key = cur
}

// Build sorts the keys and returns an array telling you how to sort the values, you must do this yourself.
func (t *KeyUint64) Build() []int {
	l := len(t.key)
	temp := make([]sortIntUint64.KeyVal, l)
	var i int
	var k uint64
	for i, k = range t.key {
		temp[i] = sortIntUint64.KeyVal{i, k}
	}
	sortIntUint64.Asc(temp)
	imap := make([]int, l)
	newkey := t.key
	for i, obj := range temp {
		imap[i] = obj.K
		newkey[i] = obj.V
	}
	return imap
}

func (t *KeyUint64) Optimize() {
	temp := make([]uint64, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyUint64) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyUint64) Next() (uint64, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v, true
	}
	return v, false
}

func (t *KeyUint64) Keys() []uint64 {
	return t.key
}

// ---------- KeyValUint64 ----------

// Add this to any struct to make it binary searchable.
type KeyValUint64 struct {
 key []sortIntUint64.KeyVal
 cursor int
}

func (t *KeyValUint64) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyValUint64) Find(thekey uint64) (int, bool) {
	var min, at int
	var current uint64
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function
func (t *KeyValUint64) Update(thekey uint64, fn func(int) int) bool {
	var min, at int
	var current uint64
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function
func (t *KeyValUint64) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// Add is equivalent to Find and then AddAt
func (t *KeyValUint64) Add(thekey uint64, theval int) bool {
	var min, at int
	var current uint64
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				if t.key[at].K != theval {
					t.key[at].K = theval
				}
				return true // found
			}
		}
	}
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]sortIntUint64.KeyVal, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:min])
		copy(tmp[min+1:], cur[min:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[min+1:], cur[min:])
	}
	cur[min] = sortIntUint64.KeyVal{theval, thekey}
	t.key = cur
	return false
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValUint64) AddUnsorted(thekey uint64, theval int) {
	t.key = append(t.key, sortIntUint64.KeyVal{theval, thekey})
	return
}

// Build sorts the keys and values.
func (t *KeyValUint64) Build() {
	sortIntUint64.Asc(t.key)
}

func (t *KeyValUint64) Optimize() {
	temp := make([]sortIntUint64.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyValUint64) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyValUint64) Next() (uint64, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *KeyValUint64) Keys() []uint64 {
	keys := make([]uint64, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ---------- CounterUint64 ----------

// Add this to any struct to make it binary searchable.
type CounterUint64 struct {
 key []sortIntUint64.KeyVal
 cursor int
}

func NewCounterUint64(ar []sortIntUint64.KeyVal) *CounterUint64 {
	return &CounterUint64{key: ar[0:0]}
}

func (t *CounterUint64) RawKey() []sortIntUint64.KeyVal {
	return t.key
}

func (t *CounterUint64) KeyUint64() *KeyUint64 {
	obj := new(KeyUint64)
	key := make([]uint64, len(t.key))
	for i, v := range t.key {
		key[i] = v.V
	}
	obj.key = key
	return obj
}

func (t *CounterUint64) KeyValUint64() *KeyValUint64 {
	obj := new(KeyValUint64)
	key := make([]sortIntUint64.KeyVal, len(t.key))
	copy(key, t.key)
	obj.key = key
	return obj
}

func (t *CounterUint64) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *CounterUint64) Find(thekey uint64) (int, bool) {
	var min, at int
	var current uint64
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function.
func (t *CounterUint64) Update(thekey uint64, fn func(int) int) bool {
	var min, at int
	var current uint64
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function.
func (t *CounterUint64) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *CounterUint64) Add(thekey uint64, theval int) {
	t.key = append(t.key, sortIntUint64.KeyVal{theval, thekey})
}

// Build sorts the keys and values.
func (t *CounterUint64) Build() {
	var temp = t.key
	if len(temp) == 0 {
		return
	}
	sortIntUint64.Asc(temp)
	this := t.key[0].V
	n := t.key[0].K
	var on int
	for _, k := range t.key[1:] {
		if k.V == this {
			n += k.K
		} else {
			temp[on] = sortIntUint64.KeyVal{n, this}
			on++
			this = k.V
			n = k.K
		}
	}
	temp[on] = sortIntUint64.KeyVal{n, this}
	t.key = temp[0:on+1]
}

func (t *CounterUint64) Optimize() {
	temp := make([]sortIntUint64.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *CounterUint64) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *CounterUint64) Next() (uint64, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *CounterUint64) Keys() []uint64 {
	keys := make([]uint64, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ------------- export ---------------

func (t *KeyUint64) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(v)
	}
}

func (t *KeyUint64) Read(r *custom.Reader) {
	l := int(r.ReadUint64Variable())
	tmp := make([]uint64, l)
	for i:=0; i<l; i++ {
		tmp[i] = r.ReadUint64Variable()
	}
	t.key = tmp
}

func (t *KeyValUint64) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint64Variable(v.V)
	}
}

func (t *KeyValUint64) Read(r *custom.Reader) {
	var k int
	var v uint64
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint64.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = r.ReadUint64Variable()
		tmp[i] = sortIntUint64.KeyVal{k, v}
	}
	t.key = tmp
}

func (t *CounterUint64) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint64Variable(v.V)
	}
}

func (t *CounterUint64) Read(r *custom.Reader) {
	var k int
	var v uint64
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint64.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = r.ReadUint64Variable()
		tmp[i] = sortIntUint64.KeyVal{k, v}
	}
	t.key = tmp
}

// ====================== uint32 ======================
// ---------- KeyUint32 ----------

// Add this to any struct to make it binary searchable.
type KeyUint32 struct {
 key []uint32
 cursor int
}

func (t *KeyUint32) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyUint32) Find(thekey uint32) (int, bool) {
	var min, at int
	var current uint32
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at]; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return at, true // found
			}
		}
	}
	return min, false // doesn't exist
}

// Add is equivalent to Find and then AddAt
func (t *KeyUint32) Add(thekey uint32) (int, bool) {
	i, ok := t.Find(thekey)
	if !ok {
		t.AddAt(thekey, i)
	}
	return i, ok
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyUint32) AddUnsorted(thekey uint32) {
	t.key = append(t.key, thekey)
	return
}

// AddAt adds this key to the index in this exact position, so it does not require later rebuilding.
func (t *KeyUint32) AddAt(thekey uint32, i int) {
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]uint32, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:i])
		copy(tmp[i+1:], cur[i:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[i+1:], cur[i:])
	}
	cur[i] = thekey
	t.key = cur
}

// Build sorts the keys and returns an array telling you how to sort the values, you must do this yourself.
func (t *KeyUint32) Build() []int {
	l := len(t.key)
	temp := make([]sortIntUint32.KeyVal, l)
	var i int
	var k uint32
	for i, k = range t.key {
		temp[i] = sortIntUint32.KeyVal{i, k}
	}
	sortIntUint32.Asc(temp)
	imap := make([]int, l)
	newkey := t.key
	for i, obj := range temp {
		imap[i] = obj.K
		newkey[i] = obj.V
	}
	return imap
}

func (t *KeyUint32) Optimize() {
	temp := make([]uint32, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyUint32) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyUint32) Next() (uint32, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v, true
	}
	return v, false
}

func (t *KeyUint32) Keys() []uint32 {
	return t.key
}

// ---------- KeyValUint32 ----------

// Add this to any struct to make it binary searchable.
type KeyValUint32 struct {
 key []sortIntUint32.KeyVal
 cursor int
}

func (t *KeyValUint32) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyValUint32) Find(thekey uint32) (int, bool) {
	var min, at int
	var current uint32
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function
func (t *KeyValUint32) Update(thekey uint32, fn func(int) int) bool {
	var min, at int
	var current uint32
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function
func (t *KeyValUint32) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// Add is equivalent to Find and then AddAt
func (t *KeyValUint32) Add(thekey uint32, theval int) bool {
	var min, at int
	var current uint32
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				if t.key[at].K != theval {
					t.key[at].K = theval
				}
				return true // found
			}
		}
	}
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]sortIntUint32.KeyVal, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:min])
		copy(tmp[min+1:], cur[min:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[min+1:], cur[min:])
	}
	cur[min] = sortIntUint32.KeyVal{theval, thekey}
	t.key = cur
	return false
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValUint32) AddUnsorted(thekey uint32, theval int) {
	t.key = append(t.key, sortIntUint32.KeyVal{theval, thekey})
	return
}

// Build sorts the keys and values.
func (t *KeyValUint32) Build() {
	sortIntUint32.Asc(t.key)
}

func (t *KeyValUint32) Optimize() {
	temp := make([]sortIntUint32.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyValUint32) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyValUint32) Next() (uint32, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *KeyValUint32) Keys() []uint32 {
	keys := make([]uint32, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ---------- CounterUint32 ----------

// Add this to any struct to make it binary searchable.
type CounterUint32 struct {
 key []sortIntUint32.KeyVal
 cursor int
}

func NewCounterUint32(ar []sortIntUint32.KeyVal) *CounterUint32 {
	return &CounterUint32{key: ar[0:0]}
}

func (t *CounterUint32) RawKey() []sortIntUint32.KeyVal {
	return t.key
}

func (t *CounterUint32) KeyUint32() *KeyUint32 {
	obj := new(KeyUint32)
	key := make([]uint32, len(t.key))
	for i, v := range t.key {
		key[i] = v.V
	}
	obj.key = key
	return obj
}

func (t *CounterUint32) KeyValUint32() *KeyValUint32 {
	obj := new(KeyValUint32)
	key := make([]sortIntUint32.KeyVal, len(t.key))
	copy(key, t.key)
	obj.key = key
	return obj
}

func (t *CounterUint32) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *CounterUint32) Find(thekey uint32) (int, bool) {
	var min, at int
	var current uint32
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function.
func (t *CounterUint32) Update(thekey uint32, fn func(int) int) bool {
	var min, at int
	var current uint32
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function.
func (t *CounterUint32) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *CounterUint32) Add(thekey uint32, theval int) {
	t.key = append(t.key, sortIntUint32.KeyVal{theval, thekey})
}

// Build sorts the keys and values.
func (t *CounterUint32) Build() {
	var temp = t.key
	if len(temp) == 0 {
		return
	}
	sortIntUint32.Asc(temp)
	this := t.key[0].V
	n := t.key[0].K
	var on int
	for _, k := range t.key[1:] {
		if k.V == this {
			n += k.K
		} else {
			temp[on] = sortIntUint32.KeyVal{n, this}
			on++
			this = k.V
			n = k.K
		}
	}
	temp[on] = sortIntUint32.KeyVal{n, this}
	t.key = temp[0:on+1]
}

func (t *CounterUint32) Optimize() {
	temp := make([]sortIntUint32.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *CounterUint32) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *CounterUint32) Next() (uint32, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *CounterUint32) Keys() []uint32 {
	keys := make([]uint32, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ------------- export ---------------

func (t *KeyUint32) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v))
	}
}

func (t *KeyUint32) Read(r *custom.Reader) {
	l := int(r.ReadUint64Variable())
	tmp := make([]uint32, l)
	for i:=0; i<l; i++ {
		tmp[i] = uint32(r.ReadUint64Variable())
	}
	t.key = tmp
}

func (t *KeyValUint32) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint64Variable(uint64(v.V))
	}
}

func (t *KeyValUint32) Read(r *custom.Reader) {
	var k int
	var v uint32
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint32.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = uint32(r.ReadUint64Variable())
		tmp[i] = sortIntUint32.KeyVal{k, v}
	}
	t.key = tmp
}

func (t *CounterUint32) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint64Variable(uint64(v.V))
	}
}

func (t *CounterUint32) Read(r *custom.Reader) {
	var k int
	var v uint32
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint32.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = uint32(r.ReadUint64Variable())
		tmp[i] = sortIntUint32.KeyVal{k, v}
	}
	t.key = tmp
}

// ====================== uint16 ======================
// ---------- KeyUint16 ----------

// Add this to any struct to make it binary searchable.
type KeyUint16 struct {
 key []uint16
 cursor int
}

func (t *KeyUint16) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyUint16) Find(thekey uint16) (int, bool) {
	var min, at int
	var current uint16
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at]; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return at, true // found
			}
		}
	}
	return min, false // doesn't exist
}

// Add is equivalent to Find and then AddAt
func (t *KeyUint16) Add(thekey uint16) (int, bool) {
	i, ok := t.Find(thekey)
	if !ok {
		t.AddAt(thekey, i)
	}
	return i, ok
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyUint16) AddUnsorted(thekey uint16) {
	t.key = append(t.key, thekey)
	return
}

// AddAt adds this key to the index in this exact position, so it does not require later rebuilding.
func (t *KeyUint16) AddAt(thekey uint16, i int) {
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]uint16, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:i])
		copy(tmp[i+1:], cur[i:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[i+1:], cur[i:])
	}
	cur[i] = thekey
	t.key = cur
}

// Build sorts the keys and returns an array telling you how to sort the values, you must do this yourself.
func (t *KeyUint16) Build() []int {
	l := len(t.key)
	temp := make([]sortIntUint16.KeyVal, l)
	var i int
	var k uint16
	for i, k = range t.key {
		temp[i] = sortIntUint16.KeyVal{i, k}
	}
	sortIntUint16.Asc(temp)
	imap := make([]int, l)
	newkey := t.key
	for i, obj := range temp {
		imap[i] = obj.K
		newkey[i] = obj.V
	}
	return imap
}

func (t *KeyUint16) Optimize() {
	temp := make([]uint16, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyUint16) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyUint16) Next() (uint16, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v, true
	}
	return v, false
}

func (t *KeyUint16) Keys() []uint16 {
	return t.key
}

// ---------- KeyValUint16 ----------

// Add this to any struct to make it binary searchable.
type KeyValUint16 struct {
 key []sortIntUint16.KeyVal
 cursor int
}

func (t *KeyValUint16) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyValUint16) Find(thekey uint16) (int, bool) {
	var min, at int
	var current uint16
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function
func (t *KeyValUint16) Update(thekey uint16, fn func(int) int) bool {
	var min, at int
	var current uint16
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function
func (t *KeyValUint16) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// Add is equivalent to Find and then AddAt
func (t *KeyValUint16) Add(thekey uint16, theval int) bool {
	var min, at int
	var current uint16
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				if t.key[at].K != theval {
					t.key[at].K = theval
				}
				return true // found
			}
		}
	}
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]sortIntUint16.KeyVal, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:min])
		copy(tmp[min+1:], cur[min:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[min+1:], cur[min:])
	}
	cur[min] = sortIntUint16.KeyVal{theval, thekey}
	t.key = cur
	return false
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValUint16) AddUnsorted(thekey uint16, theval int) {
	t.key = append(t.key, sortIntUint16.KeyVal{theval, thekey})
	return
}

// Build sorts the keys and values.
func (t *KeyValUint16) Build() {
	sortIntUint16.Asc(t.key)
}

func (t *KeyValUint16) Optimize() {
	temp := make([]sortIntUint16.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyValUint16) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyValUint16) Next() (uint16, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *KeyValUint16) Keys() []uint16 {
	keys := make([]uint16, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ---------- CounterUint16 ----------

// Add this to any struct to make it binary searchable.
type CounterUint16 struct {
 key []sortIntUint16.KeyVal
 cursor int
}

func NewCounterUint16(ar []sortIntUint16.KeyVal) *CounterUint16 {
	return &CounterUint16{key: ar[0:0]}
}

func (t *CounterUint16) RawKey() []sortIntUint16.KeyVal {
	return t.key
}

func (t *CounterUint16) KeyUint16() *KeyUint16 {
	obj := new(KeyUint16)
	key := make([]uint16, len(t.key))
	for i, v := range t.key {
		key[i] = v.V
	}
	obj.key = key
	return obj
}

func (t *CounterUint16) KeyValUint16() *KeyValUint16 {
	obj := new(KeyValUint16)
	key := make([]sortIntUint16.KeyVal, len(t.key))
	copy(key, t.key)
	obj.key = key
	return obj
}

func (t *CounterUint16) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *CounterUint16) Find(thekey uint16) (int, bool) {
	var min, at int
	var current uint16
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function.
func (t *CounterUint16) Update(thekey uint16, fn func(int) int) bool {
	var min, at int
	var current uint16
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function.
func (t *CounterUint16) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *CounterUint16) Add(thekey uint16, theval int) {
	t.key = append(t.key, sortIntUint16.KeyVal{theval, thekey})
}

// Build sorts the keys and values.
func (t *CounterUint16) Build() {
	var temp = t.key
	if len(temp) == 0 {
		return
	}
	sortIntUint16.Asc(temp)
	this := t.key[0].V
	n := t.key[0].K
	var on int
	for _, k := range t.key[1:] {
		if k.V == this {
			n += k.K
		} else {
			temp[on] = sortIntUint16.KeyVal{n, this}
			on++
			this = k.V
			n = k.K
		}
	}
	temp[on] = sortIntUint16.KeyVal{n, this}
	t.key = temp[0:on+1]
}

func (t *CounterUint16) Optimize() {
	temp := make([]sortIntUint16.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *CounterUint16) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *CounterUint16) Next() (uint16, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *CounterUint16) Keys() []uint16 {
	keys := make([]uint16, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ------------- export ---------------

func (t *KeyUint16) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint16(v)
	}
}

func (t *KeyUint16) Read(r *custom.Reader) {
	l := int(r.ReadUint64Variable())
	tmp := make([]uint16, l)
	for i:=0; i<l; i++ {
		tmp[i] = r.ReadUint16()
	}
	t.key = tmp
}

func (t *KeyValUint16) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint16(v.V)
	}
}

func (t *KeyValUint16) Read(r *custom.Reader) {
	var k int
	var v uint16
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint16.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = r.ReadUint16()
		tmp[i] = sortIntUint16.KeyVal{k, v}
	}
	t.key = tmp
}

func (t *CounterUint16) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint16(v.V)
	}
}

func (t *CounterUint16) Read(r *custom.Reader) {
	var k int
	var v uint16
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint16.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = r.ReadUint16()
		tmp[i] = sortIntUint16.KeyVal{k, v}
	}
	t.key = tmp
}

// ====================== uint8 ======================
// ---------- KeyUint8 ----------

// Add this to any struct to make it binary searchable.
type KeyUint8 struct {
 key []uint8
 cursor int
}

func (t *KeyUint8) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyUint8) Find(thekey uint8) (int, bool) {
	var min, at int
	var current uint8
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at]; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return at, true // found
			}
		}
	}
	return min, false // doesn't exist
}

// Add is equivalent to Find and then AddAt
func (t *KeyUint8) Add(thekey uint8) (int, bool) {
	i, ok := t.Find(thekey)
	if !ok {
		t.AddAt(thekey, i)
	}
	return i, ok
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyUint8) AddUnsorted(thekey uint8) {
	t.key = append(t.key, thekey)
	return
}

// AddAt adds this key to the index in this exact position, so it does not require later rebuilding.
func (t *KeyUint8) AddAt(thekey uint8, i int) {
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]uint8, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:i])
		copy(tmp[i+1:], cur[i:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[i+1:], cur[i:])
	}
	cur[i] = thekey
	t.key = cur
}

// Build sorts the keys and returns an array telling you how to sort the values, you must do this yourself.
func (t *KeyUint8) Build() []int {
	l := len(t.key)
	temp := make([]sortIntUint8.KeyVal, l)
	var i int
	var k uint8
	for i, k = range t.key {
		temp[i] = sortIntUint8.KeyVal{i, k}
	}
	sortIntUint8.Asc(temp)
	imap := make([]int, l)
	newkey := t.key
	for i, obj := range temp {
		imap[i] = obj.K
		newkey[i] = obj.V
	}
	return imap
}

func (t *KeyUint8) Optimize() {
	temp := make([]uint8, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyUint8) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyUint8) Next() (uint8, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v, true
	}
	return v, false
}

func (t *KeyUint8) Keys() []uint8 {
	return t.key
}

// ---------- KeyValUint8 ----------

// Add this to any struct to make it binary searchable.
type KeyValUint8 struct {
 key []sortIntUint8.KeyVal
 cursor int
}

func (t *KeyValUint8) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyValUint8) Find(thekey uint8) (int, bool) {
	var min, at int
	var current uint8
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function
func (t *KeyValUint8) Update(thekey uint8, fn func(int) int) bool {
	var min, at int
	var current uint8
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function
func (t *KeyValUint8) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// Add is equivalent to Find and then AddAt
func (t *KeyValUint8) Add(thekey uint8, theval int) bool {
	var min, at int
	var current uint8
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				if t.key[at].K != theval {
					t.key[at].K = theval
				}
				return true // found
			}
		}
	}
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]sortIntUint8.KeyVal, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:min])
		copy(tmp[min+1:], cur[min:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[min+1:], cur[min:])
	}
	cur[min] = sortIntUint8.KeyVal{theval, thekey}
	t.key = cur
	return false
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValUint8) AddUnsorted(thekey uint8, theval int) {
	t.key = append(t.key, sortIntUint8.KeyVal{theval, thekey})
	return
}

// Build sorts the keys and values.
func (t *KeyValUint8) Build() {
	sortIntUint8.Asc(t.key)
}

func (t *KeyValUint8) Optimize() {
	temp := make([]sortIntUint8.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyValUint8) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyValUint8) Next() (uint8, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *KeyValUint8) Keys() []uint8 {
	keys := make([]uint8, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ---------- CounterUint8 ----------

// Add this to any struct to make it binary searchable.
type CounterUint8 struct {
 key []sortIntUint8.KeyVal
 cursor int
}

func NewCounterUint8(ar []sortIntUint8.KeyVal) *CounterUint8 {
	return &CounterUint8{key: ar[0:0]}
}

func (t *CounterUint8) RawKey() []sortIntUint8.KeyVal {
	return t.key
}

func (t *CounterUint8) KeyUint8() *KeyUint8 {
	obj := new(KeyUint8)
	key := make([]uint8, len(t.key))
	for i, v := range t.key {
		key[i] = v.V
	}
	obj.key = key
	return obj
}

func (t *CounterUint8) KeyValUint8() *KeyValUint8 {
	obj := new(KeyValUint8)
	key := make([]sortIntUint8.KeyVal, len(t.key))
	copy(key, t.key)
	obj.key = key
	return obj
}

func (t *CounterUint8) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *CounterUint8) Find(thekey uint8) (int, bool) {
	var min, at int
	var current uint8
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function.
func (t *CounterUint8) Update(thekey uint8, fn func(int) int) bool {
	var min, at int
	var current uint8
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function.
func (t *CounterUint8) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *CounterUint8) Add(thekey uint8, theval int) {
	t.key = append(t.key, sortIntUint8.KeyVal{theval, thekey})
}

// Build sorts the keys and values.
func (t *CounterUint8) Build() {
	var temp = t.key
	if len(temp) == 0 {
		return
	}
	sortIntUint8.Asc(temp)
	this := t.key[0].V
	n := t.key[0].K
	var on int
	for _, k := range t.key[1:] {
		if k.V == this {
			n += k.K
		} else {
			temp[on] = sortIntUint8.KeyVal{n, this}
			on++
			this = k.V
			n = k.K
		}
	}
	temp[on] = sortIntUint8.KeyVal{n, this}
	t.key = temp[0:on+1]
}

func (t *CounterUint8) Optimize() {
	temp := make([]sortIntUint8.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *CounterUint8) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *CounterUint8) Next() (uint8, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *CounterUint8) Keys() []uint8 {
	keys := make([]uint8, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ------------- export ---------------

func (t *KeyUint8) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteByte(v)
	}
}

func (t *KeyUint8) Read(r *custom.Reader) {
	l := int(r.ReadUint64Variable())
	tmp := make([]uint8, l)
	for i:=0; i<l; i++ {
		tmp[i] = r.ReadByte()
	}
	t.key = tmp
}

func (t *KeyValUint8) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteByte(v.V)
	}
}

func (t *KeyValUint8) Read(r *custom.Reader) {
	var k int
	var v uint8
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint8.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = r.ReadByte()
		tmp[i] = sortIntUint8.KeyVal{k, v}
	}
	t.key = tmp
}

func (t *CounterUint8) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteByte(v.V)
	}
}

func (t *CounterUint8) Read(r *custom.Reader) {
	var k int
	var v uint8
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntUint8.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = r.ReadByte()
		tmp[i] = sortIntUint8.KeyVal{k, v}
	}
	t.key = tmp
}

// ====================== int ======================
// ---------- KeyInt ----------

// Add this to any struct to make it binary searchable.
type KeyInt struct {
 key []int
 cursor int
}

func (t *KeyInt) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyInt) Find(thekey int) (int, bool) {
	var min, at int
	var current int
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at]; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return at, true // found
			}
		}
	}
	return min, false // doesn't exist
}

// Add is equivalent to Find and then AddAt
func (t *KeyInt) Add(thekey int) (int, bool) {
	i, ok := t.Find(thekey)
	if !ok {
		t.AddAt(thekey, i)
	}
	return i, ok
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyInt) AddUnsorted(thekey int) {
	t.key = append(t.key, thekey)
	return
}

// AddAt adds this key to the index in this exact position, so it does not require later rebuilding.
func (t *KeyInt) AddAt(thekey int, i int) {
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]int, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:i])
		copy(tmp[i+1:], cur[i:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[i+1:], cur[i:])
	}
	cur[i] = thekey
	t.key = cur
}

// Build sorts the keys and returns an array telling you how to sort the values, you must do this yourself.
func (t *KeyInt) Build() []int {
	l := len(t.key)
	temp := make([]sortIntInt.KeyVal, l)
	var i int
	var k int
	for i, k = range t.key {
		temp[i] = sortIntInt.KeyVal{i, k}
	}
	sortIntInt.Asc(temp)
	imap := make([]int, l)
	newkey := t.key
	for i, obj := range temp {
		imap[i] = obj.K
		newkey[i] = obj.V
	}
	return imap
}

func (t *KeyInt) Optimize() {
	temp := make([]int, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyInt) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyInt) Next() (int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v, true
	}
	return v, false
}

func (t *KeyInt) Keys() []int {
	return t.key
}

// ---------- KeyValInt ----------

// Add this to any struct to make it binary searchable.
type KeyValInt struct {
 key []sortIntInt.KeyVal
 cursor int
}

func (t *KeyValInt) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *KeyValInt) Find(thekey int) (int, bool) {
	var min, at int
	var current int
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function
func (t *KeyValInt) Update(thekey int, fn func(int) int) bool {
	var min, at int
	var current int
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function
func (t *KeyValInt) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// Add is equivalent to Find and then AddAt
func (t *KeyValInt) Add(thekey int, theval int) bool {
	var min, at int
	var current int
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				if t.key[at].K != theval {
					t.key[at].K = theval
				}
				return true // found
			}
		}
	}
	cur := t.key
	lc := len(cur)
	if lc == cap(cur) {
		tmp := make([]sortIntInt.KeyVal, lc + 1, (lc * 2) + 1)
		copy(tmp, cur[0:min])
		copy(tmp[min+1:], cur[min:])
		cur = tmp
	} else {
		cur = cur[0:lc+1]
		copy(cur[min+1:], cur[min:])
	}
	cur[min] = sortIntInt.KeyVal{theval, thekey}
	t.key = cur
	return false
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *KeyValInt) AddUnsorted(thekey int, theval int) {
	t.key = append(t.key, sortIntInt.KeyVal{theval, thekey})
	return
}

// Build sorts the keys and values.
func (t *KeyValInt) Build() {
	sortIntInt.Asc(t.key)
}

func (t *KeyValInt) Optimize() {
	temp := make([]sortIntInt.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *KeyValInt) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *KeyValInt) Next() (int, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *KeyValInt) Keys() []int {
	keys := make([]int, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ---------- CounterInt ----------

// Add this to any struct to make it binary searchable.
type CounterInt struct {
 key []sortIntInt.KeyVal
 cursor int
}

func NewCounterInt(ar []sortIntInt.KeyVal) *CounterInt {
	return &CounterInt{key: ar[0:0]}
}

func (t *CounterInt) RawKey() []sortIntInt.KeyVal {
	return t.key
}

func (t *CounterInt) KeyInt() *KeyInt {
	obj := new(KeyInt)
	key := make([]int, len(t.key))
	for i, v := range t.key {
		key[i] = v.V
	}
	obj.key = key
	return obj
}

func (t *CounterInt) KeyValInt() *KeyValInt {
	obj := new(KeyValInt)
	key := make([]sortIntInt.KeyVal, len(t.key))
	copy(key, t.key)
	obj.key = key
	return obj
}

func (t *CounterInt) Len() int {
	return len(t.key)
}

// Find returns the index based on the key.
func (t *CounterInt) Find(thekey int) (int, bool) {
	var min, at int
	var current int
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				return t.key[at].K, true // found
			}
		}
	}
	return 0, false // doesn't exist
}

// Modifies the value of the key by running it through the provided function.
func (t *CounterInt) Update(thekey int, fn func(int) int) bool {
	var min, at int
	var current int
	max := len(t.key) - 1
	for min <= max {
		at = min + ((max - min) / 2)
		if current=t.key[at].V; thekey < current {
			max = at - 1
		} else {
		if thekey > current {
			min = at + 1
			} else {
				t.key[at].K = fn(t.key[at].K)
				return true // found
			}
		}
	}
	return false // doesn't exist
}

// Modifies all values by running each through the provided function.
func (t *CounterInt) UpdateAll(fn func(int) int) {
	tmp := t.key
	l := len(tmp)
	for i:=0; i<l; i++ {
		tmp[i].K = fn(tmp[i].K)
	}
}

// AddUnsorted adds this key to the end of the index for later building with Build.
func (t *CounterInt) Add(thekey int, theval int) {
	t.key = append(t.key, sortIntInt.KeyVal{theval, thekey})
}

// Build sorts the keys and values.
func (t *CounterInt) Build() {
	var temp = t.key
	if len(temp) == 0 {
		return
	}
	sortIntInt.Asc(temp)
	this := t.key[0].V
	n := t.key[0].K
	var on int
	for _, k := range t.key[1:] {
		if k.V == this {
			n += k.K
		} else {
			temp[on] = sortIntInt.KeyVal{n, this}
			on++
			this = k.V
			n = k.K
		}
	}
	temp[on] = sortIntInt.KeyVal{n, this}
	t.key = temp[0:on+1]
}

func (t *CounterInt) Optimize() {
	temp := make([]sortIntInt.KeyVal, len(t.key))
	copy(temp, t.key)
	t.key = temp
}

func (t *CounterInt) Reset() bool {
	t.cursor = 0
	if len(t.key) == 0 {
		return false
	}
	return true
}

func (t *CounterInt) Next() (int, int, bool) {
	v := t.key[t.cursor]
	if t.cursor++; t.cursor == len(t.key) {
		t.cursor = 0
		return v.V, v.K, true
	}
	return v.V, v.K, false
}

func (t *CounterInt) Keys() []int {
	keys := make([]int, len(t.key))
	for i, v := range t.key {
		keys[i] = v.V
	}
	return keys
}

// ------------- export ---------------

func (t *KeyInt) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v))
	}
}

func (t *KeyInt) Read(r *custom.Reader) {
	l := int(r.ReadUint64Variable())
	tmp := make([]int, l)
	for i:=0; i<l; i++ {
		tmp[i] = int(r.ReadUint64Variable())
	}
	t.key = tmp
}

func (t *KeyValInt) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint64Variable(uint64(v.V))
	}
}

func (t *KeyValInt) Read(r *custom.Reader) {
	var k int
	var v int
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntInt.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = int(r.ReadUint64Variable())
		tmp[i] = sortIntInt.KeyVal{k, v}
	}
	t.key = tmp
}

func (t *CounterInt) Write(w *custom.Writer) {
	w.WriteUint64Variable(uint64(len(t.key)))
	for _, v := range t.key {
		w.WriteUint64Variable(uint64(v.K))
		w.WriteUint64Variable(uint64(v.V))
	}
}

func (t *CounterInt) Read(r *custom.Reader) {
	var k int
	var v int
	l := int(r.ReadUint64Variable())
	tmp := make([]sortIntInt.KeyVal, l)
	for i:=0; i<l; i++ {
		k = int(r.ReadUint64Variable())
		v = int(r.ReadUint64Variable())
		tmp[i] = sortIntInt.KeyVal{k, v}
	}
	t.key = tmp
}
