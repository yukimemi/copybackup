package main

import ( // {{{
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
) // }}}

// for debug {{{
type debugT bool

var debug = debugT(true)

func (d debugT) Println(args ...interface{}) { // {{{
	if d {
		log.Println(args...)
	}
} // }}}

func (d debugT) Printf(format string, v ...interface{}) { // {{{
	if d {
		log.Printf(format, v...)
	}
} // }}}

func (d debugT) PrintValue(vs string, v interface{}) { // {{{
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

func failOnError(e error) { // {{{
	if e != nil {
		log.Fatal("Error:", e)
	}
} // }}}

type options struct {
	path       string
	generation int
	backup     string
	sleep      int
}

func NewOptions(path string, generation int, backup string, sleep int) *options {
	return &options{path, generation, backup, sleep}
}

func main() {
	defer un(trace("main"))

	var root string
	var e error
	var wg sync.WaitGroup

	g := flag.Int("g", -1, "バックアップする世代。")
	b := flag.String("b", "_old", "バックアップを保存する先。絶対パスでの指定も可能。 ([デフォルト _old])")
	s := flag.Int("s", 60*5, "バックアップ間s。 (秒 [デフォルト 5分])")
	flag.Parse()

	debug.Println("g = ", *g)
	debug.Println("b = ", *b)
	debug.Println("s = ", *s)

	debug.PrintValue("args", flag.Args())
	debug.PrintValue("args", os.Args)
	if flag.NArg() != 0 {
		root = flag.Arg(0)
	} else {
		root, e = os.Getwd()
		failOnError(e)
		flag.Usage()
		os.Exit(1)
	}

	debug.PrintValue("root", root)

	files, e := ioutil.ReadDir(root)
	failOnError(e)

	for _, f := range files {
		if !f.IsDir() {
			src := filepath.Join(root, f.Name())
			dst, e := makeDstPath(src, "_old")
			failOnError(e)
			wg.Add(1)
			go func(src, dst string) {
				defer wg.Done()
				backup(src, dst)
			}(src, dst)
		}
		wg.Wait()
	}
}

func makeDstPath(path, bkpath string) (string, error) {
	_, e := os.Stat(path)
	parent, base := filepath.Split(path)
	ext := filepath.Ext(base)
	basename := strings.TrimSuffix(base, ext)
	now := time.Now()
	dst := filepath.Join(parent, bkpath, basename+"_"+now.Format("20060102_150405")+ext)

	return dst, e
}

func backup(src, dst string) {
	dstParent := filepath.Dir(dst)
	os.MkdirAll(dstParent, 0755)
	e := cp(dst, src)
	debug.PrintValue("src", src)
	debug.PrintValue("dst", dst)
	failOnError(e)
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
