package i18n

import (
	"io/ioutil"
	"os"
	"path"
)

// Library is a collection of PO files providing translations in different languages
type Library struct {
	path string
}

func NewLibrary(path string) *Library {
	return &Library{path: path}
}

func (l *Library) Path() string {
	return l.path
}

func (l *Library) Save(language, domain string, po *PO) error {
	f, err := os.Create(l.poPath(language, domain))
	if err != nil {
		return err
	}

	defer f.Close()
	po.Write(f)
	return nil
}

func (l *Library) Load(language, domain string) (*PO, error) {
	f, err := os.Open(l.poPath(language, domain))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadPO(f)
}

func (l *Library) poPath(language, domain string) string {
	return path.Join(l.Path(), language, domain+".po")
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
