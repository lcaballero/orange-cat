package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/noraesae/orange-cat"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"sort"
)

func main() {
	var paths sort.StringSlice
	css, paths, err := Css()
	if err != nil {
		panic(err)
	}
	paths.Sort()
	fmt.Printf("Css Path: %s\n", css)
	fmt.Printf("Path Count: %d\n", len(paths))
	for _, p := range paths {
		fmt.Println(p)
	}
}

func CssPath() string {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, ".orange")
	return dir
}

func Css() (string, []string, error) {
	dir := CssPath()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", nil, err
	}
	paths := make([]string, 0)
	for _,f := range files {
		path := filepath.Join(dir, f.Name())
		paths = append(paths, path)
	}
	return dir, paths, nil
}

func pp() {
	app := cli.NewApp()
	app.Name = "orange"
	app.Version = orange.Version
	app.Usage = `orange is a Markdown previewer written in Go.
   Its main goal is to be used with any editor you love.
   For information, please visit https://github.com/noraesae/orange-cat`
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "basic, b",
			Usage: "Use Markdown Basic(Markdown Common by default).",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 6060,
			Usage: "Port to listen.",
		},
	}
	app.Action = func(c *cli.Context) {
		args := c.Args()

		orange := orange.NewOrange(c.Int("port"))

		if c.Bool("basic") {
			orange.UseBasic()
		}

		orange.Run(args...)
	}

	// codegangsta/cli help template
	cli.AppHelpTemplate = `orange-cat
   {{.Usage}}

USAGE:
   {{.Name}} [global options] [command] file

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}

GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

	app.Run(os.Args)
}
