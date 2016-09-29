##BinSearch

BinSearch is a super-efficient, in-memory key/value data structure for Go. In future it could also expand to be disk-based easily enough.

##Features
* Supports keys in the following types: `[]byte`, `[]rune`, `int`, `uint64`, `uint32`, `uint16`, `uint8`.
* Supports the following data structures: Key/Index store, Key/Val store, Counter (Accumulator).
* Key/Index store allows for any value structure to be used along with the key.
* Includes Read and Write functions for reading and writing the structure to disk.
* Backend is binary search with a great number of optimizations.
* Written with focus on high speed and low memory footprint.

##Advantages
* Incredibly memory efficient, max 5KB of memory overhead.
* `[]byte` and `[]rune` keys are compacted and therefore use *less* memory than the original slice, while remaining lossless and retrievable.
* Very, very fast. Much faster than the native `map` for almost all applications.

##Disadvantages
* Slow for sequential: find, add key, find, add key... in those cases the native `map` structure is faster. However, this problem can usually be solved by using Counter and Key/Index stores together, which restores the superiority to BinSearch for both speed and memory efficiency.
* Maximum key size for `[]byte` and `[]rune` types is 64 bytes.

##Installation

    go get github.com/AlasdairF/BinSearch
	
##Usage

The different structure names are one of `Key`, `KeyVal`, `Counter`, followed by one of `Bytes`, `Runes`, `Int`, `Uint64`, `Uint32`, `Uint16`, `Uint8`. E.g. `KeyValBytes`, `CounterUint32`.

`Key` and `KeyVal` types should never have duplicate keys added. It is important to either use `Find(key)` to check if a key exists before adding it (see Example 4), or to use the `Counter` structure to remove duplicates first (see Example 10).

Most structures require the keys to be added first and then the `Build()` function executed before any `Find(key)` is performed. The exception to this is the `Add(key)` function from `Key` and `KeyVal` types, which does not require the use of `Build()` and which does allow for `Find(key)` to be performed at any time, but the insertion of the keys is considerably slower. In most cases this is not necessary and there is usually a way to avoid using `Add(key)`.

###Key Type

The Key store records only the keys and its `Find(key)` function will return an index number, which is the position of the key. You can store any value in a separate slice under this index. The Key store is also useful if you don't need any values, for example you want to see only if a word exists in a dictionary.

###KeyVal Type

The KeyVal store is similar to the Key store but it also features an `int` value associated with each key. This value is always returned by the `Find(key)` function instead of the index. If you do not want an `int` value then use the Key store and implement your own value slice.

###Counter Type

The Counter type adds up all of the values associated with every identical key element. It is therefore useful for removing duplicates from a list, for tallying up scores, and for counting the number of occurances of each key. Values are `int` and so may be positive or negative; the value is irrelevant if `Counter` is being used to remove duplicates.	
	
##Index

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

##Examples

###1. Basic KeyVal usage with AddUnsorted & Build (fast)

	obj := new(binsearch.KeyValBytes) // create BinSearch structure
	obj.AddUnsorted([]byte("something"), 277) // add at the end, Find() will not yet work
	obj.AddUnsorted([]byte("word"), 1000)
	obj.AddUnsorted([]byte("hello"), 55)
	obj.Build() // when using AddUnsorted, Build() must be executed before Find()
	if val, exists := obj.Find([]byte("word")); exists {
		fmt.Println(`word is`, val) // word is 1000
	}

###2. Basic KeyVal usage with Add (slow)

	obj := new(binsearch.KeyValBytes) // create BinSearch structure
	obj.Add([]byte("something"), 277) // Add() does the add in the correct position so Find() can be used already
	obj.Add([]byte("word"), 1000)
	if val, exists := obj.Find([]byte("word")); exists { // Find() can be used without Build() because we used Add()
		fmt.Println(`word is`, val) // word is 1000
	}
	obj.Add([]byte("hello"), 55)
	if val, exists := obj.Find([]byte("hello")); exists {
		fmt.Println(`hello is`, val) // hello is 55
	}

###3. Using a custom value with AddUnsorted & Build (fast)

	obj := new(binsearch.KeyBytes)
	valstore := make([]mystruct, 0, 100)
	obj.AddUnsorted([]byte("first")) // Add the key to the end of the Key structure
	valstore = append(valstore, mystruct{123}) // Add the value using append, to the end of the value structure
	obj.AddUnsorted([]byte("second"))
	valstore = append(valstore, mystruct{456})
	newindexes, err := obj.Build() // only with Key stores, Build() returns a map of old indexes to new indexes to reorder the 
	temp := make([]mystruct, len(valstore)) // make a new slice for reordering the values
	for indx_new, indx_old := range newindexes { // loop through the mapping of old indexes to new indexes
		temp[indx_new] = valstore[indx_old] // copy the value from the old index location to the new index location on the new slice
	}
	valstore = temp // replace the old slice with the new sorted slice
	if indx, exists := obj.Find([]byte("second")); exists {
		val := valstore[indx]
		fmt.Println(val) // mystruct{456}
	}

###4. Using a custom value with Add (slow)

*Note: this can be done much more efficiently using a Counter. See Example 10 below.*

	key := []byte("test")
	obj := new(binsearch.KeyBytes)
	val := make([]mystruct, 0, 100)
	if indx, exists := obj.Add(key); exists {
		fmt.Println(`Value is`, val[indx])
	} else {
		val = append(val, mystruct{}) // Enlarge by 1
		copy(val[indx+1:], val[indx:]) // Make space at indx
		val[indx] = mystruct{123} // Add your value into the correct position
	}

###5. Using KeyBytes as a dictionary

	dictionary := new(binsearch.KeyBytes)
	dictionary.AddUnsorted([]byte("aardvark"))
	dictionary.AddUnsorted([]byte("pineapple"))
	dictionary.AddUnsorted([]byte("zulu"))
	dictionary.Build()
	if _, exists := dictionary.Find([]byte("pineapple")); exists {
		fmt.Println(`It's in the dictionary!`)
	} else {
		fmt.Println(`It's not in the dictionary!`)
	}
	
###6. Using CounterBytes to tally scores

	obj := new(binsearch.CounterBytes)
	obj.Add([]byte("john"), 5) // Counter structures only have the Add() function and it is always very fast
	obj.Add([]byte("fred"), 2)
	obj.Add([]byte("fred"), 4)
	obj.Add([]byte("bill"), 1)
	obj.Build()
	if val, exists := obj.Find([]byte("fred")); exists {
		fmt.Println(`Fred scored`, val) // Fred scored 6
	}
	
###7. Using CounterBytes to count occurances

	obj := new(binsearch.CounterBytes)
	obj.Add([]byte("the"), 1)
	obj.Add([]byte("the"), 1)
	obj.Add([]byte("the"), 1)
	obj.Add([]byte("and"), 1)
	obj.Build()
	if val, exists := obj.Find([]byte("the")); exists {
		fmt.Println(`the appears`, val, `times`) // the appears 3 times
	}
	
###8. Using CounterBytes to remove duplicates

	obj := new(binsearch.CounterBytes)
	obj.Add([]byte("the"), 0)
	obj.Add([]byte("the"), 0)
	obj.Add([]byte("the"), 0)
	obj.Add([]byte("and"), 0)
	obj.Build()
	uniques := obj.Keys()
	
###9. Creating a dictionary from results which may contain duplicates
	
	obj := new(binsearch.CounterBytes)
	obj.Add([]byte("aardvark"), 0)
	obj.Add([]byte("aardvark"), 0)
	obj.Add([]byte("zulu"), 0)
	obj.Add([]byte("zulu"), 0)
	obj.Build()
	dictionary := obj.Copy() // Copy() only exists for Counter structures, it copies the structure into a Key structure, this is only really necessary to save memory
	if _, exists := dictionary.Find([]byte("pineapple")); exists {
		fmt.Println(`It's in the dictionary!`)
	} else {
		fmt.Println(`It's not in the dictionary!`)
	}
	
###10. Creating a custom Key/Value structure from results which may contain duplicates

*Note: the result is equivalent to Example 4, but this is much faster.*

	origdata := [][]byte{[]byte("a"), []byte("b"), byte("c"), byte("c"), byte("d")}
	origval := []mystruct{mystruct{10}, mystruct{20}, mystruct{30}, mystruct{40}}
	obj := new(binsearch.CounterBytes)
	for _, word := range origdata {
		obj.Add(word, 0) // add all of the keys to the counter, ignoring the values
	}
	obj.Build()
	newobj := obj.Copy() // copy CounterBytes to KeyBytes
	valstore := make([]mystruct, newobj.Len()) // we already know the number of unique keys, so the size of the value storage structure is known
	var indx int
	for i, word := range origdata { // loop through all keys again
		indx, _ = newobj.Find(word) // find the index position of this key
		valstore[indx] = origval[i] // add the value into the appropriate position
	}

###11. Idiomatic way to range over the structure

	obj := new(binsearch.KeyValBytes)
	// ... pretend keys are added and Build() is executed here
	if obj.Reset() { // Reset() must be called first and it must be checked as Next() will panic if called on an empty structure
		var key []byte
		var val int
		for eof := false; !eof; {
			key, val, eof = obj.Next() // Bytes and Runes keys are converted back to the original `[]byte` or `[]rune` from their compacted version when using Next()
		}
	}

###12. Saving to file
	
	import "github.com/AlasdairF/Custom"
	func save(filename string, obj *binsearch.KeyValBytes) error {
		fi, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer fi.Close()
		w := custom.NewWriter(fi)
		defer w.Close()
		obj.Write(w)
		return nil
	}
	
###12. Reading from file
	
	import "github.com/AlasdairF/Custom"
	func load(filename string) (*binsearch.KeyValBytes, error) {
		// Open file for reading
		fi, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer fi.Close()
		// Attach reader
		r := custom.NewReader(fi, 20480)
		// Load binsearch.KeyRunes
		obj := new(binsearch.KeyValBytes)
		obj.Read(r) // do the reading
		// Make sure we're at the end and the checksum is OK
		if r.EOF() != nil {
			return nil, errors.New(`Not a valid binsearch structure.`)
		}
		return obj, nil
	}
