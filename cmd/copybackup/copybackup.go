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
	Generation int    `short:"g" long:"generation" description:"バックアップする世代。 (デフォルト [5])" default:"5"`
	BackupDst  string `short:"b" long:"backup" description:"バックアップを保存する先。 (デフォルト [_old])" default:"_old"`
	Sleep      int    `short:"s" long:"sleep" description:"バックアップ間隔。 (秒 [デフォルト 5分])" default:"60*5"`
	Csv        string `short:"c" long:"csv" description:"設定ファイルから読み込み。(csvファイル)"`
}

func main() {
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)

	var e error
	var wg sync.WaitGroup
	var cg *cb.CopyGroup

	opts := &Options{}
	parser := flags.NewParser(opts, flags.Default)
	args, e := parser.Parse()
	core.FailOnError(e)

	if len(args) == 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(0)
	}

	for _, arg := range args {
		if f, _ := os.Stat(arg); f.IsDir() {
			files, e := ioutil.ReadDir(arg)
			core.FailOnError(e)
			for _, f := range files {
				if !f.IsDir() {
					cg = cb.NewCopyGroup(filepath.Join(arg, f.Name()), opts.BackupDst, opts.Generation, opts.Sleep)
					wg.Add(1)
					go func() {
						defer wg.Done()
						cg.Backup()
					}()
				}
			}

		} else {
			cg = cb.NewCopyGroup(arg, opts.BackupDst, opts.Generation, opts.Sleep)
			wg.Add(1)
			go func() {
				defer wg.Done()
				cg.Backup()
			}()
		}
	}
	wg.Wait()

}
