package i18n

import (
	"io/ioutil"
	"path"
)

type Library struct {
	path string
}

func NewLibrary(path string) *Library {
	return &Library{path: path}
}

func (l *Library) Path() string {
	return l.path
}

func (l *Library) POPath(language, domain string) string {
	return path.Join(l.Path(), language, domain+".po")
}

func (l *Library) Activate(language, domain string) {
	//gotext.Configure(l.path, language, domain)
}

func (l *Library) Languages() []string {
	directory, err := ioutil.ReadDir(l.path)
	if err != nil {
		panic(err)
	}

	var languages []string
	for _, file := range directory {
		if file.IsDir() {
			languages = append(languages, file.Name())
		}
	}

	return languages
}
