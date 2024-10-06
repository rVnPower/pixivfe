package main

import (
	"os"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/soluble-ai/go-jnode"
	"golang.org/x/net/html"
)

func ProcessFile_html(filename string, result *jnode.Node) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	content2 := stripComments(string(content))

	root, err := html.Parse(strings.NewReader(content2))
	if err != nil {
		panic(err)
	}

	body := cascadia.Query(root, cascadia.MustCompile("body"))
	if body == nil {
		panic("can't find <body>")
	}

	recurse_html(body, func(node *html.Node) {
		s := strings.TrimSpace(node.Data)

		if s != "" {
			object := jnode.NewObjectNode()
			object.Put("msg", s)
			object.Put("file", filename)
			result.Append(object)
		}
	})
}

func recurse_html(node *html.Node, record func(*html.Node)) {
	for it := node.FirstChild; it != nil; it = it.NextSibling {
		switch it.Type {
		case html.TextNode:
			if re_command_fullmatch.MatchString(it.Data) {
				// is command, skip
				continue
			}
			// this node contains text, record the text
			record(it)

		case html.ElementNode:
			recurse_html(it, record)
		}
	}
}
