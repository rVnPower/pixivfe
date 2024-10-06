package main

import (
	"os"
	"regexp"

	"github.com/soluble-ai/go-jnode"
	"github.com/yargevad/filepathx"
)

var re_command_fullmatch = regexp.MustCompile(`\A([\s\n# =]*\{\{[^\{\}]*\}\})*[\s\n]*\z`)
var re_comment = regexp.MustCompile(`\{\*[\s\S]*?\*\}`)

func stripComments(s string) string {
	return re_comment.ReplaceAllString(s, "")
}

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
		ProcessFile_html(file, result)
		// ProcessFile_jet(file, result)
	}

	os.Stdout.WriteString(result.String())
}
