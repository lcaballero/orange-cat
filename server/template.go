package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"text/template"
	"io/ioutil"
	"sort"
)


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

func ShowCssPaths() {
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

func Template(w http.ResponseWriter, filepath string) {
	var style string = `<style type="text/css">` + DefaultStyle + "</style>"
	var customStyle string
	if css, err := CustomCSS(); err == nil {
		customStyle = *css
	}
	fmt.Println(style)
	fmt.Println(customStyle)

	templateStr := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset='UTF-8' />
  <title>%s</title>
  %s
  %s
</head>
<body>
  <div id='md' class='markdown-body'></div>
  <script>
    (function () {
      var markdown = document.getElementById("md");
      var conn = new WebSocket("ws://" + location.host + "/%[1]s");
      conn.onmessage = function (evt) {
        markdown.innerHTML = evt.data;
      };
    })();
  </script>
</body>`, filepath, style, customStyle)

	var (
		t   *template.Template
		err error
	)

	if t, err = template.New("template").Parse(templateStr); err != nil {
		panic(err)
	}

	if err = t.Execute(w, nil); err != nil {
		panic(err)
	}
}

func CustomCSS() (*string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	customCSSPath := filepath.Join(usr.HomeDir, ".orange/orange-cat.css")

	stat, err := os.Stat(customCSSPath)
	if err != nil || !stat.Mode().IsRegular() {
		return nil, errors.New("No custom CSS")
	}

	customCSS := "<link rel='stylesheet' href='" + customCSSPath + "' />"
	return &customCSS, nil
}

var DefaultStyle = `
