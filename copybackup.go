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
	sleep      int
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
	latestFile, e := cg.getLatestFile()
	core.FailOnError(e)
	if latestFile != "" {
		l, _ := os.Stat(latestFile)
		s, _ := os.Stat(cg.src)
		core.Logger.Debug("cg.src ModTime = ", s.ModTime().Format(time.StampMilli))
		core.Logger.Debug("latestFile ModTime = ", l.ModTime().Format(time.StampMilli))
		if l.ModTime().Equal(s.ModTime()) {
			core.Logger.Debugf("[%s] is same as [%s]", cg.src, latestFile)
			return nil
		}
	}
	core.Logger.Infof("%s -> %s", cg.src, cg.dst)
	e = cp(cg.dst, cg.src)
	cg.deleteOldFile()
	return e
} // }}}

func (cg *CopyGroup) deleteOldFile() { // {{{
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

	latestFile = files[0]
	for _, file := range files {
		m, e := regexp.MatchString(matchs, file.Name())
		core.FailOnError(e)
		if m {
			if latestFile.ModTime().Before(file.ModTime()) {
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

	oldestFile = files[0]
	for _, file := range files {
		m, e := regexp.MatchString(matchs, file.Name())
		core.FailOnError(e)
		if m {
			if oldestFile.ModTime().After(file.ModTime()) {
				oldestFile = file
			}
		}
	}
	return filepath.Join(dstDir, oldestFile.Name()), e
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
