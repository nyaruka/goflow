package i18n

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Library is a collection of PO files providing translations in different languages
type Library struct {
	path        string
	srcLanguage string
}

// NewLibrary creates new library from directory structure in path
func NewLibrary(path, srcLanguage string) *Library {
	return &Library{path: path, srcLanguage: srcLanguage}
}

// Path returns the root path of this library
func (l *Library) Path() string {
	return l.path
}

// SrcLanguage returns the source language of this library
func (l *Library) SrcLanguage() string {
	return l.srcLanguage
}

// Update updates the message IDs in the default language from the given PO,
// and merges those changes into the other PO files
func (l *Library) Update(domain string, pot *PO) error {
	// update the PO file for the source language (i.e. our POT)
	f, err := os.Create(l.poPath(l.srcLanguage, domain))
	if err != nil {
		return err
	}

	defer f.Close()
	pot.Write(f)

	// merge the ID changes into the PO files for the translation languages
	for _, lc := range l.locales(false) {
		args := []string{
			"-q",
			"--previous",
			l.poPath(lc, domain),
			l.poPath(l.srcLanguage, domain),
			"-o",
			l.poPath(lc, domain),
			"--no-wrap",
			"--sort-output",
		}

		b := &strings.Builder{}

		cmd := exec.Command("msgmerge", args...)
		cmd.Stderr = b
		if err := cmd.Run(); err != nil {
			return errors.New(b.String())
		}
	}

	return nil
}

// Load loads the PO for the given language and domain
func (l *Library) Load(language, domain string) (*PO, error) {
	f, err := os.Open(l.poPath(language, domain))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadPO(f)
}

// Locales returns the names of the locales included in this library
func (l *Library) Locales() []string {
	return l.locales(true)
}

func (l *Library) locales(includeSrc bool) []string {
	directory, err := os.ReadDir(l.path)
	if err != nil {
		panic(err)
	}

	var languages []string
	for _, file := range directory {
		if file.IsDir() {
			lang := file.Name()
			if includeSrc || lang != l.srcLanguage {
				languages = append(languages, lang)
			}
		}
	}

	return languages
}

// returns the path of the PO file for the given language and domain
func (l *Library) poPath(language, domain string) string {
	return path.Join(l.Path(), language, domain+".po")
}
