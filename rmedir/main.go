package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/badugisoft/tools/lib"
	"gopkg.in/urfave/cli.v2"
)

var context = struct {
	Recursive   bool
	ExcludeSelf bool
}{}

func main() {
	app := &cli.App{
		Name:    "rmedir",
		Version: "0.0.1",
		Usage:   "remove empty directory",
		Action:  rmedir,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Destination: &context.Recursive,
			},
			&cli.BoolFlag{
				Name:        "exclude-self",
				Aliases:     []string{"e"},
				Destination: &context.ExcludeSelf,
			},
		},
	}

	app.Run(os.Args)
}

func rmedir(c *cli.Context) (errRet error) {
	defer func() {
		if r := recover(); r != nil {
			errRet = cli.Exit(r.(error).Error(), 1)
		}
	}()

	if c.Args().Len() == 0 {
		cli.ShowAppHelp(c)
		lib.Panic("not enough arguments")
	}

	if context.ExcludeSelf && !context.Recursive {
		return
	}

	dirname := c.Args().First()

	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		lib.Panicf("dir not exist: '%v'", dirname)
	}
	lib.PanicfIf(err, "open dir failed: '%v'", dirname)

	if !info.IsDir() {
		lib.Panicf("not a directory: '%v'", dirname)
	}

	if rmedirRecur(dirname) && !context.ExcludeSelf {
		lib.PanicfIf(os.Remove(dirname), "remove dir failed : '%v'", dirname)
		fmt.Println("deleted: ", dirname)
	}

	return
}

func rmedirRecur(dirname string) bool {
	isEmpty := true

	files, err := ioutil.ReadDir(dirname)
	lib.PanicfIf(err, "read dir failed: '%v'", dirname)

	for _, f := range files {
		if f.IsDir() && context.Recursive {
			subdirname := filepath.Join(dirname, f.Name())
			if rmedirRecur(subdirname) {
				lib.PanicfIf(os.Remove(subdirname), "remove dir failed : '%v'", subdirname)
				fmt.Println("deleted:", subdirname)
			} else {
				isEmpty = false
			}
		} else {
			isEmpty = false
		}
	}

	return isEmpty
}
