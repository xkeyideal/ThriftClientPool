package custom

import (
 "github.com/AlasdairF/Conv"
 "unicode/utf8"
 "math"
 "io"
 "os"
 "errors"
 "reflect"
 "sync"
 "github.com/klauspost/compress/zlib"
 "github.com/AlasdairF/snappy"
)

const (
	bufferLen = 65536 // determined in trials on writing to disk and writing to memory
	bufferLenMinus1  = bufferLen - 1
	bufferLenMinus2  = bufferLen - 2
	bufferLenMinus3  = bufferLen - 3
	bufferLenMinus4  = bufferLen - 4
	bufferLenMinus5  = bufferLen - 5
	bufferLenMinus6  = bufferLen - 6
	bufferLenMinus7  = bufferLen - 7
	bufferLenMinus8  = bufferLen - 8
	bufferLenMinus512 = bufferLen - 512
)

// Constants stolen from unicode/utf8 for WriteRune
const (
	maxRune   = '\U0010FFFF'
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000
	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111
	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1
)

var ErrNotEOF = errors.New(`Not EOF`)

// -------- INTERFACE --------

type Interface interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteByte(byte) error
	WriteRune(rune) (int, error)
	Write2Bytes(byte, byte) error
	Write3Bytes(byte, byte, byte) error
	Write4Bytes(byte, byte, byte, byte) error
	Write5Bytes(byte, byte, byte, byte, byte) error
	Write6Bytes(byte, byte, byte, byte, byte, byte) error
	Write7Bytes(byte, byte, byte, byte, byte, byte, byte) error
	Write8Bytes(byte, byte, byte, byte, byte, byte, byte, byte) error
	Write9Bytes(byte, byte, byte, byte, byte, byte, byte, byte, byte) error
	WriteBool(bool) error
	Write2Bools(bool, bool) error
	Write8Bools(bool, bool, bool, bool, bool, bool, bool, bool) error
	Write2Uint4s(uint8, uint8) error
	WriteUint16(uint16) error
	WriteUint16Variable(uint16) error
	WriteInt16Variable(int16) error
	WriteUint24(uint32) error
	WriteUint32(uint32) error
	WriteUint48(uint64) error
	WriteUint64(uint64) error
	WriteUint64Variable(uint64) error
	Write2Uint64sVariable(uint64, uint64) error
	WriteFloat32(float32) error
	WriteFloat64(float64) error
	WriteString8(string) (int, error)
	WriteString16(string) (int, error)
	WriteString32(string) (int, error)
	WriteBytes8([]byte) (int, error)
	WriteBytes16([]byte) (int, error)
	WriteBytes32([]byte) (int, error)
	WriteAll(a ...interface{}) (int, error)
	WriteInt(int) (int, error)
	Close() error
}

// -------- POOL -------

var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, bufferLen)
    },
}

// -------- COPY --------

func Copy(w io.Writer, r io.Reader) (t int, err error) {
	b := pool.Get().([]byte)
	defer pool.Put(b)
	var m, n int
	for {
		m, err = r.Read(b[n:])
		n += m
		if err == io.EOF {
			_, err = w.Write(b[0:n])
			t += n
			return
		}
		if err != nil {
			t += n
			return
		}
		if n >= bufferLenMinus512 {
			_, err = w.Write(b[0:n])
			t += n
			if err != nil {
				return
			}
			n = 0
		}
	}
}

func CopyFile(w io.Writer, filename string) (t int, err error) {
	var r *os.File
	r, err = os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()
	b := pool.Get().([]byte)
	defer pool.Put(b)
	var m, n int
	for {
		m, err = r.Read(b[n:])
		n += m
		if err == io.EOF {
			_, err = w.Write(b[0:n])
			t += n
			return
		}
		if err != nil {
			t += n
			return
		}
		if n >= bufferLenMinus512 {
			_, err = w.Write(b[0:n])
			t += n
			if err != nil {
				return
			}
			n = 0
		}
	}
}

// -------- FIXED BUFFER WRITER --------

type Writer struct {
	w io.Writer
	data []byte
	cursor int
	close bool
}

// Creates a new buffered writer wrapping an io.Writer
func NewWriter(f io.Writer) *Writer {
	if nf, ok := f.(*Writer); ok { // If it's already a Writer then don't create a new writer around it
		return nf
	} else {
		return &Writer{w: f, data: pool.Get().([]byte)}
	}
}

// Creates a new buffered writer wrapping an io.Writer which attempts to close the underlying writer when the custom.Writer is closed
func NewWriterCloser(f io.Writer) *Writer {
	return &Writer{w: f, data: pool.Get().([]byte), close: true}
}

// Creates a new buffered Zlib writer wrapping an io.Writer
func NewZlibWriter(f io.Writer) *Writer {
	return &Writer{w: zlib.NewWriter(f), data: pool.Get().([]byte), close: true}
}

// Creates a new buffered Snappy writer wrapping an io.Writer
func NewSnappyWriter(f io.Writer) *Writer {
	return &Writer{w: snappy.NewWriter(f), data: pool.Get().([]byte), close: true}
}

// Write a slice of bytes to the buffer. Implements io.Writer interface
func (w *Writer) Write(p []byte) (int, error) {
	l := len(p)
	if w.cursor + l > bufferLen {
		var err error
		if w.cursor > 0 {
			_, err = w.w.Write(w.data[0:w.cursor]) // flush
		}
		if l > bufferLen { // data to write is longer than the length of the Writer
			w.cursor = 0
			return w.w.Write(p)
		}
		copy(w.data[0:l], p)
		w.cursor = l
		return l, err
	}
	copy(w.data[w.cursor:], p)
	w.cursor += l
	return l, nil
}

// Write a string to the buffer
func (w *Writer) WriteString(p string) (int, error) {
	l := len(p)
	if w.cursor + l > bufferLen {
		var err error
		if w.cursor > 0 {
			_, err = w.w.Write(w.data[0:w.cursor]) // flush
		}
		if l > bufferLen { // data to write is longer than the length of the Writer
			w.cursor = 0
			return w.w.Write([]byte(p))
		}
		copy(w.data[0:l], p)
		w.cursor = l
		return l, err
	}
	copy(w.data[w.cursor:], p)
	w.cursor += l
	return l, nil
}

// Write a byte to the buffer
func (w *Writer) WriteByte(p byte) error {
	if w.cursor < bufferLen {
		w.data[w.cursor] = p
		w.cursor++
		return nil
	}
	var err error
	if w.cursor > 0 {
		_, err = w.w.Write(w.data[0:w.cursor]) // flush
	}
	w.data[0] = p
	w.cursor = 1
	return err
}

// Write a newline /n to the buffer
func (w *Writer) Writeln() error {
	if w.cursor < bufferLen {
		w.data[w.cursor] = '\n'
		w.cursor++
		return nil
	}
	var err error
	if w.cursor > 0 {
		_, err = w.w.Write(w.data[0:w.cursor]) // flush
	}
	w.data[0] = '\n'
	w.cursor = 1
	return err
}

// Write a rune to the buffer as UTF8
func (w *Writer) WriteRune(r rune) (int, error) {
	switch i := uint32(r); {
	case i <= rune1Max:
		err := w.WriteByte(byte(r))
		return 1, err
	case i <= rune2Max:
		err := w.Write2Bytes(t2 | byte(r>>6), tx | byte(r)&maskx)
		return 2, err
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = '\uFFFD'
		fallthrough
	case i <= rune3Max:
		err := w.Write3Bytes(t3 | byte(r>>12), tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 3, err
	default:
		err := w.Write4Bytes(t4 | byte(r>>18), tx | byte(r>>12)&maskx, tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 4, err
	}
}

// Write an integer in its ASCII form (not bitpacked)
func (w *Writer) WriteInt(p int) (int, error) {
	return conv.Write(w, p, 0)
}

// Write 2 bytes to the buffer
func (w *Writer) Write2Bytes(p1, p2 byte) error {
	if w.cursor < bufferLenMinus1 {
		w.data[w.cursor] = p1
		w.data[w.cursor + 1] = p2
		w.cursor += 2
		return nil
	}
	var err error
	if w.cursor > 0 {
		_, err = w.w.Write(w.data[0:w.cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.cursor = 2
	return err
}

// Write 3 bytes to the buffer
func (w *Writer) Write3Bytes(p1, p2, p3 byte) error {
	cursor := w.cursor
	if cursor < bufferLenMinus2 {
		w.data[cursor] = p1
		w.data[cursor + 1] = p2
		w.data[cursor + 2] = p3
		w.cursor += 3
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = w.w.Write(w.data[0:cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.data[2] = p3
	w.cursor = 3
	return err
}

// Write 4 bytes to the buffer
func (w *Writer) Write4Bytes(p1, p2, p3, p4 byte) error {
	cursor := w.cursor
	if cursor < bufferLenMinus3 {
		w.data[cursor] = p1
		w.data[cursor + 1] = p2
		w.data[cursor + 2] = p3
		w.data[cursor + 3] = p4
		w.cursor += 4
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = w.w.Write(w.data[0:cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.data[2] = p3
	w.data[3] = p4
	w.cursor = 4
	return err
}

// Write 5 bytes to the buffer
func (w *Writer) Write5Bytes(p1, p2, p3, p4, p5 byte) error {
	cursor := w.cursor
	if cursor < bufferLenMinus4 {
		w.data[cursor] = p1
		w.data[cursor + 1] = p2
		w.data[cursor + 2] = p3
		w.data[cursor + 3] = p4
		w.data[cursor + 4] = p5
		w.cursor += 5
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = w.w.Write(w.data[0:cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.data[2] = p3
	w.data[3] = p4
	w.data[4] = p5
	w.cursor = 5
	return err
}

// Write 6 bytes to the buffer
func (w *Writer) Write6Bytes(p1, p2, p3, p4, p5, p6 byte) error {
	cursor := w.cursor
	if cursor < bufferLenMinus5 {
		w.data[cursor] = p1
		w.data[cursor + 1] = p2
		w.data[cursor + 2] = p3
		w.data[cursor + 3] = p4
		w.data[cursor + 4] = p5
		w.data[cursor + 5] = p6
		w.cursor += 6
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = w.w.Write(w.data[0:cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.data[2] = p3
	w.data[3] = p4
	w.data[4] = p5
	w.data[5] = p6
	w.cursor = 6
	return err
}

// Write 7 bytes to the buffer
func (w *Writer) Write7Bytes(p1, p2, p3, p4, p5, p6, p7 byte) error {
	cursor := w.cursor
	if cursor < bufferLenMinus6 {
		w.data[cursor] = p1
		w.data[cursor + 1] = p2
		w.data[cursor + 2] = p3
		w.data[cursor + 3] = p4
		w.data[cursor + 4] = p5
		w.data[cursor + 5] = p6
		w.data[cursor + 6] = p7
		w.cursor += 7
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = w.w.Write(w.data[0:cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.data[2] = p3
	w.data[3] = p4
	w.data[4] = p5
	w.data[5] = p6
	w.data[6] = p7
	w.cursor = 7
	return err
}

// Write 8 bytes to the buffer
func (w *Writer) Write8Bytes(p1, p2, p3, p4, p5, p6, p7, p8 byte) error {
	cursor := w.cursor
	if cursor < bufferLenMinus7 {
		w.data[cursor] = p1
		w.data[cursor + 1] = p2
		w.data[cursor + 2] = p3
		w.data[cursor + 3] = p4
		w.data[cursor + 4] = p5
		w.data[cursor + 5] = p6
		w.data[cursor + 6] = p7
		w.data[cursor + 7] = p8
		w.cursor += 8
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = w.w.Write(w.data[0:cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.data[2] = p3
	w.data[3] = p4
	w.data[4] = p5
	w.data[5] = p6
	w.data[6] = p7
	w.data[7] = p8
	w.cursor = 8
	return err
}

// Write 9 bytes to the buffer
func (w *Writer) Write9Bytes(p1, p2, p3, p4, p5, p6, p7, p8, p9 byte) error {
	cursor := w.cursor
	if cursor < bufferLenMinus8 {
		w.data[cursor] = p1
		w.data[cursor + 1] = p2
		w.data[cursor + 2] = p3
		w.data[cursor + 3] = p4
		w.data[cursor + 4] = p5
		w.data[cursor + 5] = p6
		w.data[cursor + 6] = p7
		w.data[cursor + 7] = p8
		w.data[cursor + 8] = p9
		w.cursor += 9
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = w.w.Write(w.data[0:cursor]) // flush
	}
	w.data[0] = p1
	w.data[1] = p2
	w.data[2] = p3
	w.data[3] = p4
	w.data[4] = p5
	w.data[5] = p6
	w.data[6] = p7
	w.data[7] = p8
	w.data[8] = p9
	w.cursor = 9
	return err
}

// Encode a bool in 1 byte and write it to the buffer
func (w *Writer) WriteBool(v bool) error {
	if v {
		return w.WriteByte(1)
	} else {
		return w.WriteByte(0)
	}
}

// Encode 2 bools in 1 byte and write it to the buffer
func (w *Writer) Write2Bools(v1, v2 bool) error {
	var b byte
	if v1 {
		b = 1
	}
	if v2 {
		b |= 2
	}
	return w.WriteByte(b)
}

// Encode 8 bools in 1 byte and write it to the buffer
func (w *Writer) Write8Bools(v1, v2, v3, v4, v5, v6, v7, v8 bool) error {
	var b byte
	if v1 {
		b = 1
	}
	if v2 {
		b |= 2
	}
	if v3 {
		b |= 4
	}
	if v4 {
		b |= 8
	}
	if v5 {
		b |= 16
	}
	if v6 {
		b |= 32
	}
	if v7 {
		b |= 64
	}
	if v8 {
		b |= 128
	}
	return w.WriteByte(b)
}

// Encode 2 uint8s (assuming each is maximum 127) in 1 byte and write it to the buffer
func (w *Writer) Write2Uint4s(v1, v2 uint8) error {
	v1 |= v2 << 4
	return w.WriteByte(v1)
}

// Encode a uint16 in 2 bytes and write it to the buffer
func (w *Writer) WriteUint16(v uint16) error {
	return w.Write2Bytes(byte(v), byte(v >> 8))
}

// Encode a uint16 in 1-3 bytes and write it to the buffer. If the uint16 < 255 then it's encoded in 1 byte, otherwise 3 bytes.
func (w *Writer) WriteUint16Variable(v uint16) error {
	if v < 255 {
		return w.WriteByte(byte(v))
	}
	return w.Write3Bytes(255, byte(v), byte(v >> 8))
}

// Encode an int16 in 1-3 bytes and write it to the buffer. If the uint16 > -128 < 128 then it's encoded in 1 byte, otherwise 3 bytes.
func (w *Writer) WriteInt16Variable(v int16) error {
	if v > -128 && v < 128 {
		return w.WriteByte(byte(v + 127))
	}
	v2 := uint16(v)
	return w.Write3Bytes(255, byte(v2), byte(v2 >> 8))
}

// Encode a uint32 in 3 bytes and write it to the buffer
func (w *Writer) WriteUint24(v uint32) error {
	return w.Write3Bytes(byte(v), byte(v >> 8), byte(v >> 16))
}

// Encode a uint32 in 4 bytes and write it to the buffer
func (w *Writer) WriteUint32(v uint32) error {
	return w.Write4Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
}

// Encode a uint64 in 6 bytes and write it to the buffer
func (w *Writer) WriteUint48(v uint64) error {
	return w.Write6Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
}

// Encode a uint64 in 8 bytes and write it to the buffer
func (w *Writer) WriteUint64(v uint64) error {
	return w.Write8Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
}

// Encode a uint64 in 2-7 bytes and write it to the buffer. The length is always 1 byte more than the minimum representation of the uint64.
func (w *Writer) WriteUint64Variable(v uint64) error {
	switch numbytes(v) {
		case 0: return w.WriteByte(0)
		case 1: return w.Write2Bytes(1, byte(v))
		case 2: return w.Write3Bytes(2, byte(v), byte(v >> 8))
		case 3: return w.Write4Bytes(3, byte(v), byte(v >> 8), byte(v >> 16))
		case 4: return w.Write5Bytes(4, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
		case 5: return w.Write6Bytes(5, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32))
		case 6: return w.Write7Bytes(6, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
		case 7: return w.Write8Bytes(7, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48))
		case 8: return w.Write9Bytes(8, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 25), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
	}
	return nil
}

// Encode 2 uint64s in 3-13 bytes and write it to the buffer. The length is always 1 byte more than the minimum representation of both the uint64s.
func (w *Writer) Write2Uint64sVariable(v1 uint64, v2 uint64) error {
	s2 := numbytes(v2)
	switch numbytes(v1) {
		case 0: w.WriteByte(s2)
		case 1: w.Write2Bytes(16 | s2, byte(v1))
		case 2: w.Write3Bytes(32 | s2, byte(v1), byte(v1 >> 8))
		case 3: w.Write4Bytes(48 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16))
		case 4: w.Write5Bytes(64 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24))
		case 5: w.Write6Bytes(80 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32))
		case 6: w.Write7Bytes(96 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40))
		case 7: w.Write8Bytes(112 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48))
		case 8: w.Write9Bytes(128 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 25), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48), byte(v1 >> 56))
	}
	switch s2 {
		case 0: return nil
		case 1: return w.WriteByte(byte(v2))
		case 2: return w.Write2Bytes(byte(v2), byte(v2 >> 8))
		case 3: return w.Write3Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16))
		case 4: return w.Write4Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24))
		case 5: return w.Write5Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32))
		case 6: return w.Write6Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40))
		case 7: return w.Write7Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48))
		case 8: return w.Write8Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 25), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48), byte(v2 >> 56))
	}
	return nil
}

// Encode a float32 in 4 bytes and write it to the buffer
func (w *Writer) WriteFloat32(flt float32) error {
	return w.WriteUint32(math.Float32bits(flt))
}

// Encode a float64 in 8 bytes and write it to the buffer
func (w *Writer) WriteFloat64(flt float64) error {
	return w.WriteUint64(math.Float64bits(flt))
}

// Write a string to the buffer with maximum length 255
func (w *Writer) WriteString8(s string) (int, error) {
	if len(s) >= 255 {
		w.WriteByte(255)
		_, err := w.WriteString(s[0:255])
		return 256, err
	} else {
		w.WriteByte(uint8(len(s)))
		_, err := w.WriteString(s)
		return len(s) + 1, err
	}
}

// Write a string to the buffer with maximum length 65,535
func (w *Writer) WriteString16(s string) (int, error) {
	if len(s) >= 65535 {
		w.WriteUint16(65535)
		_, err := w.WriteString(s[0:65535])
		return 65537, err
	} else {
		w.WriteUint16(uint16(len(s)))
		_, err := w.WriteString(s)
		return len(s) + 2, err
	}
}

// Write a string to the buffer with maximum length 4,294,967,295
func (w *Writer) WriteString32(s string) (int, error) {
	if len(s) >= 4294967295 {
		w.WriteUint32(4294967295)
		_, err := w.WriteString(s[0:4294967295])
		return 4294967299, err
	} else {
		w.WriteUint32(uint32(len(s)))
		_, err := w.WriteString(s)
		return len(s) + 4, err
	}
}

// Write a slice of bytes to the buffer with maximum length 255
func (w *Writer) WriteBytes8(s []byte) (int, error) {
	if len(s) >= 255 {
		w.WriteByte(255)
		_, err := w.Write(s[0:255])
		return 256, err
	} else {
		w.WriteByte(uint8(len(s)))
		_, err := w.Write(s)
		return len(s) + 1, err
	}
}

// Write a slice of bytes to the buffer with maximum length 65,535
func (w *Writer) WriteBytes16(s []byte) (int, error) {
	if len(s) >= 65535 {
		w.WriteUint16(65535)
		_, err := w.Write(s[0:65535])
		return 65537, err
	} else {
		w.WriteUint16(uint16(len(s)))
		_, err := w.Write(s)
		return len(s) + 2, err
	}
}

// Write a slice of bytes to the buffer with maximum length 4,294,967,295
func (w *Writer) WriteBytes32(s []byte) (int, error) {
	if len(s) >= 4294967295 {
		w.WriteUint32(4294967295)
		_, err := w.Write(s[0:4294967295])
		return 4294967299, err
	} else {
		w.WriteUint32(uint32(len(s)))
		_, err := w.Write(s)
		return len(s) + 4, err
	}
}

// Reflects on the values and writes them all out. Not particularly safe.
// This function only works with: integer, slice of bytes, string, byte.
// uint8 is written as byte. int32 is written as rune. Other integers are written as ASCII representation of their number.
// A slice of anything other than bytes could cause unknown behavior.
func (w *Writer) WriteAll(a ...interface{}) (n int, err error) {
	var i int
	for _, p := range a {
		switch reflect.TypeOf(p).Kind() {
			case reflect.String:
				i, err = w.WriteString(reflect.ValueOf(p).String())
			case reflect.Slice: // all slices are assumed to be slices of bytes
				i, err = w.Write(reflect.ValueOf(p).Bytes())
			case reflect.Uint8: // byte
				i, err = 1, w.WriteByte(byte(reflect.ValueOf(p).Uint()))
			case reflect.Int32: // rune
				i, err = w.WriteRune(rune(reflect.ValueOf(p).Int()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
				i, err = conv.Write(w, int(reflect.ValueOf(p).Int()), 0)
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err = conv.Write(w, int(reflect.ValueOf(p).Uint()), 0)
			default:
				err = errors.New("custom.Buffer.WriteAll: not a supported type")
				return
		}
		if err != nil {
			return
		}
		n += i
	}
	return
}

// Flush the buffer and close the custom.Writer
func (w *Writer) Close() (err error) {
	if w.cursor > 0 {
		_, err = w.w.Write(w.data[0:w.cursor])
		w.cursor = 0
	}
	if len(w.data) == bufferLen {
		pool.Put(w.data)
		w.data = nil
	}
	if w.close {
		if sw, ok := w.w.(io.Closer); ok { // Attempt to close underlying writer if it has a Close() method
			if err == nil {
				err = sw.Close()
			} else {
				sw.Close()
			}
		}
	}
	w.w = nil
	return
}

// Flush the buffer of custom.Writer to the underlying io.Writer. This is not usually necessary as long as you remember to Close()
func (w *Writer) Flush() (err error) {
	if w.cursor > 0 {
		_, err = w.w.Write(w.data[0:w.cursor])
		w.cursor = 0
	}
	return
}

// Flushes the buffer to the underlying writer, closing it if this is a WriterCloser and then transfers to a new writer (no longer a WriterCloser)
func (w *Writer) Reset(newwriter io.Writer) (err error) {
	if w.cursor > 0 {
		_, err = w.w.Write(w.data[0:w.cursor])
		w.cursor = 0
	}
	if w.close {
		if sw, ok := w.w.(io.Closer); ok { // Attempt to close underlying writer if it has a Close() method
			if err == nil {
				err = sw.Close()
			} else {
				sw.Close()
			}
		}
		w.close = false
	}
	w.w = newwriter
	return
}

// -------- GROWING BUFFER --------

type Buffer struct {
	data []byte
	cursor, length int
}

// Creates a new buffer
func NewBuffer(l int) *Buffer {
	if l <= bufferLen {
		return &Buffer{data: pool.Get().([]byte), length: bufferLen}
	} else {
		return &Buffer{data: make([]byte, l), length: l}
	}
}

// Reads all of r into the buffer
func (w *Buffer) ReadFrom(r io.Reader) (int, error) {
	var i, n int
	var err error
	for {
		if w.cursor >= w.length - 512 {
			w.grow(1)
		}
		i, err = r.Read(w.data[w.cursor:])
		w.cursor += i
		n += i
		if err == io.EOF {
			break
		}
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (w *Buffer) grow(l int) {
	newLength := (w.length + l) * 2
	newAr := make([]byte, newLength)
	copy(newAr, w.data)
	if w.length == bufferLen {
		pool.Put(w.data)
	}
	w.length = newLength
	w.data = newAr
}

// Write a slice of bytes to the buffer. Implements io.Writer interface
func (w *Buffer) Write(p []byte) (int, error) {
	l := len(p)
	if w.cursor + l > w.length {
		w.grow(l)
	}
	copy(w.data[w.cursor:], p)
	w.cursor += l
	return l, nil
}

// Write a string to the buffer
func (w *Buffer) WriteString(p string) (int, error) {
	l := len(p)
	if w.cursor + l > w.length {
		w.grow(l)
	}
	copy(w.data[w.cursor:], p)
	w.cursor += l
	return l, nil
}

// Write a byte to the buffer
func (w *Buffer) WriteByte(p byte) error {
	if w.cursor >= w.length {
		w.grow(1)
	}
	w.data[w.cursor] = p
	w.cursor++
	return nil
}

// Write a line \n to the buffer
func (w *Buffer) Writeln() error {
	if w.cursor >= w.length {
		w.grow(1)
	}
	w.data[w.cursor] = '\n'
	w.cursor++
	return nil
}

// Write a rune to the buffer as UTF8
func (w *Buffer) WriteRune(r rune) (int, error) {
	switch i := uint32(r); {
	case i <= rune1Max:
		err := w.WriteByte(byte(r))
		return 1, err
	case i <= rune2Max:
		err := w.Write2Bytes(t2 | byte(r>>6), tx | byte(r)&maskx)
		return 2, err
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = '\uFFFD'
		fallthrough
	case i <= rune3Max:
		err := w.Write3Bytes(t3 | byte(r>>12), tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 3, err
	default:
		err := w.Write4Bytes(t4 | byte(r>>18), tx | byte(r>>12)&maskx, tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 4, err
	}
}

// Write an integer in its ASCII form (not bitpacked)
func (w *Buffer) WriteInt(p int) (int, error) {
	return conv.Write(w, p, 0)
}

// Write 2 bytes to the buffer
func (w *Buffer) Write2Bytes(p1, p2 byte) error {
	c := w.cursor
	if c + 2 > w.length {
		w.grow(2)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.cursor += 2
	return nil
}

// Write 3 bytes to the buffer
func (w *Buffer) Write3Bytes(p1, p2, p3 byte) error {
	c := w.cursor
	if c + 3 > w.length {
		w.grow(3)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.data[c + 2] = p3
	w.cursor += 3
	return nil
}

// Write 4 bytes to the buffer
func (w *Buffer) Write4Bytes(p1, p2, p3, p4 byte) error {
	c := w.cursor
	if c + 4 > w.length {
		w.grow(4)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.data[c + 2] = p3
	w.data[c + 3] = p4
	w.cursor += 4
	return nil
}

// Write 5 bytes to the buffer
func (w *Buffer) Write5Bytes(p1, p2, p3, p4, p5 byte) error {
	c := w.cursor
	if c + 5 > w.length {
		w.grow(5)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.data[c + 2] = p3
	w.data[c + 3] = p4
	w.data[c + 4] = p5
	w.cursor += 5
	return nil
}

// Write 6 bytes to the buffer
func (w *Buffer) Write6Bytes(p1, p2, p3, p4, p5, p6 byte) error {
	c := w.cursor
	if c + 6 > w.length {
		w.grow(6)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.data[c + 2] = p3
	w.data[c + 3] = p4
	w.data[c + 4] = p5
	w.data[c + 5] = p6
	w.cursor += 6
	return nil
}

// Write 7 bytes to the buffer
func (w *Buffer) Write7Bytes(p1, p2, p3, p4, p5, p6, p7 byte) error {
	c := w.cursor
	if c + 7 > w.length {
		w.grow(7)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.data[c + 2] = p3
	w.data[c + 3] = p4
	w.data[c + 4] = p5
	w.data[c + 5] = p6
	w.data[c + 6] = p7
	w.cursor += 7
	return nil
}

// Write 8 bytes to the buffer
func (w *Buffer) Write8Bytes(p1, p2, p3, p4, p5, p6, p7, p8 byte) error {
	c := w.cursor
	if c + 8 > w.length {
		w.grow(8)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.data[c + 2] = p3
	w.data[c + 3] = p4
	w.data[c + 4] = p5
	w.data[c + 5] = p6
	w.data[c + 6] = p7
	w.data[c + 7] = p8
	w.cursor += 8
	return nil
}

// Write 9 bytes to the buffer
func (w *Buffer) Write9Bytes(p1, p2, p3, p4, p5, p6, p7, p8, p9 byte) error {
	c := w.cursor
	if c + 9 > w.length {
		w.grow(9)
	}
	w.data[c] = p1
	w.data[c + 1] = p2
	w.data[c + 2] = p3
	w.data[c + 3] = p4
	w.data[c + 4] = p5
	w.data[c + 5] = p6
	w.data[c + 6] = p7
	w.data[c + 7] = p8
	w.data[c + 8] = p9
	w.cursor += 9
	return nil
}

// Encode a bool in 1 byte and write it to the buffer
func (w *Buffer) WriteBool(v bool) error {
	if v {
		return w.WriteByte(1)
	} else {
		return w.WriteByte(0)
	}
}

// Encode 2 bools in 1 byte and write it to the buffer
func (w *Buffer) Write2Bools(v1, v2 bool) error {
	var b byte
	if v1 {
		b = 1
	}
	if v2 {
		b |= 2
	}
	return w.WriteByte(b)
}

// Encode 8 bools in 1 byte and write it to the buffer
func (w *Buffer) Write8Bools(v1, v2, v3, v4, v5, v6, v7, v8 bool) error {
	var b byte
	if v1 {
		b = 1
	}
	if v2 {
		b |= 2
	}
	if v3 {
		b |= 4
	}
	if v4 {
		b |= 8
	}
	if v5 {
		b |= 16
	}
	if v6 {
		b |= 32
	}
	if v7 {
		b |= 64
	}
	if v8 {
		b |= 128
	}
	return w.WriteByte(b)
}

// Encode 2 uint8s (assuming each is maximum 127) in 1 byte and write it to the buffer
func (w *Buffer) Write2Uint4s(v1, v2 uint8) error {
	v1 |= v2 << 4
	return w.WriteByte(v1)
}

// Encode a uint16 in 2 bytes and write it to the buffer
func (w *Buffer) WriteUint16(v uint16) error {
	return w.Write2Bytes(byte(v), byte(v >> 8))
}

// Encode a uint16 in 1-3 bytes and write it to the buffer. If the uint16 < 255 then it's encoded in 1 byte, otherwise 3 bytes.
func (w *Buffer) WriteUint16Variable(v uint16) error {
	if v < 255 {
		return w.WriteByte(byte(v))
	}
	return w.Write3Bytes(255, byte(v), byte(v >> 8))
}

// Encode an int16 in 1-3 bytes and write it to the buffer. If the uint16 > -128 < 128 then it's encoded in 1 byte, otherwise 3 bytes.
func (w *Buffer) WriteInt16Variable(v int16) error {
	if v > -128 && v < 128 {
		return w.WriteByte(byte(v + 127))
	}
	v2 := uint16(v)
	return w.Write3Bytes(255, byte(v2), byte(v2 >> 8))
}

// Encode a uint32 in 3 bytes and write it to the buffer
func (w *Buffer) WriteUint24(v uint32) error {
	return w.Write3Bytes(byte(v), byte(v >> 8), byte(v >> 16))
}

// Encode a uint32 in 4 bytes and write it to the buffer
func (w *Buffer) WriteUint32(v uint32) error {
	return w.Write4Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
}

// Encode a uint64 in 6 bytes and write it to the buffer
func (w *Buffer) WriteUint48(v uint64) error {
	return w.Write6Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
}

// Encode a uint64 in 8 bytes and write it to the buffer
func (w *Buffer) WriteUint64(v uint64) error {
	return w.Write8Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
}

// Encode a uint64 in 2-7 bytes and write it to the buffer. The length is always 1 byte more than the minimum representation of the uint64.
func (w *Buffer) WriteUint64Variable(v uint64) error {
	switch numbytes(v) {
		case 0: return w.WriteByte(0)
		case 1: return w.Write2Bytes(1, byte(v))
		case 2: return w.Write3Bytes(2, byte(v), byte(v >> 8))
		case 3: return w.Write4Bytes(3, byte(v), byte(v >> 8), byte(v >> 16))
		case 4: return w.Write5Bytes(4, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
		case 5: return w.Write6Bytes(5, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32))
		case 6: return w.Write7Bytes(6, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
		case 7: return w.Write8Bytes(7, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48))
		case 8: return w.Write9Bytes(8, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 25), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
	}
	return nil
}

// Encode 2 uint64s in 3-13 bytes and write it to the buffer. The length is always 1 byte more than the minimum representation of both the uint64s.
func (w *Buffer) Write2Uint64sVariable(v1 uint64, v2 uint64) error {
	s2 := numbytes(v2)
	switch numbytes(v1) {
		case 0: w.WriteByte(s2)
		case 1: w.Write2Bytes(16 | s2, byte(v1))
		case 2: w.Write3Bytes(32 | s2, byte(v1), byte(v1 >> 8))
		case 3: w.Write4Bytes(48 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16))
		case 4: w.Write5Bytes(64 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24))
		case 5: w.Write6Bytes(80 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32))
		case 6: w.Write7Bytes(96 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40))
		case 7: w.Write8Bytes(112 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48))
		case 8: w.Write9Bytes(128 | s2, byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 25), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48), byte(v1 >> 56))
	}
	switch s2 {
		case 0: return nil
		case 1: return w.WriteByte(byte(v2))
		case 2: return w.Write2Bytes(byte(v2), byte(v2 >> 8))
		case 3: return w.Write3Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16))
		case 4: return w.Write4Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24))
		case 5: return w.Write5Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32))
		case 6: return w.Write6Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40))
		case 7: return w.Write7Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48))
		case 8: return w.Write8Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 25), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48), byte(v2 >> 56))
	}
	return nil
}

// Encode a float32 in 4 bytes and write it to the buffer
func (w *Buffer) WriteFloat32(flt float32) error {
	return w.WriteUint32(math.Float32bits(flt))
}

// Encode a float64 in 8 bytes and write it to the buffer
func (w *Buffer) WriteFloat64(flt float64) error {
	return w.WriteUint64(math.Float64bits(flt))
}

// Write a string to the buffer with maximum length 255
func (w *Buffer) WriteString8(s string) (int, error) {
	if len(s) >= 255 {
		w.WriteByte(255)
		_, err := w.WriteString(s[0:255])
		return 256, err
	} else {
		w.WriteByte(uint8(len(s)))
		_, err := w.WriteString(s)
		return len(s) + 1, err
	}
}

// Write a string to the buffer with maximum length 65,535
func (w *Buffer) WriteString16(s string) (int, error) {
	if len(s) >= 65535 {
		w.WriteUint16(65535)
		_, err := w.WriteString(s[0:65535])
		return 65537, err
	} else {
		w.WriteUint16(uint16(len(s)))
		_, err := w.WriteString(s)
		return len(s) + 2, err
	}
}

// Write a string to the buffer with maximum length 4,294,967,295
func (w *Buffer) WriteString32(s string) (int, error) {
	if len(s) >= 4294967295 {
		w.WriteUint32(4294967295)
		_, err := w.WriteString(s[0:4294967295])
		return 4294967299, err
	} else {
		w.WriteUint32(uint32(len(s)))
		_, err := w.WriteString(s)
		return len(s) + 4, err
	}
}

// Write a slice of bytes to the buffer with maximum length 255
func (w *Buffer) WriteBytes8(s []byte) (int, error) {
	if len(s) >= 255 {
		w.WriteByte(255)
		_, err := w.Write(s[0:255])
		return 256, err
	} else {
		w.WriteByte(uint8(len(s)))
		_, err := w.Write(s)
		return len(s) + 1, err
	}
}

// Write a slice of bytes to the buffer with maximum length 65,535
func (w *Buffer) WriteBytes16(s []byte) (int, error) {
	if len(s) >= 65535 {
		w.WriteUint16(65535)
		_, err := w.Write(s[0:65535])
		return 65537, err
	} else {
		w.WriteUint16(uint16(len(s)))
		_, err := w.Write(s)
		return len(s) + 2, err
	}
}

// Write a slice of bytes to the buffer with maximum length 4,294,967,295
func (w *Buffer) WriteBytes32(s []byte) (int, error) {
	if len(s) >= 4294967295 {
		w.WriteUint32(4294967295)
		_, err := w.Write(s[0:4294967295])
		return 4294967299, err
	} else {
		w.WriteUint32(uint32(len(s)))
		_, err := w.Write(s)
		return len(s) + 4, err
	}
}

// Reflects on the values and writes them all out. Not particularly safe.
// This function only works with: integer, slice of bytes, string, byte.
// uint8 is written as byte. int32 is written as rune. Other integers are written as ASCII representation of their number.
// A slice of anything other than bytes could cause unknown behavior.
func (w *Buffer) WriteAll(a ...interface{}) (n int, err error) {
	var i int
	for _, p := range a {
		switch reflect.TypeOf(p).Kind() {
			case reflect.String:
				i, err = w.WriteString(reflect.ValueOf(p).String())
			case reflect.Slice: // all slices are assumed to be slices of bytes
				i, err = w.Write(reflect.ValueOf(p).Bytes())
			case reflect.Uint8: // byte
				i, err = 1, w.WriteByte(byte(reflect.ValueOf(p).Uint()))
			case reflect.Int32: // rune
				i, err = w.WriteRune(rune(reflect.ValueOf(p).Int()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
				i, err = conv.Write(w, int(reflect.ValueOf(p).Int()), 0)
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err = conv.Write(w, int(reflect.ValueOf(p).Uint()), 0)
			default:
				err = errors.New("custom.Buffer.WriteAll: not a supported type")
				return
		}
		if err != nil {
			return
		}
		n += i
	}
	return
}

// Reset (empty) the buffer
func (w *Buffer) Reset() {
	w.cursor = 0
	return
}

// Return the length of what has been written so far
func (w *Buffer) Len() int {
	return w.cursor
}

// Return the bytes written to the buffer. This slice is not safe to use once the custom.Buffer has Close()
func (w *Buffer) Bytes() []byte {
	return w.data[0:w.cursor]
}

// Return a copy of the bytes written to the buffer. This slice is safe to use even once the custom.Buffer has Close()
func (w *Buffer) BytesCopy() []byte {
	b := make([]byte, w.cursor)
	copy(b, w.data)
	return b
}

// Return a copy of the bytes written to the buffer as a string. This slice is safe to use even once the custom.Buffer has Close()
func (w *Buffer) String() string {
	return string(w.data[0:w.cursor])
}

// Releases the buffer back to the pool
func (w *Buffer) Close() error {
	if w.length == bufferLen {
		pool.Put(w.data)
		w.length = 0
		w.data = nil
	}
	return nil
}

func numbytes(v uint64) uint8 {
	switch {
		case v == 0: return 0
		case v < 256: return 1
		case v < 65536: return 2
		case v < 16777216: return 3
		case v < 4294967296: return 4
		case v < 1099511627776: return 5
		case v < 281474976710655: return 6
		case v < 72057594037927936: return 7
		default: return 8
	}
}

// -------- READER --------

type Reader struct {
	f io.Reader
	at int		// the cursor for where I am in buf
	n int		// how much uncompressed but as of yet unparsed data is left in buf
	buf []byte	// the buffer for reading data
	close, eof bool
}

// Creates a new buffered reader wrapping an io.Reader
func NewReader(f io.Reader) *Reader {
	return &Reader{f: f, buf: pool.Get().([]byte)}
}

// Creates a new buffered reader wrapping an io.Reader which contains Zlib compressed data
func NewZlibReader(f io.Reader) *Reader {
	z, err := zlib.NewReader(f)
	if err != nil {
		panic(err)
	}
	return &Reader{f: z, buf: pool.Get().([]byte), close: true}
}

// Creates a new buffered reader wrapping an io.Reader which contains Snappy compressed data
func NewSnappyReader(f io.Reader) *Reader {
	return &Reader{f: snappy.NewReader(f), buf: pool.Get().([]byte), close: true}
}

func (r *Reader) fill(x int) error {
	copy(r.buf, r.buf[r.at:r.at+r.n])
	r.at = 0
	m, err := r.f.Read(r.buf[r.n:])
	r.n += m
	if err != nil {
		return err
	}
	for r.n < x {
		m, err = r.f.Read(r.buf[r.n:])
		r.n += m
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) fill1() {
	r.at = 0
	m, err := r.f.Read(r.buf)
	r.n = m
	if err != nil {
		panic(err)
	}
}

// Populate slice of bytes
func (r *Reader) Read(b []byte) (int, error) {
	x := len(b)
	if x > bufferLen { // the user has requested more data than the buffer size
		n := r.n
		copy(b, r.buf[r.at:r.at+n]) // copy what we have in the buffer
		r.at, r.n = 0, 0 // buffer is now empty
		if i, err := io.ReadAtLeast(r.f, b[n:], x-n); err != nil { // then read the remainder directly from the src
			return n+i, err
		}
		return x, nil
	}
	var err error
	if r.n < x {
		if err = r.fill(x); err == io.EOF {
			x = r.n
		}
	}
	copy(b, r.buf[r.at:r.at+x]) // must be copied to avoid memory leak
	r.at += x
	r.n -= x
	return x, err
}

// Reads x bytes and returns this slice of bytes as a copy.
func (r *Reader) Readx(x int) []byte {
	b := make([]byte, x)
	if x > bufferLen { // the user has requested more data than the buffer size
		n := r.n
		copy(b, r.buf[r.at:r.at+n]) // copy what we have in the buffer
		r.at, r.n = 0, 0 // buffer is now empty
		if _, err := io.ReadAtLeast(r.f, b[n:], x-n); err != nil { // then read the remainder directly from the src
			panic(err)
		}
		return b
	}
	if r.n < x {
		if err := r.fill(x); err != nil {
			panic(err)
		}
	}
	copy(b, r.buf[r.at:r.at+x]) // must be copied to avoid memory leak
	r.at += x
	r.n -= x
	return b
}

// Reads x bytes and returns a slice of the buffer. This slice is not a copy and so must be used or copied before the next read.
func (r *Reader) ReadxRaw(x int) []byte {
	if x > bufferLen { // the user has requested more data than the buffer size
		b := make([]byte, x)
		n := r.n
		copy(b, r.buf[r.at:r.at+n]) // copy what we have in the buffer
		r.at, r.n = 0, 0 // buffer is now empty
		if _, err := io.ReadAtLeast(r.f, b[n:], x-n); err != nil { // then read the remainder directly from the src
			panic(err)
		}
		return b
	}
	if r.n < x {
		if err := r.fill(x); err != nil {
			panic(err)
		}
	}
	r.at += x
	r.n -= x
	return r.buf[r.at-x:r.at]
}

// Read 1 byte
func (r *Reader) ReadByte() uint8 {
	if r.n == 0 {
		r.fill1()
	}
	r.at++
	r.n--
	return r.buf[r.at-1]
}

// Read and decode a boolean encoded with WriteBool
func (r *Reader) ReadBool() (b1 bool) {
	if r.n == 0 {
		r.fill1()
	}
	if r.buf[r.at] > 0 {
		b1 = true
	}
	r.at++
	r.n--
	return
}

// Read and decode 2 booleans encoded with Write2Bools
func (r *Reader) Read2Bools() (b1 bool, b2 bool) {
	if r.n == 0 {
		r.fill1()
	}
	switch r.buf[r.at] {
		case 1: b1 = true
		case 2: b2 = true
		case 3: b1, b2 = true, true
	}
	r.at++
	r.n--
	return
}

// Read and decode 8 booleans encoded with Write8Bools
func (r *Reader) Read8Bools() (b1 bool, b2 bool, b3 bool, b4 bool, b5 bool, b6 bool, b7 bool, b8 bool) {
	if r.n == 0 {
		r.fill1()
	}
	c := r.buf[r.at]
	if c & 1 > 0 {
		b1 = true
	}
	if c & 2 > 0 {
		b2 = true
	}
	if c & 4 > 0 {
		b3 = true
	}
	if c & 8 > 0 {
		b4 = true
	}
	if c & 16 > 0 {
		b5 = true
	}
	if c & 32 > 0 {
		b6 = true
	}
	if c & 64 > 0 {
		b7 = true
	}
	if c & 128 > 0 {
		b8 = true
	}
	r.at++
	r.n--
	return
}

// Read and decode 2 uint8s encoded with Read2Uint4s
func (r *Reader) Read2Uint4s() (uint8, uint8) {
	if r.n == 0 {
		r.fill1()
	}
	res1, res2 := r.buf[r.at] & 15, r.buf[r.at] >> 4
	r.at++
	r.n--
	return res1, res2
}

// Read a UTF8 and return it in a new slice of bytes
func (r *Reader) ReadUTF8() []byte {
	first := r.ReadByte()
	if first < 128 { // length 1
		return []byte{first}
	}
	if first & 32 == 0 { // length 2
		return []byte{first, r.ReadByte()}
	} else {
		b := make([]byte, 3)
		b[0] = first
		b[1] = r.ReadByte()
		b[2] = r.ReadByte()
		return b
	}
}

// Read a UTF8 and return it as a slice of the buffer. This slice is not a copy and so must be used or copied before the next read.
func (r *Reader) ReadUTF8Raw() []byte {
	if r.n < 3 {
		if err := r.fill(3); err != nil && err != io.EOF {
			panic(err)
		}
	}
	first := r.buf[r.at]
	if first < 128 { // length 1
		r.at++
		r.n--
		return r.buf[r.at-1:r.at]
	}
	if first & 32 == 0 { // length 2
		r.at += 2
		r.n -= 2
		return r.buf[r.at-2:r.at]
	} else {
		r.at += 3
		r.n -= 3
		return r.buf[r.at-3:r.at]
	}
}

// Read a rune which was encoded as UTF8 (or with WriteRune)
func (r *Reader) ReadRune() rune {
	if r.n < 3 {
		if err := r.fill(3); err != nil && err != io.EOF {
			panic(err)
		}
	}
	first := r.buf[r.at]
	if first < 128 { // length 1
		r.at++
		r.n--
		return rune(first)
	}
	if first & 32 == 0 { // length 2
		rn, _ := utf8.DecodeRune(r.buf[r.at:r.at+2])
		r.at += 2
		r.n -= 2
		return rn
	} else {
		rn, _ := utf8.DecodeRune(r.buf[r.at:r.at+3])
		r.at += 3
		r.n -= 3
		return rn
	}
}

// Read and decode a uint16 encoded with WriteUint16
func (r *Reader) ReadUint16() uint16 {
	if r.n < 2 {
		if err := r.fill(2); err != nil {
			panic(err)
		}
	}
	r.at += 2
	r.n -= 2
	return uint16(r.buf[r.at-2]) | uint16(r.buf[r.at-1])<<8
}

// Read and decode a uint16 encoded with WriteUint16Variable
func (r *Reader) ReadUint16Variable() uint16 {
	v := r.ReadByte()
	if v < 255 {
		return uint16(v)
	}
	return r.ReadUint16()
}

// Read and decode an int16 encoded with WriteInt16Variable
func (r *Reader) ReadInt16Variable() int16 {
	v := r.ReadByte()
	if v < 255 {
		return int16(v) - 127
	}
	return int16(r.ReadUint16())
}

// Read and decode an uint32 encoded with WriteUint24
func (r *Reader) ReadUint24() uint32 {
	if r.n < 3 {
		if err := r.fill(3); err != nil {
			panic(err)
		}
	}
	r.at += 3
	r.n -= 3
	return uint32(r.buf[r.at-3]) | uint32(r.buf[r.at-2])<<8 | uint32(r.buf[r.at-1])<<16
}

// Read and decode an uint32 encoded with WriteUint32
func (r *Reader) ReadUint32() uint32 {
	if r.n < 4 {
		if err := r.fill(4); err != nil {
			panic(err)
		}
	}
	r.at += 4
	r.n -= 4
	return uint32(r.buf[r.at-4]) | uint32(r.buf[r.at-3])<<8 | uint32(r.buf[r.at-2])<<16 | uint32(r.buf[r.at-1])<<24
}

// Read and decode an uint64 encoded with WriteUint48
func (r *Reader) ReadUint48() uint64 {
	if r.n < 6 {
		if err := r.fill(6); err != nil {
			panic(err)
		}
	}
	r.at += 6
	r.n -= 6
	return uint64(r.buf[r.at-6]) | uint64(r.buf[r.at-5])<<8 | uint64(r.buf[r.at-4])<<16 | uint64(r.buf[r.at-3])<<24 | uint64(r.buf[r.at-2])<<32 | uint64(r.buf[r.at-1])<<40
}

// Read and decode an uint64 encoded with WriteUint64
func (r *Reader) ReadUint64() uint64 {
	if r.n < 8 {
		if err := r.fill(8); err != nil {
			panic(err)
		}
	}
	r.at += 8
	r.n -= 8
	return uint64(r.buf[r.at-8]) | uint64(r.buf[r.at-7])<<8 | uint64(r.buf[r.at-6])<<16 | uint64(r.buf[r.at-5])<<24 | uint64(r.buf[r.at-4])<<32 | uint64(r.buf[r.at-3])<<40 | uint64(r.buf[r.at-2])<<48 | uint64(r.buf[r.at-1])<<56
}

// Read and decode an uint64 encoded with WriteUint64Variable
func (r *Reader) ReadUint64Variable() uint64 {
	s1 := int(r.ReadByte())
	if r.n < s1 {
		if err := r.fill(s1); err != nil {
			panic(err)
		}
	}
	var res1 uint64
	switch s1 {
		case 1: res1 = uint64(r.buf[r.at])
		case 2: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8
		case 3: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16
		case 4: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24
		case 5: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32
		case 6: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40
		case 7: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48
		case 8: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48 | uint64(r.buf[r.at+7])<<56
	}
	r.at += s1
	r.n -= s1
	return res1
}

// Read and decode 2 uint64s encoded with Write2Uint64sVariable
func (r *Reader) Read2Uint64sVariable() (uint64, uint64) {
	s2 := r.ReadByte()
	s1 := s2 >> 4
	s2 &= 15
	x := int(s1 + s2)
	if r.n < x {
		if err := r.fill(x); err != nil {
			panic(err)
		}
	}
	var res1, res2 uint64
	switch s1 {
		case 1: res1 = uint64(r.buf[r.at])
		case 2: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8
		case 3: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16
		case 4: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24
		case 5: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32
		case 6: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40
		case 7: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48
		case 8: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48 | uint64(r.buf[r.at+7])<<56
	}
	r.at += int(s1)
	r.n -= int(s1)
	switch s2 {
		case 1: res2 = uint64(r.buf[r.at])
		case 2: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8
		case 3: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16
		case 4: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24
		case 5: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32
		case 6: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40
		case 7: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48
		case 8: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48 | uint64(r.buf[r.at+7])<<56
	}
	r.at += int(s2)
	r.n -= int(s2)
	return res1, res2
}

// Read and decode a float32 encoded with WriteFloat32
func (r *Reader) ReadFloat32() float32 {
	return math.Float32frombits(r.ReadUint32())
}

// Read and decode a float64 encoded with WriteFloat64
func (r *Reader) ReadFloat64() float64 {
	return math.Float64frombits(r.ReadUint64())
}

// Read and decode a string encoded with WriteString8
func (r *Reader) ReadString8() string {
	return string(r.ReadxRaw(int(r.ReadByte())))
}

// Read and decode a string encoded with WriteString16
func (r *Reader) ReadString16() string {
	return string(r.ReadxRaw(int(r.ReadUint16())))
}

// Read and decode a string encoded with WriteString32
func (r *Reader) ReadString32() string {
	return string(r.ReadxRaw(int(r.ReadUint32())))
}

// Read and decode a string encoded with WriteString8 as slice of bytes
func (r *Reader) ReadBytes8() []byte {
	return r.Readx(int(r.ReadByte()))
}

// Read and decode a string encoded with WriteString16 as slice of bytes
func (r *Reader) ReadBytes16() []byte {
	return r.Readx(int(r.ReadUint16()))
}

// Read and decode a string encoded with WriteString32 as slice of bytes
func (r *Reader) ReadBytes32() []byte {
	return r.Readx(int(r.ReadUint32()))
}

func (r *Reader) Discard(x int) {
	if r.n < x {
		if err := r.fill(x); err != nil {
			panic(err)
		}
	}
	r.at += x
	r.n -= x
}

// Seeks on the underlying io.Reader
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if sw, ok := r.f.(io.Seeker); ok {
		r.at, r.n = 0, 0
		return sw.Seek(offset, whence)
	}
	return 0, errors.New(`Does not implement io.Seeker`)
}

// Checks whether the end of the underlying io.Reader has been reached. Returns nil if this is already the end. This is safe to do at any time whilst reading to check if the end is reached.
func (r *Reader) EOF() error {
	if r.n > 0 {
		return ErrNotEOF
	}
	m, err := r.f.Read(r.buf)
	r.n = m
	if err == io.EOF {
		return nil
	}
	if err == nil {
		return ErrNotEOF
	}
	return err
}

// Releases the buffer back to the pool
func (r *Reader) Close() error {
	pool.Put(r.buf)
	r.buf = nil
	if r.close {
		if sw, ok := r.f.(io.Closer); ok { // Attempt to close underlying reader if it has a Close() method
			return sw.Close()
		}
	}
	r.f = nil
	return nil
}

// -------- BYTES READER --------

type BytesReader struct {
	data []byte
	cursor, length int
}

// Creates a reader wrapping a slice of bytes
func NewBytesReader(p []byte) *BytesReader {
	return &BytesReader{data: p, length: len(p)}
}

// Populate slice of bytes (copy)
func (r *BytesReader) Read(p []byte) (n int, err error) {
	if to := r.cursor + len(p); to < r.length {
		n = copy(p, r.data[r.cursor:to])
		r.cursor += n
	} else {
		if r.cursor < r.length {
			n = copy(p, r.data[r.cursor:])
			r.cursor += n
		}
		err = io.EOF
	}
	return
}

// Read x bytes and returns this slice of bytes as a copy
func (r *BytesReader) Readx(x int) []byte {
	p := make([]byte, x)
	r.cursor += copy(p, r.data[r.cursor:r.cursor+x])
	return p
}

// Returns a slice of the original. This slice is not a copy and so should not be modified
func (r *BytesReader) ReadxRaw(x int) []byte {
	r.cursor += x
	return r.data[r.cursor-x:r.cursor]
}

// Read 1 byte
func (r *BytesReader) ReadByte() uint8 {
	r.cursor++
	return r.data[r.cursor-1]
}

// Read and decode a boolean encoded with WriteBool
func (r *BytesReader) ReadBool() (b1 bool) {
	if r.data[r.cursor] > 0 {
		b1 = true
	}
	r.cursor++
	return
}

// Read and decode 2 booleans encoded with Write2Bools
func (r *BytesReader) Read2Bools() (b1 bool, b2 bool) {
	switch r.data[r.cursor] {
		case 1: b1 = true
		case 2: b2 = true
		case 3: b1, b2 = true, true
	}
	r.cursor++
	return
}

// Read and decode 8 booleans encoded with Write8Bools
func (r *BytesReader) Read8Bools() (b1 bool, b2 bool, b3 bool, b4 bool, b5 bool, b6 bool, b7 bool, b8 bool) {
	c := r.data[r.cursor]
	if c & 1 > 0 {
		b1 = true
	}
	if c & 2 > 0 {
		b2 = true
	}
	if c & 4 > 0 {
		b3 = true
	}
	if c & 8 > 0 {
		b4 = true
	}
	if c & 16 > 0 {
		b5 = true
	}
	if c & 32 > 0 {
		b6 = true
	}
	if c & 64 > 0 {
		b7 = true
	}
	if c & 128 > 0 {
		b8 = true
	}
	r.cursor++
	return
}

// Read and decode 2 uint8s encoded with Read2Uint4s
func (r *BytesReader) Read2Uint4s() (uint8, uint8) {
	res1, res2 := r.data[r.cursor] & 15, r.data[r.cursor] >> 4
	r.cursor++
	return res1, res2
}

// Read a UTF8 and return it in a new slice of bytes
func (r *BytesReader) ReadUTF8() []byte {
	if r.data[r.cursor] < 128 { // length 1
		r.cursor++
		return []byte{r.data[r.cursor-1]}
	}
	if r.data[r.cursor] & 32 == 0 { // length 2
		r.cursor += 2
		return []byte{r.data[r.cursor-2], r.data[r.cursor-1]}
	} else {
		r.cursor += 3
		return []byte{r.data[r.cursor-3], r.data[r.cursor-2], r.data[r.cursor-1]}
	}
}

// Read a UTF8 and return a reslice of the original slice. This slice is not a copy and so should not be modified
func (r *BytesReader) ReadUTF8Raw() []byte {
	if r.data[r.cursor] < 128 { // length 1
		r.cursor++
		return r.data[r.cursor-1:r.cursor]
	}
	if r.data[r.cursor] & 32 == 0 { // length 2
		r.cursor += 2
		return r.data[r.cursor-2:r.cursor]
	} else {
		r.cursor += 3
		return r.data[r.cursor-3:r.cursor]
	}
}

// Read a rune which was encoded as UTF8 (or with WriteRune)
func (r *BytesReader) ReadRune() rune {
	if r.data[r.cursor] < 128 { // length 1
		r.cursor++
		return rune(r.data[r.cursor-1])
	}
	if r.data[r.cursor] & 32 == 0 { // length 2
		r.cursor += 2
		rn, _ := utf8.DecodeRune(r.data[r.cursor-2:r.cursor])
		return rn
	} else {
		r.cursor += 3
		rn, _ := utf8.DecodeRune(r.data[r.cursor-3:r.cursor])
		return rn
	}
}

// Read and decode a uint16 encoded with WriteUint16
func (r *BytesReader) ReadUint16() uint16 {
	r.cursor += 2
	return uint16(r.data[r.cursor-2]) | uint16(r.data[r.cursor-1])<<8
}

// Read and decode a uint16 encoded with WriteUint16Variable
func (r *BytesReader) ReadUint16Variable() uint16 {
	v := r.ReadByte()
	if v < 255 {
		return uint16(v)
	}
	return r.ReadUint16()
}

// Read and decode an int16 encoded with WriteInt16Variable
func (r *BytesReader) ReadInt16Variable() int16 {
	v := r.ReadByte()
	if v < 255 {
		return int16(v) - 127
	}
	return int16(r.ReadUint16())
}

// Read and decode an uint32 encoded with WriteUint24
func (r *BytesReader) ReadUint24() uint32 {
	r.cursor += 3
	return uint32(r.data[r.cursor-3]) | uint32(r.data[r.cursor-2])<<8 | uint32(r.data[r.cursor-1])<<16
}

// Read and decode an uint32 encoded with WriteUint32
func (r *BytesReader) ReadUint32() uint32 {
	r.cursor += 4
	return uint32(r.data[r.cursor-4]) | uint32(r.data[r.cursor-3])<<8 | uint32(r.data[r.cursor-2])<<16 | uint32(r.data[r.cursor-1])<<24
}

// Read and decode an uint64 encoded with WriteUint48
func (r *BytesReader) ReadUint48() uint64 {
	r.cursor += 6
	return uint64(r.data[r.cursor-6]) | uint64(r.data[r.cursor-5])<<8 | uint64(r.data[r.cursor-4])<<16 | uint64(r.data[r.cursor-3])<<24 | uint64(r.data[r.cursor-2])<<32 | uint64(r.data[r.cursor-1])<<40
}

// Read and decode an uint64 encoded with WriteUint64
func (r *BytesReader) ReadUint64() uint64 {
	r.cursor += 8
	return uint64(r.data[r.cursor-8]) | uint64(r.data[r.cursor-7])<<8 | uint64(r.data[r.cursor-6])<<16 | uint64(r.data[r.cursor-5])<<24 | uint64(r.data[r.cursor-4])<<32 | uint64(r.data[r.cursor-3])<<40 | uint64(r.data[r.cursor-2])<<48 | uint64(r.data[r.cursor-1])<<56
}

// Read and decode an uint64 encoded with WriteUint64Variable
func (r *BytesReader) ReadUint64Variable() uint64 {
	s1 := int(r.ReadByte())
	var res1 uint64
	switch s1 {
		case 1: res1 = uint64(r.data[r.cursor])
		case 2: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8
		case 3: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16
		case 4: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24
		case 5: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32
		case 6: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40
		case 7: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48
		case 8: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48 | uint64(r.data[r.cursor+7])<<56
	}
	r.cursor += s1
	return res1
}

// Read and decode 2 uint64s encoded with Write2Uint64sVariable
func (r *BytesReader) Read2Uint64sVariable() (uint64, uint64) {
	s2 := r.ReadByte()
	s1 := s2 >> 4
	s2 &= 15
	var res1, res2 uint64
	switch s1 {
		case 1: res1 = uint64(r.data[r.cursor])
		case 2: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8
		case 3: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16
		case 4: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24
		case 5: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32
		case 6: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40
		case 7: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48
		case 8: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48 | uint64(r.data[r.cursor+7])<<56
	}
	r.cursor += int(s1)
	switch s2 {
		case 1: res2 = uint64(r.data[r.cursor])
		case 2: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8
		case 3: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16
		case 4: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24
		case 5: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32
		case 6: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40
		case 7: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48
		case 8: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48 | uint64(r.data[r.cursor+7])<<56
	}
	r.cursor += int(s2)
	return res1, res2
}

// Read and decode a float32 encoded with WriteFloat32
func (r *BytesReader) ReadFloat32() float32 {
	return math.Float32frombits(r.ReadUint32())
}

// Read and decode a float64 encoded with WriteFloat64
func (r *BytesReader) ReadFloat64() float64 {
	return math.Float64frombits(r.ReadUint64())
}

// Read and decode a string encoded with WriteString8
func (r *BytesReader) ReadString8() string {
	return string(r.ReadxRaw(int(r.ReadByte())))
}

// Read and decode a string encoded with WriteString16
func (r *BytesReader) ReadString16() string {
	return string(r.ReadxRaw(int(r.ReadUint16())))
}

// Read and decode a string encoded with WriteString32
func (r *BytesReader) ReadString32() string {
	return string(r.ReadxRaw(int(r.ReadUint32())))
}

// Read and decode a string encoded with WriteString8 as slice of bytes
func (r *BytesReader) ReadBytes8() []byte {
	return r.Readx(int(r.ReadByte()))
}

// Read and decode a string encoded with WriteString16 as slice of bytes
func (r *BytesReader) ReadBytes16() []byte {
	return r.Readx(int(r.ReadUint16()))
}

// Read and decode a string encoded with WriteString32 as slice of bytes
func (r *BytesReader) ReadBytes32() []byte {
	return r.Readx(int(r.ReadUint32()))
}

// Moves the cursor forward x bytes without returning anything
func (r *BytesReader) Discard(x int) {
	r.cursor += x
}

// Implements io.Seeker (see io.Seeker for usage)
func (r *BytesReader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
		case 0:
			abs = offset
		case 1:
			abs = int64(r.cursor) + offset
		case 2:
			abs = int64(r.length) + offset
		default:
			return 0, errors.New("buffer.BytesReader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("buffer.BytesReader.Seek: negative position")
	}
	r.cursor = int(abs)
	return abs, nil
}

// Returns nil if the cursor is at the end of the slice of bytes, otherwise returns an error
func (r *BytesReader) EOF() error {
	if r.cursor >= len(r.data) {
		return nil
	}
	return ErrNotEOF
}

