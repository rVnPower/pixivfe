package template

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

type LocalizedFSLoader struct {
	Dir string
}

func NewLocalizedFSLoader(dir string) *LocalizedFSLoader {
	return &LocalizedFSLoader{
		Dir: dir,
	}
}

func (l *LocalizedFSLoader) Exists(templatePath string) bool {
	templatePath = filepath.Join(l.Dir, filepath.FromSlash(templatePath))
	stat, err := os.Stat(templatePath)
	if err == nil && !stat.IsDir() {
		return true
	}
	return false
}

func (l *LocalizedFSLoader) Open(templatePath string) (io.ReadCloser, error) {
	locale := i18n.GetLocale()
	i18n_path := path.Join(l.Dir, templatePath)
	templatePath = filepath.Join(l.Dir, filepath.FromSlash(templatePath))

	//println("load replacer:", i18n_path)

	replacer := i18n.Replacer(locale, i18n_path)
	if replacer == nil {
		return os.Open(templatePath)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(strings.NewReader(replacer.Replace(string(content)))), nil
}
