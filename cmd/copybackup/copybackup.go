package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/umisama/golog"
	cb "github.com/yukimemi/copybackup"
	core "github.com/yukimemi/gocore"
)

type Options struct {
}

func main() {
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)

	var root string
	var e error
	var wg sync.WaitGroup

	g := flag.Int("g", -1, "バックアップする世代。")
	b := flag.String("b", "_old", "バックアップを保存する先。絶対パスでの指定も可能。 ([デフォルト _old])")
	s := flag.Int("s", 60*5, "バックアップ間隔。 (秒 [デフォルト 5分])")

	core.Logger.Debugf("g = %d", *g)
	core.Logger.Debugf("b = %s", *b)
	core.Logger.Debugf("s = %d", *s)
	flag.Parse()

	if flag.NArg() != 0 {
		root = flag.Arg(0)
	} else {
		root, e = os.Getwd()
		core.FailOnError(e)
		flag.Usage()
		os.Exit(0)
	}

	files, e := ioutil.ReadDir(root)
	core.FailOnError(e)

	for _, f := range files {
		if !f.IsDir() {
			src := filepath.Join(root, f.Name())
			dst, e := cb.MakeDstPath(src, "_old")
			core.FailOnError(e)
			wg.Add(1)
			go func(src, dst string) {
				defer wg.Done()
				cb.Backup(src, dst)
			}(src, dst)
		}
		wg.Wait()
	}
}
