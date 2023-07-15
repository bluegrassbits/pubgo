package content

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"time"

	"pubgo/config"

	"gopkg.in/yaml.v2"
)

// Entry represents a single content entry
type Entry struct {
	FileName string
	Body     template.HTML

	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Date        time.Time `yaml:"date"`
	Author      string    `yaml:"author"`
	Page        string    `yaml:"page"`
}

func (e Entry) StaticFileName() string {
	return strings.Replace(e.FileName, ".md", ".html", 1)
}

// Entries is a slice of Entry
type Entries []Entry

// Page is the data passed to the template
type Content struct {
	Site        config.Site
	Page        config.Page
	BasePath    string
	RequestPath string
	Mode        string
	Title       string
	Collection  bool
	Entry       Entry
	Entries     Entries
}

// ParseEntry parses a file and returns an Entry struct
// and the remaining data

func ParseEntry(data []byte) (Entry, []byte, error) {
	entry := Entry{Title: "No Title"}

	yamlStart := bytes.Index(data, []byte("---"))
	if yamlStart < 0 {
		return entry, data, errors.New("No yaml found")
	}
	yamlEnd := bytes.Index(data[yamlStart+3:], []byte("---"))
	if yamlEnd < 0 {
		return entry, data, errors.New("No yaml found")
	}
	err := yaml.Unmarshal(data[yamlStart+3:yamlStart+3+yamlEnd], &entry)

	// Get data after yaml
	data = data[yamlStart+3+yamlEnd+3:]

	if err != nil {
		return entry, data, err
	}
	return entry, data, nil
}
