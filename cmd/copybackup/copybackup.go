package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/umisama/golog"
	cb "github.com/yukimemi/copybackup"
	core "github.com/yukimemi/gocore"
)

type Options struct {
	Generation int    `short:"g" long:"generation" description:"バックアップする世代。(デフォルト: 5)" default:"5"`
	BackupDst  string `short:"b" long:"backup" description:"バックアップを保存する先。(デフォルト: _old)" default:"_old"`
	Sleep      int    `short:"s" long:"sleep" description:"バックアップ間隔。[秒] (デフォルト: 300)" default:"60*5"`
	Csv        string `short:"c" long:"csv" description:"設定ファイルから読み込み。(csvファイル)"`
	Help       bool   `short:"h" long:"help" description:"このヘルプを表示。"`
}

func showHelp() {

	os.Stderr.WriteString(` 
Usage: copybackup [options] [FILE | DIRECTORY | input csv file]

Options:
`)

	t := reflect.TypeOf(Options{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag

		var o string
		if s := tag.Get("short"); s != "" {
			o = fmt.Sprintf("-%s, --%s", tag.Get("short"), tag.Get("long"))
		} else {
			o = fmt.Sprintf("--%s", tag.Get("long"))
		}

		fmt.Fprintf(
			os.Stderr,
			"  %-21s %s\n",
			o,
			tag.Get("description"),
		)
	}
}

func main() {
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)

	var e error
	var st int
	defer func() { os.Exit(st) }()
	var wg sync.WaitGroup

	cpus := runtime.NumCPU()
	semaphore := make(chan int, cpus)

	opts := new(Options)
	parser := flags.NewParser(opts, flags.PrintErrors)

	args, e := parser.Parse()
	if e != nil {
		st = 1
		showHelp()
		return
	}

	if opts.Help {
		// parser.WriteHelp(os.Stderr)
		showHelp()
		return
	}

	if len(args) == 0 {
		// parser.WriteHelp(os.Stderr)
		showHelp()
		return
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
