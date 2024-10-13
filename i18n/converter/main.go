package main

import (
	"encoding/json"
	"os"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"github.com/soluble-ai/go-jnode"
)

func main() {
	root := &jnode.Node{}
	err := json.NewDecoder(os.Stdin).Decode(root)
	if err != nil {
		panic(err)
	}
	translation_map := make(map[string]string)
	for i := 0; i < root.Size(); i++ {
		o := root.Get(i)
		msg := o.ToMap()["msg"].(string)
		file := o.ToMap()["file"].(string)
		translation_map[i18n.SuccintId(file, msg)] = msg
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	encoder.Encode(translation_map)
}
