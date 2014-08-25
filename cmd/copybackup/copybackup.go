package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/umisama/golog"
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
	parser := flags.NewParser(opts, flags.Default)
	args, e := parser.Parse()
	core.FailOnError(e)

	if len(args) == 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(0)
	}

	root = args[0]

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
