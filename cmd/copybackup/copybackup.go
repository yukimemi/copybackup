package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/umisama/golog"
	cb "github.com/yukimemi/copybackup"
	core "github.com/yukimemi/gocore"
)

type Options struct {
	Generation int    `short:"g" long:"generation" description:"バックアップする世代。" default:"5"`
	BackupDst  string `short:"b" long:"backup" description:"バックアップを保存する先。" default:"_old"`
	Sleep      int    `short:"s" long:"sleep" description:"バックアップ間隔。" default:"60*5"`
	Csv        string `short:"c" long:"csv" description:"設定ファイルから読み込み。(csvファイル)"`
}

func main() {
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)

	var e error
	var wg sync.WaitGroup

	cpus := runtime.NumCPU()
	semaphore := make(chan int, cpus)

	opts := new(Options)
	parser := flags.NewParser(opts, flags.Default)
	args, e := parser.Parse()
	core.FailOnError(e)

	if len(args) == 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(0)
	}

	for _, arg := range args {
		if f, e := os.Stat(arg); f.IsDir() {
			files, e := ioutil.ReadDir(arg)
			core.FailOnError(e)
			for _, f := range files {
				if !f.IsDir() {
					wg.Add(1)
					go func() {
						defer wg.Done()
						semaphore <- 1
						cg := cb.NewCopyGroup(filepath.Join(arg, f.Name()), opts.BackupDst, opts.Generation, opts.Sleep)
						cg.Backup()
						cg.DeleteOldFile()
						<-semaphore
					}()
				}
			}
		} else if e != nil {
			core.FailOnError(e)
		} else {
			cg := cb.NewCopyGroup(arg, opts.BackupDst, opts.Generation, opts.Sleep)
			wg.Add(1)
			go func() {
				defer wg.Done()
				semaphore <- 1
				cg.Backup()
				cg.DeleteOldFile()
				<-semaphore
			}()
		}
	}
	wg.Wait()
}
