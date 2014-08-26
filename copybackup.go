package copybackup

import ( // {{{
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	core "github.com/yukimemi/gocore"
) // }}}

type CopyGroup struct { // {{{
	src        string
	dst        string
	bkpath     string
	generation int
	sleep      int
} // }}}

func makeDstPath(src, bkpath string) (string, error) { // {{{
	var dst string

	f, e := os.Stat(src)
	if e == nil && f.IsDir() {
		// core.Logger.Warnf("%s is Directory !", src)
		return "", fmt.Errorf("%s is Directory !", src)
	}
	parent, base := filepath.Split(src)
	ext := filepath.Ext(base)
	basename := strings.TrimSuffix(base, ext)
	now := time.Now()
	if filepath.IsAbs(bkpath) {
		dst = filepath.Join(bkpath, basename+"_"+now.Format("20060102-150405")+ext)
	} else {
		dst = filepath.Join(parent, bkpath, basename+"_"+now.Format("20060102-150405")+ext)
	}

	return dst, e
} // }}}

func NewCopyGroup(src, bkpath string, generation, sleep int) *CopyGroup { // {{{
	dst, e := makeDstPath(src, bkpath)
	core.FailOnError(e)
	return &CopyGroup{src, dst, bkpath, generation, sleep}
} // }}}

func (cg *CopyGroup) Backup() error { // {{{
	var e error
	dstParent := filepath.Dir(cg.dst)
	e = os.MkdirAll(dstParent, os.ModePerm)
	core.FailOnError(e)
	core.Logger.Infof("%s -> %s", cg.src, cg.dst)
	return cp(cg.dst, cg.src)
} // }}}

func cp(dst, src string) error { // {{{
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
} // }}}
