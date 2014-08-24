package copybackup

import ( // {{{
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	core "github.com/yukimemi/gocore"
) // }}}

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
	var e error
	dstParent := filepath.Dir(dst)
	e = os.MkdirAll(dstParent, os.ModePerm)
	core.FailOnError(e)
	core.Logger.Infof("%s -> %s", src, dst)
	e = cp(dst, src)
	core.FailOnError(e)
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
