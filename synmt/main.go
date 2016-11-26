package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	cli "gopkg.in/urfave/cli.v2"
)

var context = struct {
	// Recursive bool
}{}

func main() {
	app := &cli.App{
		Name:   "synmt",
		Usage:  "synchronize modified time",
		Action: synmt,
		Flags:  []cli.Flag{
		// &cli.BoolFlag{
		// 	Name:        "recursive",
		// 	Aliases:     []string{"r"},
		// 	Destination: &context.Recursive,
		// },
		},
	}

	app.Run(os.Args)
}

func synmt(c *cli.Context) (errRet error) {
	defer func() {
		if r := recover(); r != nil {
			errRet = cli.Exit(r.(error).Error(), -1)
		}
	}()

	if c.Args().Len() == 0 {
		panicf("not enough arguments")
	}

	dirname := c.Args().First()

	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		panicf("dir not exist: '%v'", dirname)
	}
	_panicf(err, "open dir failed: '%v'", dirname)

	if !info.IsDir() {
		panicf("not a directory: '%v'", dirname)
	}

	modified := synmtRecur(dirname, info.ModTime())
	if modified.Before(info.ModTime()) {
		_panicf(os.Chtimes(dirname, modified, modified),
			"change time failed: '%v'", dirname)
	}
	return
}

var minTime = time.Unix(0, 0)

func synmtRecur(dirname string, dirModified time.Time) time.Time {
	files, err := ioutil.ReadDir(dirname)
	_panicf(err, "read dir failed: '%v'", dirname)

	maxModified := minTime

	for _, f := range files {
		modified := f.ModTime()

		if f.IsDir() {
			subdirname := filepath.Join(dirname, f.Name())
			modified = synmtRecur(subdirname, modified)
			if modified.Before(f.ModTime()) {
				_panicf(os.Chtimes(subdirname, modified, modified),
					"change time failed: '%v'", subdirname)
			}
		}

		if modified.After(maxModified) {
			maxModified = modified
		}
	}

	if maxModified.Equal(minTime) {
		return dirModified
	}

	return maxModified
}
