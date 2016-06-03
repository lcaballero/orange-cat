package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"io/ioutil"
	"sort"
	"strings"
	"encoding/json"
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

func Links() []string {
	links := make([]string, 0)
	_, paths, err := Css()
	if err != nil {
		panic(err)
	}
	for _,p := range paths {
		link := fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="%s" />`, p)
		links = append(links, link)
	}
	return links
}

func WriteMd(path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	mdPath := filepath.Join(cwd, path)
	fmt.Printf("Markdown file: %s\n", mdPath)

	_, err = os.Stat(mdPath)
	if err != nil && os.IsNotExist(err) {
		return ""
	}

	mdBytes, err := ioutil.ReadFile(mdPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	b := MdConverter.Convert(mdBytes)
	return string(b)
}

func Template(w http.ResponseWriter, req *http.Request, filepath string) {
	var style string = strings.Join(Links(), "\n")
	fmt.Println("styles:")
	fmt.Println(style)
	urlJson, err := json.MarshalIndent(req.URL, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("req.URL", string(urlJson))
	markdownHtml := WriteMd("/" + req.URL.Path)

	templateStr := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset='UTF-8' />
  <title>%s</title>
  %s
</head>
<body>
  <div id='md' class='markdown-body'>%s</div>
</body>
</html>`, filepath, style, markdownHtml)

	w.Write([]byte(templateStr))
}

var websocket_script = `
  <script>
    (function () {
      var markdown = document.getElementById("md");
      var conn = new WebSocket("ws://" + location.host + "/%[1]s");
      conn.onmessage = function (evt) {
        markdown.innerHTML = evt.data;
      };
    })();
  </script>
`
