package content

import (
	"bytes"
	"errors"
	"html/template"
	"sort"
	"strings"
	"time"

	"pubgo/config"

	"gopkg.in/yaml.v2"
)

// Entry represents a single content entry
type Entry struct {
	FileName string
	Body     template.HTML
	Page     string `yaml:"page"`

	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	Author      string    `yaml:"author"`
	Description string    `yaml:"description"`

	IncludeToc   bool `yaml:"include_toc"`
	ShowComments bool `yaml:"show_comments"`
}

func (e Entry) StaticFileName() string {
	return strings.Replace(e.FileName, ".md", ".html", 1)
}

// Entries is a slice of Entry
type Entries []Entry

// sort entries by date
func (e Entries) SortByDate() Entries {
	var entries Entries
	for _, entry := range e {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.UnixMicro() < entries[j].Date.UnixMicro()
	})
	return entries
}

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
