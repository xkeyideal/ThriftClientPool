## Custom

Custom provides very optimized writers and readers, speeding up writing and reading to and from both disk and memory.

### Structures
- **custom.Writer** wraps an io.Writer, optimizing the writes (replaces bufio.Writer)
- **custom.Reader** wraps an io.Reader, optimizing the reads (replaces bufio.Reader)
- **custom.Buffer** replaces bytes.Buffer
- **custom.BytesReader** replaces bytes.Reader

### Features
- Highly optimized with focus on speed and efficiency for both disk and memory applications
- Buffers are pooled and reused so do not need to be allocated more than once across the runtime of the application
- Read and writes to the underlying reader/writer are buffered, improving the read/write speed
- Optimized functions with encoding & decoding for writing/reading slices, strings, integers, floats and booleans
- Built-in support for zlib and snappy compression
- Satisfies io.Reader, io.ReadCloser, io.ReadSeeker, io.RuneReader, io.Writer, io.WriteCloser, io.WriteSeeker

### Documentation
[View documentation on Godoc](https://godoc.org/github.com/AlasdairF/Custom)

### Writing Example
     fi, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
     if err != nil {
        panic(err)
     }
     defer fi.Close()
     w := custom.NewZlibWriter(fi)
     defer w.Close()
     w.WriteUint32(uint32(len(mydata)))
     for _, b := range mydata {
        w.WriteUint64Variable(b)
     }

### Reading Example
     fi, err := os.Open(filename)
     if err != nil {
        panic(err)
     }
     defer fi.Close()
     r := custom.NewZlibReader(fi)
     defer r.Close()
     l := int(r.ReadUint32())
     mydata := make([]uint64, l)
     for i:=0; i<l; i++ {
        mydata[i] = r.ReadUint64Variable()
     }
     
