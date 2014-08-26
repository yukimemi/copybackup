package copybackup

import (
	"os"
	"regexp"
	"testing"

	"github.com/umisama/golog"
	core "github.com/yukimemi/gocore"
)

func TestMakeDstPath1(t *testing.T) { // {{{
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)
	src := "/tmp/"
	bkpath := "_old"

	expected := "^/tmp/_old/test_[0-9]{8}-[0-9]{6}$"
	dst, e := makeDstPath(src, bkpath)
	if e != nil {
		t.Error(e)
	}

	if m, e := regexp.MatchString(expected, dst); m && e == nil {
	} else {
		t.Errorf("expected = [%s] but dst is = [%s]", expected, dst)
	}
} // }}}

func TestMakeDstPath2(t *testing.T) { // {{{
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)
	src := "/tmp"
	bkpath := "/Users/yukimemi/backup"

	expected := "^/Users/yukimemi/backup/test_[0-9]{8}-[0-9]{6}$"
	dst, e := makeDstPath(src, bkpath)
	if e != nil {
		t.Error(e)
	}

	if m, e := regexp.MatchString(expected, dst); m && e == nil {
		core.Logger.Infof("dst := [%s]", dst)
	} else {
		t.Errorf("expected = [%s] but dst is = [%s]", expected, dst)
	}
} // }}}
