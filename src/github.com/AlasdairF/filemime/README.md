## FileMime

While Go does already include a [mime.TypeByExtension](https://golang.org/pkg/mime/#TypeByExtension) function, this FileMime package is more complete (for use in web apps), works with filename, path or extension, and it's much faster.

## Usage Examples

    f1 := []byte(`html`)
    mimetype := filemime.Get(f1) // text/html
	
    f2 := `/home/mystuff/work.pdf`
    mimetype = filemime.Get([]byte(f2)) // application/pdf
	
	// filemime.Get(f) is the same as filemime.Ext2Mime(filemime.GetExt(f))
    f3 := `http://www.whatever.com/somefile.jpg`
	ext := filemime.GetExt([]byte(f3)) // jpg (in bytes)
	mimetype = filemime.Ext2Mime(ext) // image/jpeg
