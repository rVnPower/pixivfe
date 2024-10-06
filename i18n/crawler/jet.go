package main

import (
	"strings"

	"github.com/CloudyKit/jet/v6"
	jetutils "github.com/CloudyKit/jet/v6/utils"
	"github.com/soluble-ai/go-jnode"
)

func ProcessFile_jet(filename string, result *jnode.Node) {
	parser := jet.NewSet(jet.NewOSFileSystemLoader("."))
	template, err := parser.GetTemplate(filename)
	if err != nil {
		panic(err)
	}

	jetutils.Walk(template, jetutils.VisitorFunc(func(vc jetutils.VisitorContext, node_ jet.Node) {
		switch node := node_.(type) {
		case *jet.TextNode:
			if re_command.Match(node.Text) {
				// is command, skip
			} else {
				s := strings.TrimSpace(string(node.Text))
				// s = strings.ReplaceAll(s, "\n", " ")
				// println(s)
				object := jnode.NewObjectNode()
				object.Put("msg", s)
				object.Put("file", filename)
				object.Put("line", node.Line)
				object.Put("offset", int(node.Pos))
				result.Append(object)
			}
		case *jet.ListNode:
			for _, n := range node.Nodes {
				vc.Visit(n)
			}
		}
	}))
}
