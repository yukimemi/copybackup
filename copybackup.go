package copybackup

import ( // {{{
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
) // }}}

// for debug {{{
type DebugT bool

var debug = DebugT(true)

func (d DebugT) Println(args ...interface{}) { // {{{
	if d {
		log.Println(args...)
	}
} // }}}

func (d DebugT) Printf(format string, v ...interface{}) { // {{{
	if d {
		log.Printf(format, v...)
	}
} // }}}

func (d DebugT) PrintValue(vs string, v interface{}) { // {{{
	if d {
		log.Printf("%s = [%s]", vs, v)
	}
} // }}}

func trace(s string) string { // {{{
	debug.Println("ENTER:", s)
	return s
} // }}}

func un(s string) { // {{{
	debug.Println("LEAVE:", s)
} // }}}

// }}}

func FailOnError(e error) { // {{{
	if e != nil {
		log.Fatal("Error:", e)
	}
} // }}}

type Options struct {
	path       string
	generation int
	backup     string
	sleep      int
}

func NewOptions(path string, generation int, backup string, sleep int) *Options {
	return &Options{path, generation, backup, sleep}
}

func MakeDstPath(path, bkpath string) (string, error) {
	_, e := os.Stat(path)
	parent, base := filepath.Split(path)
	ext := filepath.Ext(base)
	basename := strings.TrimSuffix(base, ext)
	now := time.Now()
	dst := filepath.Join(parent, bkpath, basename+"_"+now.Format("20060102_150405")+ext)

	return dst, e
}

func Backup(src, dst string) {
	dstParent := filepath.Dir(dst)
	os.MkdirAll(dstParent, 0755)
	e := cp(dst, src)
	debug.PrintValue("src", src)
	debug.PrintValue("dst", dst)
	FailOnError(e)
}

func cp(dst, src string) error {
	// https://gist.github.com/elazarl/5507969
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
