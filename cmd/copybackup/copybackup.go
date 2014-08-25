package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jessevdk/go-flags"
	cb "github.com/yukimemi/copybackup"
	core "github.com/yukimemi/gocore"
)

type Options struct {
	Generation int    `short:"g" long:"generation" description:"バックアップする世代。"`
	BackupDst  string `short:"b" long:"backup" description:"" `
	Sleep      int    `short:"s" long:"sleep" description:"バックアップ間隔。 (秒 [デフォルト 5分])"`
}

func main() {
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)

	var root string
	var e error
	var wg sync.WaitGroup

	opts := &Options{}
	p := flags.NewParser(opts, flags.PrintErrors)

	core.Logger.Debugf("g = %d", *g)
	core.Logger.Debugf("b = %s", *b)
	core.Logger.Debugf("s = %d", *s)
	flags.Parse()

	if flags.NArg() != 0 {
		root = flags.Arg(0)
	} else {
		root, e = os.Getwd()
		core.FailOnError(e)
		flags.Usage()
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
