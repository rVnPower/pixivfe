package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/soluble-ai/go-jnode"
	"github.com/yargevad/filepathx"
	"golang.org/x/net/html"
)

func main() {
	result := jnode.NewArrayNode()

	files, err := filepathx.Glob("assets/views/**/*.html")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file == "assets/views/temp.jet.html" {
			continue
		}
		processFile(file, result)
		// processFile_jet(file, result)
	}


	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}

var re_command_fullmatch = regexp.MustCompile(`\A([\s\n# =]*\{\{[^\{\}]*\}\})*[\s\n]*\z`)
var re_comment = regexp.MustCompile(`\{\*[\s\S]*?\*\}`)

func stripComments(s string) string {
	return re_comment.ReplaceAllString(s, "")
}

func processFile(filename string, result *jnode.Node) {
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

	visit(body, func(nodes []*html.Node) {
		builder := strings.Builder{}
		for _, node := range nodes {
			html.Render(&builder, node)
		}
		s := strings.TrimSpace(builder.String())

		if s != "" {
			object := jnode.NewObjectNode()
			object.Put("msg", s)
			object.Put("file", filename)
			result.Append(object)
		}
	})
}

func visit(node *html.Node, record func([]*html.Node)) {
	stash := []*html.Node{}
	for it := node.FirstChild; it != nil; it = it.NextSibling {
		switch it.Type {
		case html.TextNode:
			if re_command_fullmatch.MatchString(it.Data) {
				continue
			}
			// this node contains text, record the text
			stash = append(stash, it)
		case html.ElementNode:
			tag := it.Data

			if tag == "a" && containsOnlyText(it) {
				stash = append(stash, it)
			} else {
				clearTo(&stash, record)
				visit(it, record)
			}
		}
	}
	clearTo(&stash, record)
}

func containsOnlyText(node *html.Node) bool {
	for it := node.FirstChild; it != nil; it = it.NextSibling {
		switch it.Type {
		case html.TextNode:
			if re_command_fullmatch.MatchString(it.Data) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// func hasAttr(attrs []html.Attribute, key string) (bool, string) {
// 	for _, attr := range attrs {
// 		if attr.Key == key {
// 			return true, attr.Val
// 		}
// 	}
// 	return false, ""
// }

func clearTo(stash *[]*html.Node, record func([]*html.Node)) {
	all_element_node := true
	for _, node := range *stash {
		if node.Type != html.ElementNode {
			all_element_node = false
			break
		}
	}
	if all_element_node {
		for _, node := range *stash {
			visit(node, record)
		}
	} else {
		record(*stash)
	}
	*stash = []*html.Node{}
}
