package copybackup

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/umisama/golog"
	core "github.com/yukimemi/gocore"
)

const THIS_FILE_NAME string = "copybackup_test.go"
const DUMMY_FILE_NAME string = "DUMMY.dmy"
const GENERATION = 3
const TEST_FILES_COUNT = 6

var cg *CopyGroup
var files []string
var oldestFileExpected string
var latestFileExpected string
var dummy1 string
var dummy2 string

func init() { // {{{
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)
	pwd, _ := os.Getwd()
	src := filepath.Join(pwd, THIS_FILE_NAME)
	files = make([]string, 0)
	cg = NewCopyGroup(src, "_old", GENERATION)
	dstDir := filepath.Dir(cg.dst)
	if _, e := os.Stat(dstDir); e == nil {
		os.RemoveAll(dstDir)
	}
	os.MkdirAll(dstDir, os.ModePerm)

	dummy1 = filepath.Join(pwd, "_old", DUMMY_FILE_NAME+"1")
	core.Logger.Debugf("dummy1 = [%s]", dummy1)
	ioutil.WriteFile(dummy1, []byte("DUMMY"), os.ModePerm)

	for i := 0; i < TEST_FILES_COUNT; i++ {
		cg = NewCopyGroup(src, "_old", GENERATION)
		cp(cg.dst, cg.src)
		os.Chtimes(cg.dst, time.Now(), time.Now())
		if i == 0 {
			oldestFileExpected = cg.dst
		} else if i == TEST_FILES_COUNT-1 {
			latestFileExpected = cg.dst
		}
		files = append(files, cg.dst)
		time.Sleep(time.Second)
	}

	dummy2 = filepath.Join(pwd, "_old", DUMMY_FILE_NAME+"2")
	core.Logger.Debugf("dummy2 = [%s]", dummy2)
	ioutil.WriteFile(dummy2, []byte("DUMMY"), os.ModePerm)

} // }}}

func TestMakeDstPath1(t *testing.T) { // {{{
	pwd, _ := os.Getwd()
	src := filepath.Join(pwd, THIS_FILE_NAME)
	bkpath := "_old"

	expected := "^" + pwd + `/_old/copybackup_test_[0-9]{8}-[0-9]{6}\.go$`
	dst, e := makeDstPath(src, bkpath)
	if e != nil {
		t.Error(e)
	}

	if m, e := regexp.MatchString(expected, dst); m && e == nil {
		// core.Logger.Infof("dst = [%s]", dst)
	} else {
		t.Errorf("expected = [%s] but dst is = [%s]", expected, dst)
	}
} // }}}

func TestMakeDstPath2(t *testing.T) { // {{{
	pwd, _ := os.Getwd()
	src := filepath.Join(pwd, THIS_FILE_NAME)
	bkpath := "/Users/yukimemi/backup"

	expected := `^/Users/yukimemi/backup/copybackup_test_[0-9]{8}-[0-9]{6}\.go$`
	dst, e := makeDstPath(src, bkpath)
	if e != nil {
		t.Error(e)
	}

	if m, e := regexp.MatchString(expected, dst); m && e == nil {
		// core.Logger.Infof("dst = [%s]", dst)
	} else {
		t.Errorf("expected = [%s] but dst is = [%s]", expected, dst)
	}
} // }}}

func TestCountMatchFiles(t *testing.T) { // {{{
	count := cg.countMatchFiles()
	expected := TEST_FILES_COUNT

	if count != expected {
		t.Errorf("expected = [%d] but count = [%d]", expected, count)
	}
} // }}}

func TestGetOldestFile(t *testing.T) { // {{{
	oldestFile, _ := cg.getOldestFile()
	if oldestFile != oldestFileExpected {
		t.Errorf("expected = [%s] but oldestFile = [%s]", oldestFileExpected, oldestFile)
	}
} // }}}

func TestGetLatestFile(t *testing.T) { // {{{
	latestFile, _ := cg.getLatestFile()
	if latestFile != latestFileExpected {
		t.Errorf("expected = [%s] but latestFile = [%s]", latestFileExpected, latestFile)
	}
} // }}}

func TestDeleteOldFile1(t *testing.T) { // {{{
	pwd, _ := os.Getwd()
	cg := NewCopyGroup(filepath.Join(pwd, THIS_FILE_NAME), "_old", -1)
	cg.DeleteOldFile()

	if _, e := os.Stat(oldestFileExpected); e != nil {
		t.Errorf("[%s] expects exists, but it was deleted !", oldestFileExpected)
	} else if _, e := os.Stat(dummy1); e != nil {
		t.Errorf("[%s] expects exists, but it was deleted !", dummy1)
	} else if _, e := os.Stat(dummy2); e != nil {
		t.Errorf("[%s] expects exists, but it was deleted !", dummy2)
	}
} // }}}

func TestDeleteOldFile2(t *testing.T) { // {{{
	cg.DeleteOldFile()

	if _, e := os.Stat(oldestFileExpected); e == nil {
		t.Errorf("[%s] expects deleted, but it is exists !", oldestFileExpected)
	} else if _, e := os.Stat(dummy1); e != nil {
		t.Errorf("[%s] expects exists, but it was deleted !", dummy1)
	} else if _, e := os.Stat(dummy2); e != nil {
		t.Errorf("[%s] expects exists, but it was deleted !", dummy2)
	}
} // }}}

func TestBackup(t *testing.T) { // {{{
	cg := NewCopyGroup(cg.src, cg.bkpath, cg.generation)
	cg.Backup()
	count := cg.countMatchFiles()
	expected := GENERATION + 1

	if count != expected {
		t.Errorf("expected = [%d] but count = [%d]", expected, count)
	}
} // }}}
