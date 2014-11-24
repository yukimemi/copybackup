package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/umisama/golog"
	cb "github.com/yukimemi/copybackup"
	core "github.com/yukimemi/gocore"

	"golang.org/x/text/encoding/japanese"
)

type Options struct { // {{{
	Generation int    `short:"g" long:"generation" description:"バックアップする世代。(デフォルト: 5)" default:"5"`
	BackupDst  string `short:"b" long:"backup" description:"バックアップを保存する先。(デフォルト: _old)" default:"_old"`
	Sleep      int    `short:"s" long:"sleep" description:"バックアップ間隔。[秒] (デフォルト: -1 (繰り返さない))" default:"-1"`
	Csv        string `short:"c" long:"csv" description:"設定ファイルから読み込み。(csvファイル)"`
	Help       bool   `short:"h" long:"help" description:"このヘルプを表示。"`
} // }}}

func showHelp() { // {{{

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
} // }}}

func normalizationCsv(array [][]string) [][]string {
	var result [][]string
	for _, a := range array {
		if len(a) == 1 {
			a = append(a, "10")
		}
		if len(a) == 2 {
			a = append(a, "_old")
		}
		result = append(result, a)
	}
	return result
}

func main() { // {{{
	core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Info)
	// core.Logger, _ = log.NewLogger(os.Stdout, log.TIME_FORMAT_SEC, log.LOG_FORMAT_POWERFUL, log.LogLevel_Debug)

	var e error
	var st int
	defer func() { os.Exit(st) }()
	var wg sync.WaitGroup

	cpus := runtime.NumCPU()
	core.Logger.Debugf("cpu num = [%d]", cpus)
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
		showHelp()
		return
	}

	var records [][]string

	if len(args) == 0 {
		if opts.Csv != "" {
			core.Logger.Infof("csv file = [%s]", opts.Csv)
			records, e = core.ImportCsv(opts.Csv, japanese.ShiftJIS.NewDecoder())
			core.FailOnError(e)
			records = normalizationCsv(records)
		} else {
			showHelp()
			return
		}
	} else {
		for _, arg := range args {
			records = append(records, []string{arg, string(opts.Generation), opts.BackupDst})
		}
	}

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT)
	for {
		for _, record := range records {
			gen, e := strconv.Atoi(record[1])
			core.FailOnError(e)
			if f, e := os.Stat(record[0]); f.IsDir() {
				files, e := ioutil.ReadDir(record[0])
				core.FailOnError(e)
				for _, f := range files {
					if !f.IsDir() {
						src := filepath.Join(record[0], f.Name())
						wg.Add(1)
						go func() {
							defer wg.Done()
							semaphore <- 1
							cg := cb.NewCopyGroup(src, record[2], gen)
							cg.Backup()
							cg.DeleteOldFile()
							<-semaphore
						}()
					}
				}
			} else if e != nil {
				core.FailOnError(e)
			} else {
				cg := cb.NewCopyGroup(record[0], record[2], gen)
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

		if opts.Sleep == -1 {
			return
		}
		fmt.Printf("%d 秒後に再バックアップ処理を行います。'C-c'を押すと終了します。\n", opts.Sleep)
		select {
		case <-time.After(time.Second * time.Duration(opts.Sleep)):
			break
		case <-s:
			fmt.Println("終了します...")
			return
			break
		}
	}
} // }}}
