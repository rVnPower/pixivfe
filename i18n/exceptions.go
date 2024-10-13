package i18n

var IgnoreTheseStrings = map[string]bool{
	"":                      true,
	"»":                     true,
	"▶":                     true,
	"⧉ {{ .Pages }}":        true,
	"PixivFE":               true,
	"pixiv.net/i/{{ .ID }}": true,
}
