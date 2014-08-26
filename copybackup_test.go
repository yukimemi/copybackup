package copybackup

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/umisama/golog"
	core "github.com/yukimemi/gocore"
)

func TestMakeDstPath1(t *testing.T) { // {{{
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)
	pwd, _ := os.Getwd()
	src := filepath.Join(pwd, "copybackup_test.go")
	bkpath := "_old"

	expected := "^" + pwd + `/_old/copybackup_test_[0-9]{8}-[0-9]{6}\.go$`
	dst, e := makeDstPath(src, bkpath)
	if e != nil {
		t.Error(e)
	}

	if m, e := regexp.MatchString(expected, dst); m && e == nil {
		core.Logger.Infof("dst = [%s]", dst)
	} else {
		t.Errorf("expected = [%s] but dst is = [%s]", expected, dst)
	}
} // }}}

func TestMakeDstPath2(t *testing.T) { // {{{
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)
	pwd, _ := os.Getwd()
	src := filepath.Join(pwd, "copybackup_test.go")
	bkpath := "/Users/yukimemi/backup"

	expected := `^/Users/yukimemi/backup/copybackup_test_[0-9]{8}-[0-9]{6}\.go$`
	dst, e := makeDstPath(src, bkpath)
	if e != nil {
		t.Error(e)
	}

	if m, e := regexp.MatchString(expected, dst); m && e == nil {
		core.Logger.Infof("dst = [%s]", dst)
	} else {
		t.Errorf("expected = [%s] but dst is = [%s]", expected, dst)
	}
} // }}}

func TestCopyGroup(t *testing.T) {
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)
	pwd, _ := os.Getwd()
	src := filepath.Join(pwd, "copybackup_test.go")
	cg := NewCopyGroup(src, "_old", 3, 60*10)
	core.Logger.Infof("cg = [%s]", cg)
}
