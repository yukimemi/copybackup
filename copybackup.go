package copybackup

import ( // {{{
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	core "github.com/yukimemi/gocore"
) // }}}

type CopyGroup struct { // {{{
	src        string
	dst        string
	bkpath     string
	generation int
} // }}}

func makeDstPath(src, bkpath string) (string, error) { // {{{
	var dst string

	f, e := os.Stat(src)
	if e == nil && f.IsDir() {
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

func NewCopyGroup(src, bkpath string, generation int) *CopyGroup { // {{{
	// core.Logger.Debugf("src = [%s]", src)
	dst, e := makeDstPath(src, bkpath)
	// core.Logger.Debugf("dst = [%s]", dst)
	core.FailOnError(e)
	return &CopyGroup{src, dst, bkpath, generation}
} // }}}

func (cg *CopyGroup) Backup() error { // {{{
	var e error
	dstParent := filepath.Dir(cg.dst)
	e = os.MkdirAll(dstParent, os.ModePerm)
	core.FailOnError(e)
	latestFile, e := cg.getLatestFile()
	core.FailOnError(e)
	if latestFile != "" {
		l, _ := os.Stat(latestFile)
		s, _ := os.Stat(cg.src)
		core.Logger.Debugf("src = [%s], time = [%s]", s.Name(), s.ModTime().Format(time.StampMilli))
		core.Logger.Debugf("latestFile = [%s], time = [%s]", l.Name(), l.ModTime().Format(time.StampMilli))
		if l.ModTime().Equal(s.ModTime()) {
			core.Logger.Debugf("[%s] is same as [%s]", cg.src, latestFile)
			return nil
		}
	}
	core.Logger.Infof("%s -> %s", cg.src, cg.dst)
	return cp(cg.dst, cg.src)
} // }}}

func (cg *CopyGroup) DeleteOldFile() { // {{{
	if cg.generation != -1 {
		for {
			if cg.countMatchFiles() > cg.generation {
				oldestFile, e := cg.getOldestFile()
				core.FailOnError(e)
				os.Remove(oldestFile)
			} else {
				break
			}
		}
	}
} // }}}

func (cg *CopyGroup) countMatchFiles() int { // {{{
	var count int

	basename := core.GetBaseName(cg.src)
	ext := filepath.Ext(cg.src)
	matchs := `^` + basename + `_[0-9]{8}-[0-9]{6}` + ext + `$`

	files, e := ioutil.ReadDir(filepath.Dir(cg.dst))
	core.FailOnError(e)

	for _, file := range files {
		m, e := regexp.MatchString(matchs, file.Name())
		core.FailOnError(e)
		if m {
			count++
		}
	}
	return count
} // }}}

func (cg *CopyGroup) getLatestFile() (string, error) { // {{{
	var latestFile os.FileInfo
	var e error

	basename := core.GetBaseName(cg.src)
	ext := filepath.Ext(cg.src)
	matchs := `^` + basename + `_[0-9]{8}-[0-9]{6}` + ext + `$`

	dstDir := filepath.Dir(cg.dst)
	files, e := ioutil.ReadDir(dstDir)
	core.FailOnError(e)

	if len(files) == 0 {
		return "", nil
	}

	for _, file := range files {
		m, e := regexp.MatchString(matchs, file.Name())
		core.FailOnError(e)
		if m {
			if latestFile == nil {
				latestFile = file
			} else if latestFile.ModTime().Before(file.ModTime()) {
				latestFile = file
			}
		}
	}
	return filepath.Join(dstDir, latestFile.Name()), e
} // }}}

func (cg *CopyGroup) getOldestFile() (string, error) { // {{{
	var oldestFile os.FileInfo
	var e error

	basename := core.GetBaseName(cg.src)
	ext := filepath.Ext(cg.src)
	matchs := `^` + basename + `_[0-9]{8}-[0-9]{6}` + ext + `$`

	dstDir := filepath.Dir(cg.dst)
	files, e := ioutil.ReadDir(dstDir)
	core.FailOnError(e)

	if len(files) == 0 {
		return "", nil
	}

	for _, file := range files {
		m, e := regexp.MatchString(matchs, file.Name())
		core.FailOnError(e)
		if m {
			if oldestFile == nil {
				oldestFile = file
			} else if oldestFile.ModTime().After(file.ModTime()) {
				oldestFile = file
			}
		}
	}
	return filepath.Join(dstDir, oldestFile.Name()), e
} // }}}

func cp(dst, src string) error { // {{{
	// https://gist.github.com/elazarl/5507969
	s, e := os.Open(src)
	core.FailOnError(e)
	sinfo, e := os.Stat(src)
	core.FailOnError(e)

	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, e := os.Create(dst)
	defer d.Close()
	if e != nil {
		return e
	}
	if _, e := io.Copy(d, s); e != nil {
		return e
	}
	return os.Chtimes(dst, sinfo.ModTime(), sinfo.ModTime())
} // }}}
