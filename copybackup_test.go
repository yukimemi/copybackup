package copybackup

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/umisama/golog"
	core "github.com/yukimemi/gocore"
)

const THIS_FILE_NAME string = "copybackup_test.go"
const GENERATION = 3
const TEST_FILES_COUNT = 6

var cg *CopyGroup
var files []string
var oldestFileExpected string
var latestFileExpected string

func init() { // {{{
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)
	pwd, _ := os.Getwd()
	files = make([]string, 0)
	cg = NewCopyGroup(filepath.Join(pwd, THIS_FILE_NAME), "_old", GENERATION, 60)
	dstDir := filepath.Dir(cg.dst)
	if _, e := os.Stat(dstDir); e == nil {
		os.RemoveAll(dstDir)
	}
	os.MkdirAll(dstDir, os.ModePerm)

	for i := 0; i < TEST_FILES_COUNT; i++ {
		cg = NewCopyGroup(filepath.Join(pwd, THIS_FILE_NAME), "_old", GENERATION, 60)
		cp(cg.dst, cg.src)
		if i == 0 {
			oldestFileExpected = cg.dst
		} else if i == TEST_FILES_COUNT-1 {
			latestFileExpected = cg.dst
		}
		files = append(files, cg.dst)
		time.Sleep(time.Second)
	}
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
	cg := NewCopyGroup(filepath.Join(pwd, THIS_FILE_NAME), "_old", -1, 60)
	cg.deleteOldFile()
	count := cg.countMatchFiles()
	expected := TEST_FILES_COUNT

	if count != expected {
		t.Errorf("expected = [%d] but count = [%d]", expected, count)
	}
} // }}}

func TestDeleteOldFile2(t *testing.T) { // {{{
	cg.deleteOldFile()
	count := cg.countMatchFiles()
	expected := GENERATION

	if count != expected {
		t.Errorf("expected = [%d] but count = [%d]", expected, count)
	}
} // }}}

func TestBackup(t *testing.T) {
	cg.Backup()
	count := cg.countMatchFiles()
	expected := GENERATION

	if count != expected {
		t.Errorf("expected = [%d] but count = [%d]", expected, count)
	}
}
