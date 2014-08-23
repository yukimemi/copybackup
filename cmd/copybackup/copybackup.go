package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/yukimemi/copybackup"
)

func main() {
	var debug DebugT

	var root string
	var e error
	var wg sync.WaitGroup

	g := flag.Int("g", -1, "バックアップする世代。")
	b := flag.String("b", "_old", "バックアップを保存する先。絶対パスでの指定も可能。 ([デフォルト _old])")
	s := flag.Int("s", 60*5, "バックアップ間s。 (秒 [デフォルト 5分])")
	flag.Parse()

	debug.Println("g = ", *g)
	debug.Println("b = ", *b)
	debug.Println("s = ", *s)

	debug.PrintValue("args", flag.Args())
	debug.PrintValue("args", os.Args)
	if flag.NArg() != 0 {
		root = flag.Arg(0)
	} else {
		root, e = os.Getwd()
		FailOnError(e)
		flag.Usage()
		os.Exit(1)
	}

	debug.PrintValue("root", root)

	files, e := ioutil.ReadDir(root)
	FailOnError(e)

	for _, f := range files {
		if !f.IsDir() {
			src := filepath.Join(root, f.Name())
			dst, e := copybackup.MakeDstPath(src, "_old")
			FailOnError(e)
			wg.Add(1)
			go func(src, dst string) {
				defer wg.Done()
				copybackup.Backup(src, dst)
			}(src, dst)
		}
		wg.Wait()
	}
}
