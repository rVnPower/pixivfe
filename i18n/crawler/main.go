package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/soluble-ai/go-jnode"
	"github.com/yargevad/filepathx"
	"golang.org/x/net/html"
)

var re_command = regexp.MustCompile(`\A([\p{Zs}\n]*\{\{[^\{\}]*\}\})*[\p{Zs}\n]*\z`)

func testRegex() {
	if !re_command.MatchString(`{{- extends "layout/default" }}
{{- block body() }}`) {
		panic("regex is broken")
	}
}

func main() {
	testRegex()

	result := jnode.NewArrayNode()

	files, err := filepathx.Glob("assets/views/**/*.html")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		// println(file)
		processFile(file, result)
	}

	os.Stdout.WriteString(result.String())
}

func processFile(filename string, result *jnode.Node) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	root, err := html.Parse(f)
	if err != nil {
		panic(err)
	}

	body := cascadia.Query(root, cascadia.MustCompile("body"))
	if body == nil {
		panic("can't find <body>")
	}

	recurse(body, func(node *html.Node) {
		s := strings.TrimSpace(node.Data)
		// s = strings.ReplaceAll(s, "\n", " ")
		// println(s)
		object := jnode.NewObjectNode()
		object.Put("msg", s)
		object.Put("file", filename)
		result.Append(object)
	})
}

func recurse(node *html.Node, record func(*html.Node)) {
	for it := node.FirstChild; it != nil; it = it.NextSibling {
		switch it.Type {
		case html.TextNode:
			if re_command.MatchString(it.Data) {
				// is command, skip
				continue
			}
			// this node contains text, record the text
			record(it)

		case html.ElementNode:
			recurse(it, record)
		}
	}
}
